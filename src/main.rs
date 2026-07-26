use clap::{Args, CommandFactory, Parser, Subcommand};
use godi::{
    CancellationToken, CommonOptions, CopyOptions, Event, FileFilter, Importance, SealFormat,
    SealOptions, VerifyOptions, seal, sealed_copy, verify,
};
use signal_hook::consts::{SIGINT, SIGTERM};
use std::{num::NonZeroUsize, path::PathBuf, process::ExitCode};

#[derive(Parser)]
#[command(
    name = "godi",
    version = "v1.1.0",
    author = "Sebastian Thiel & Contributors",
    about = "Verify data integrity and transfer data securely at highest speeds.",
    disable_help_subcommand = false
)]
struct Cli {
    #[arg(
        long = "streams-per-input-device",
        alias = "spid",
        default_value_t = 1,
        global = true
    )]
    input_streams: usize,
    #[arg(long, default_value = "error", global = true)]
    verbosity: Verbosity,
    #[arg(
        long = "file-exclude-patterns",
        default_value = "VOLATILE",
        global = true
    )]
    filters: String,
    #[command(subcommand)]
    command: Option<Command>,
}

#[derive(Subcommand)]
enum Command {
    /// Generate a seal for one or more files or directories.
    Seal(SealArgs),
    /// Seal data while copying it to one or more destinations.
    #[command(name = "sealed-copy")]
    SealedCopy(CopyArgs),
    /// Compare files on disk with one or more seals.
    Verify(VerifyArgs),
    #[cfg(feature = "web")]
    /// Launch the web frontend.
    Web(WebArgs),
}

#[derive(Args)]
struct SealArgs {
    #[arg(long, default_value = "gob")]
    format: SealFormat,
    #[arg(required = true)]
    sources: Vec<PathBuf>,
}

#[derive(Args)]
struct CopyArgs {
    #[arg(long)]
    verify: bool,
    #[arg(
        long = "streams-per-output-device",
        alias = "spod",
        default_value_t = 1
    )]
    output_streams: usize,
    #[arg(long, default_value = "gob")]
    format: SealFormat,
    #[arg(required = true, num_args = 2.., allow_hyphen_values = true)]
    paths: Vec<PathBuf>,
}

#[derive(Args)]
struct VerifyArgs {
    #[arg(required = true)]
    seals: Vec<PathBuf>,
}

#[cfg(feature = "web")]
#[derive(Args)]
struct WebArgs {
    #[arg(long = "no-show")]
    no_show: bool,
    #[arg(short = 'a', long = "address", default_value = "localhost:9078")]
    address: String,
}

#[derive(Clone, Copy)]
enum Verbosity {
    Statistics,
    Info,
    Warn,
    Error,
    Result,
    Off,
}

impl std::str::FromStr for Verbosity {
    type Err = String;

    fn from_str(value: &str) -> Result<Self, Self::Err> {
        match value {
            "statistics" => Ok(Self::Statistics),
            "info" => Ok(Self::Info),
            "warn" => Ok(Self::Warn),
            "error" => Ok(Self::Error),
            "result" => Ok(Self::Result),
            "off" => Ok(Self::Off),
            _ => Err(format!("Unknown verbosity level: '{value}'")),
        }
    }
}

fn main() -> ExitCode {
    let args = legacy_aliases(std::env::args_os().collect());
    let cli = match Cli::try_parse_from(args) {
        Ok(cli) => cli,
        Err(err) => {
            let code = if err.use_stderr() { 3 } else { 0 };
            let _ = err.print();
            return ExitCode::from(code);
        }
    };

    #[cfg(feature = "web")]
    if cli.command.is_none() {
        return run_web("localhost:9078".into(), true);
    }
    let Some(command) = cli.command else {
        let _ = Cli::command().print_help();
        println!();
        return ExitCode::SUCCESS;
    };

    let streams = match NonZeroUsize::new(cli.input_streams) {
        Some(value) => value,
        None => {
            eprintln!("--streams-per-input-device must not be smaller than 1");
            return ExitCode::from(3);
        }
    };
    let filters = match parse_filters(&cli.filters) {
        Ok(filters) => filters,
        Err(err) => {
            eprintln!("{err}");
            return ExitCode::from(3);
        }
    };
    let cancellation = CancellationToken::default();
    if let Err(err) = signal_hook::flag::register(SIGINT, cancellation_flag(&cancellation))
        .and_then(|_| signal_hook::flag::register(SIGTERM, cancellation_flag(&cancellation)))
    {
        eprintln!("{err}");
        return ExitCode::from(3);
    }
    let common = CommonOptions {
        input_streams: streams,
        filters,
        cancellation,
    };
    let mut log = |event: &Event| log_event(event, cli.verbosity);

    let result = match command {
        Command::Seal(args) => seal(
            &args.sources,
            SealOptions {
                common,
                format: args.format,
            },
            &mut log,
        ),
        Command::SealedCopy(args) => {
            let split = args
                .paths
                .iter()
                .position(|p| p == std::path::Path::new("__godi_destinations__"));
            let (sources, destinations) = match split {
                Some(index) if index == 0 || index == args.paths.len() - 1 => {
                    eprintln!("sources and destinations must surround --");
                    return ExitCode::from(3);
                }
                Some(index) => (&args.paths[..index], &args.paths[index + 1..]),
                None if args.paths.len() == 2 => (&args.paths[..1], &args.paths[1..]),
                None => {
                    eprintln!("specify sources -- destinations");
                    return ExitCode::from(3);
                }
            };
            let Some(output_streams) = NonZeroUsize::new(args.output_streams) else {
                eprintln!("--streams-per-output-device must not be smaller than 1");
                return ExitCode::from(3);
            };
            sealed_copy(
                sources,
                destinations,
                CopyOptions {
                    seal: SealOptions {
                        common,
                        format: args.format,
                    },
                    output_streams,
                    verify_after_copy: args.verify,
                },
                &mut log,
            )
        }
        Command::Verify(args) => verify(
            &args.seals,
            VerifyOptions {
                input_streams: streams,
                cancellation: common.cancellation,
            },
            &mut log,
        ),
        #[cfg(feature = "web")]
        Command::Web(args) => return run_web(args.address, !args.no_show),
    };
    if result.is_ok() {
        ExitCode::SUCCESS
    } else {
        ExitCode::from(1)
    }
}

fn parse_filters(value: &str) -> Result<Vec<FileFilter>, String> {
    if value.is_empty() {
        return Ok(Vec::new());
    }
    value
        .split(',')
        .map(|value| FileFilter::parse(value).map_err(|err| err.to_string()))
        .collect()
}

fn log_event(event: &Event, verbosity: Verbosity) {
    let show = match verbosity {
        Verbosity::Off => false,
        Verbosity::Statistics => event.importance >= Importance::Statistics,
        Verbosity::Info => true,
        Verbosity::Warn => event.importance >= Importance::Warn,
        Verbosity::Error => event.importance >= Importance::Error,
        Verbosity::Result => event.importance == Importance::Result,
    };
    if show {
        if event.error.is_some() {
            eprintln!("{}", event.message);
        } else {
            println!("{}", event.message);
        }
    }
}

fn cancellation_flag(token: &CancellationToken) -> std::sync::Arc<std::sync::atomic::AtomicBool> {
    // signal-hook needs the atomic itself; mirror it into the public token.
    let flag = std::sync::Arc::new(std::sync::atomic::AtomicBool::new(false));
    let source = flag.clone();
    let token = token.clone();
    std::thread::spawn(move || {
        while !source.load(std::sync::atomic::Ordering::Relaxed) {
            std::thread::park_timeout(std::time::Duration::from_millis(25));
        }
        token.cancel();
    });
    flag
}

fn legacy_aliases(mut args: Vec<std::ffi::OsString>) -> Vec<std::ffi::OsString> {
    let sealed_copy = args.iter().position(|arg| arg == "sealed-copy");
    for arg in &mut args {
        if arg == "-spid" {
            *arg = "--spid".into();
        } else if arg == "-spod" {
            *arg = "--spod".into();
        }
    }
    if let Some(index) = sealed_copy
        && let Some(separator) = args[index + 1..].iter_mut().find(|arg| *arg == "--")
    {
        *separator = "__godi_destinations__".into();
    }
    args
}

#[cfg(feature = "web")]
fn run_web(address: String, show: bool) -> ExitCode {
    match tokio::runtime::Runtime::new()
        .map_err(|err| err.to_string())
        .and_then(|runtime| runtime.block_on(godi::web::serve(&address, show)))
    {
        Ok(()) => ExitCode::SUCCESS,
        Err(err) => {
            eprintln!("{err}");
            ExitCode::from(3)
        }
    }
}
