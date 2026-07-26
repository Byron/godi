use crate::{SealEntry, SealFormat, read_seal, write_seal};
use chrono::Local;
use glob::Pattern;
use md5::Md5;
use rayon::prelude::*;
use sha1::{Digest, Sha1};
use std::{
    fmt, fs,
    fs::{File, OpenOptions},
    io::{self, Read, Write},
    num::NonZeroUsize,
    path::{Component, Path, PathBuf},
    sync::{
        Arc,
        atomic::{AtomicBool, Ordering},
    },
    time::{Duration, Instant},
};

#[derive(Debug)]
pub struct GodiError(pub String);

impl fmt::Display for GodiError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.0)
    }
}

impl std::error::Error for GodiError {}

impl From<io::Error> for GodiError {
    fn from(value: io::Error) -> Self {
        Self(value.to_string())
    }
}

#[derive(Clone, Default)]
pub struct CancellationToken(Arc<AtomicBool>);

impl CancellationToken {
    pub fn cancel(&self) {
        self.0.store(true, Ordering::Release);
    }

    pub fn is_cancelled(&self) -> bool {
        self.0.load(Ordering::Acquire)
    }
}

#[derive(Clone, Debug)]
pub enum FileFilter {
    Symlinks,
    Hidden,
    Seals,
    Volatile,
    Glob(Pattern),
}

impl FileFilter {
    pub fn parse(value: &str) -> Result<Self, GodiError> {
        match value {
            "SYMLINK" => Ok(Self::Symlinks),
            "HIDDEN" => Ok(Self::Hidden),
            "SEALS" => Ok(Self::Seals),
            "VOLATILE" => Ok(Self::Volatile),
            _ => Pattern::new(value)
                .map(Self::Glob)
                .map_err(|err| GodiError(err.to_string())),
        }
    }

    pub(crate) fn matches(&self, name: &std::ffi::OsStr, metadata: &fs::Metadata) -> bool {
        let name = name.to_string_lossy();
        match self {
            Self::Symlinks => metadata.file_type().is_symlink(),
            Self::Hidden => name.starts_with('.'),
            Self::Seals => is_seal_name(&name),
            Self::Volatile => {
                (!metadata.is_dir() && !metadata.file_type().is_symlink() && !metadata.is_file())
                    || name == ".DS_Store"
                    || (metadata.is_dir()
                        && (name.starts_with(".Trash")
                            || name == ".fseventsd"
                            || name == ".TemporaryItems"
                            || name.starts_with(".DocumentRevisions")
                            || name == "lost+found"
                            || name.starts_with(".Spotlight")
                            || name == "System Volume Information"
                            || name == "$Recycle.Bin"))
            }
            Self::Glob(pattern) => pattern.matches(&name),
        }
    }
}

impl fmt::Display for FileFilter {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Symlinks => f.write_str("SYMLINK"),
            Self::Hidden => f.write_str("HIDDEN"),
            Self::Seals => f.write_str("SEALS"),
            Self::Volatile => f.write_str("VOLATILE"),
            Self::Glob(pattern) => pattern.fmt(f),
        }
    }
}

#[derive(Clone)]
pub struct CommonOptions {
    pub input_streams: NonZeroUsize,
    pub filters: Vec<FileFilter>,
    pub cancellation: CancellationToken,
}

impl Default for CommonOptions {
    fn default() -> Self {
        Self {
            input_streams: NonZeroUsize::MIN,
            filters: vec![FileFilter::Volatile],
            cancellation: CancellationToken::default(),
        }
    }
}

#[derive(Clone, Default)]
pub struct SealOptions {
    pub common: CommonOptions,
    pub format: SealFormat,
}

#[derive(Clone)]
pub struct CopyOptions {
    pub seal: SealOptions,
    pub output_streams: NonZeroUsize,
    pub verify_after_copy: bool,
}

impl Default for CopyOptions {
    fn default() -> Self {
        Self {
            seal: SealOptions::default(),
            output_streams: NonZeroUsize::MIN,
            verify_after_copy: false,
        }
    }
}

#[derive(Clone)]
pub struct VerifyOptions {
    pub input_streams: NonZeroUsize,
    pub cancellation: CancellationToken,
}

impl Default for VerifyOptions {
    fn default() -> Self {
        Self {
            input_streams: NonZeroUsize::MIN,
            cancellation: CancellationToken::default(),
        }
    }
}

#[derive(Clone, Copy, Debug, Eq, Ord, PartialEq, PartialOrd)]
#[repr(u8)]
pub enum Importance {
    Info,
    Warn,
    Error,
    Statistics,
    Result,
}

#[derive(Clone, Debug)]
pub struct Event {
    pub message: String,
    pub importance: Importance,
    pub error: Option<String>,
}

#[derive(Clone, Debug, Default)]
pub struct Statistics {
    pub files_read: u64,
    pub files_written: u64,
    pub bytes_read: u64,
    pub bytes_written: u64,
    pub skipped: u64,
    pub errors: u64,
    pub undone: u64,
    pub cancelled: bool,
    pub elapsed: Duration,
}

#[derive(Clone, Debug, Default)]
pub struct OperationReport {
    pub seals: Vec<PathBuf>,
    pub statistics: Statistics,
}

#[derive(Clone)]
struct SourceFile {
    root: PathBuf,
    path: PathBuf,
    relative: PathBuf,
    mode: u32,
    size: i64,
    symlink: bool,
}

struct Processed {
    source: SourceFile,
    sha1: Vec<u8>,
    md5: Vec<u8>,
    copied: Vec<Option<PathBuf>>,
    errors: Vec<Option<String>>,
    bytes: u64,
}

pub fn seal(
    sources: &[PathBuf],
    options: SealOptions,
    mut emit: impl FnMut(&Event),
) -> Result<OperationReport, GodiError> {
    let started = Instant::now();
    let (sources, files, skipped) = discover(sources, &options.common.filters)?;
    let processed = process(files, &[], &options.common, 1);
    let mut report = OperationReport::default();
    report.statistics.skipped = skipped;
    for root in sources {
        let selected: Vec<_> = processed.iter().filter(|p| p.source.root == root).collect();
        if let Some(error) = selected.iter().find_map(|p| p.errors[0].as_ref()) {
            report.statistics.errors += 1;
            event(&mut emit, error, Importance::Error, true);
            continue;
        }
        let entries: Vec<_> = selected
            .iter()
            .map(|p| p.entry_at(&p.source.path))
            .collect();
        let path = index_path(&root, options.format);
        match write_seal(&path, options.format, &entries) {
            Ok(()) => {
                event(
                    &mut emit,
                    &format!("Wrote seal file to '{}'", path.display()),
                    Importance::Result,
                    false,
                );
                report.seals.push(path);
            }
            Err(err) => {
                report.statistics.errors += 1;
                event(&mut emit, &err.to_string(), Importance::Error, true);
            }
        }
    }
    fill_stats(
        &processed,
        &mut report.statistics,
        started,
        &options.common.cancellation,
    );
    finish_event(&mut emit, "SEAL", &report.statistics);
    if report.statistics.errors > 0 {
        Err(GodiError("seal failed".into()))
    } else {
        Ok(report)
    }
}

pub fn sealed_copy(
    sources: &[PathBuf],
    destinations: &[PathBuf],
    options: CopyOptions,
    mut emit: impl FnMut(&Event),
) -> Result<OperationReport, GodiError> {
    if destinations.is_empty() {
        return Err(GodiError("Please provide at least one destination".into()));
    }
    let destinations = parse_destinations(destinations, sources)?;
    let started = Instant::now();
    let (_, files, skipped) = discover(sources, &options.seal.common.filters)?;
    let processed = process(
        files,
        &destinations,
        &options.seal.common,
        options.output_streams.get(),
    );
    let mut report = OperationReport::default();
    report.statistics.skipped = skipped;

    for (index, destination) in destinations.iter().enumerate() {
        let failed = processed
            .iter()
            .any(|p| p.errors.get(index).and_then(Option::as_ref).is_some());
        if failed || options.seal.common.cancellation.is_cancelled() {
            report.statistics.errors += 1;
            rollback(
                &processed,
                index,
                destination,
                &mut report.statistics,
                &mut emit,
            );
            continue;
        }
        let entries: Vec<_> = processed
            .iter()
            .map(|p| p.entry_at(&destination.join(&p.source.relative)))
            .collect();
        let path = index_path(destination, options.seal.format);
        match write_seal(&path, options.seal.format, &entries) {
            Ok(()) => {
                report.seals.push(path.clone());
                event(
                    &mut emit,
                    &format!("Wrote seal file to '{}'", path.display()),
                    Importance::Result,
                    false,
                );
            }
            Err(err) => {
                report.statistics.errors += 1;
                rollback(
                    &processed,
                    index,
                    destination,
                    &mut report.statistics,
                    &mut emit,
                );
                event(&mut emit, &err.to_string(), Importance::Error, true);
            }
        }
    }

    fill_stats(
        &processed,
        &mut report.statistics,
        started,
        &options.seal.common.cancellation,
    );
    if options.verify_after_copy && !report.seals.is_empty() && !report.statistics.cancelled {
        let verify_result = verify(
            &report.seals,
            VerifyOptions {
                input_streams: options.seal.common.input_streams,
                cancellation: options.seal.common.cancellation.clone(),
            },
            &mut emit,
        );
        if verify_result.is_err() {
            report.statistics.errors += 1;
        }
    }
    finish_event(&mut emit, "SEAL", &report.statistics);
    if report.statistics.errors > 0 {
        Err(GodiError("sealed-copy failed".into()))
    } else {
        Ok(report)
    }
}

pub fn verify(
    seals: &[PathBuf],
    options: VerifyOptions,
    mut emit: impl FnMut(&Event),
) -> Result<OperationReport, GodiError> {
    if seals.is_empty() {
        return Err(GodiError("Please provide at least one seal file".into()));
    }
    let started = Instant::now();
    let mut tasks = Vec::new();
    let mut decode_errors = 0;
    for seal in seals {
        let directory = seal.parent().unwrap_or_else(|| Path::new("."));
        match read_seal(seal) {
            Ok(entries) => {
                for entry in entries {
                    let path = directory.join(&entry.relative_path);
                    let metadata = fs::symlink_metadata(&path);
                    tasks.push((
                        entry,
                        metadata.map(|m| SourceFile {
                            root: directory.to_path_buf(),
                            path,
                            relative: PathBuf::new(),
                            mode: mode(&m),
                            size: m.len() as i64,
                            symlink: m.file_type().is_symlink(),
                        }),
                    ));
                }
            }
            Err(err) => {
                decode_errors += 1;
                event(
                    &mut emit,
                    &format!("SEAL MISMATCH: '{}' - {err}", seal.display()),
                    Importance::Error,
                    true,
                );
            }
        }
    }

    let pool = rayon::ThreadPoolBuilder::new()
        .num_threads(options.input_streams.get())
        .build()
        .map_err(|err| GodiError(err.to_string()))?;
    let cancellation = options.cancellation.clone();
    let checked: Vec<_> = pool.install(|| {
        tasks
            .into_par_iter()
            .map(|(expected, source)| match source {
                Ok(source) => {
                    hash_file(&source, &[], &cancellation, 1).map(|actual| (expected, Ok(actual)))
                }
                Err(err) => Ok((expected, Err(err))),
            })
            .collect::<Result<Vec<_>, io::Error>>()
    })?;

    let mut report = OperationReport::default();
    report.statistics.errors = decode_errors;
    for (expected, actual) in checked {
        match actual {
            Err(err) => {
                report.statistics.errors += 1;
                event(
                    &mut emit,
                    &format!(
                        "MISSING MISMATCH: {}: {err}",
                        expected.relative_path.display()
                    ),
                    Importance::Error,
                    true,
                );
            }
            Ok(actual) => {
                report.statistics.files_read += 1;
                report.statistics.bytes_read += actual.bytes;
                if actual.source.size != expected.size {
                    report.statistics.errors += 1;
                    event(
                        &mut emit,
                        &format!(
                            "SIZE MISMATCH: {} sealed with size {}B, got size {}B",
                            actual.source.path.display(),
                            expected.size,
                            actual.source.size
                        ),
                        Importance::Error,
                        true,
                    );
                } else if (!expected.sha1.is_empty() && expected.sha1 != actual.sha1)
                    || (!expected.md5.is_empty() && expected.md5 != actual.md5)
                {
                    report.statistics.errors += 1;
                    event(
                        &mut emit,
                        &format!(
                            "HASH MISMATCH: {} flipped at least one bit",
                            actual.source.path.display()
                        ),
                        Importance::Error,
                        true,
                    );
                } else {
                    event(
                        &mut emit,
                        &format!("OK: {}", actual.source.path.display()),
                        Importance::Info,
                        false,
                    );
                }
            }
        }
    }
    report.statistics.cancelled = options.cancellation.is_cancelled();
    report.statistics.elapsed = started.elapsed();
    finish_event(&mut emit, "VERIFY", &report.statistics);
    if report.statistics.errors > 0 || report.statistics.cancelled {
        Err(GodiError("verification failed".into()))
    } else {
        Ok(report)
    }
}

impl Processed {
    fn entry_at(&self, path: &Path) -> SealEntry {
        SealEntry {
            path: path.to_path_buf(),
            relative_path: self.source.relative.clone(),
            mode: self.source.mode,
            size: self.source.size,
            sha1: self.sha1.clone(),
            md5: self.md5.clone(),
        }
    }
}

fn discover(
    sources: &[PathBuf],
    filters: &[FileFilter],
) -> Result<(Vec<PathBuf>, Vec<SourceFile>, u64), GodiError> {
    if sources.is_empty() {
        return Err(GodiError("Please provide at least one source".into()));
    }
    let source_paths = parse_sources(sources)?;
    let mut roots = Vec::new();
    let mut files = Vec::new();
    let mut skipped = 0;
    for source in &source_paths {
        let metadata = fs::symlink_metadata(source)?;
        if metadata.is_dir() {
            roots.push(source.clone());
            walk(source, source, filters, &mut files, &mut skipped)?;
        } else if metadata.is_file() || metadata.file_type().is_symlink() {
            let root = source.parent().unwrap_or_else(|| Path::new("."));
            if !roots.iter().any(|item| item == root) {
                roots.push(root.to_path_buf());
            }
            files.push(source_file(
                root,
                source,
                source.file_name().unwrap_or_default().into(),
                &metadata,
            ));
        } else {
            return Err(GodiError(format!(
                "'{}' is not a regular file or directory",
                source.display()
            )));
        }
    }
    Ok((roots, files, skipped))
}

fn walk(
    root: &Path,
    directory: &Path,
    filters: &[FileFilter],
    files: &mut Vec<SourceFile>,
    skipped: &mut u64,
) -> Result<(), GodiError> {
    let mut directories = Vec::new();
    for item in fs::read_dir(directory)? {
        let item = item?;
        let path = item.path();
        let metadata = fs::symlink_metadata(&path)?;
        if filters
            .iter()
            .any(|f| f.matches(&item.file_name(), &metadata))
        {
            *skipped += 1;
            continue;
        }
        if metadata.is_dir() {
            directories.push(path);
        } else {
            let relative = path
                .strip_prefix(root)
                .map_err(|err| GodiError(err.to_string()))?
                .to_path_buf();
            files.push(source_file(root, &path, relative, &metadata));
        }
    }
    for child in directories {
        walk(root, &child, filters, files, skipped)?;
    }
    Ok(())
}

fn source_file(root: &Path, path: &Path, relative: PathBuf, metadata: &fs::Metadata) -> SourceFile {
    SourceFile {
        root: root.to_path_buf(),
        path: path.to_path_buf(),
        relative,
        mode: mode(metadata),
        size: metadata.len() as i64,
        symlink: metadata.file_type().is_symlink(),
    }
}

fn process(
    files: Vec<SourceFile>,
    destinations: &[PathBuf],
    options: &CommonOptions,
    output_streams: usize,
) -> Vec<Processed> {
    let threads = options
        .input_streams
        .get()
        .saturating_mul(3)
        .saturating_add(output_streams.saturating_sub(1));
    let pool = rayon::ThreadPoolBuilder::new()
        .num_threads(threads)
        .build()
        .expect("valid thread count");
    let failed: Arc<Vec<AtomicBool>> = Arc::new(
        (0..destinations.len().max(1))
            .map(|_| AtomicBool::new(false))
            .collect(),
    );
    pool.install(|| {
        files
            .into_par_iter()
            .map(|file| {
                let active: Vec<_> = destinations
                    .iter()
                    .enumerate()
                    .map(|(i, root)| {
                        if failed[i].load(Ordering::Acquire) {
                            None
                        } else {
                            Some(root.join(&file.relative))
                        }
                    })
                    .collect();
                match hash_file(&file, &active, &options.cancellation, output_streams) {
                    Ok(result) => {
                        for (i, error) in result.errors.iter().enumerate() {
                            if error.is_some() {
                                failed[i].store(true, Ordering::Release);
                            }
                        }
                        result
                    }
                    Err(err) => Processed {
                        source: file,
                        sha1: Vec::new(),
                        md5: Vec::new(),
                        copied: vec![None; destinations.len().max(1)],
                        errors: vec![Some(err.to_string()); destinations.len().max(1)],
                        bytes: 0,
                    },
                }
            })
            .collect()
    })
}

fn hash_file(
    source: &SourceFile,
    output_paths: &[Option<PathBuf>],
    cancellation: &CancellationToken,
    _output_streams: usize,
) -> io::Result<Processed> {
    let mut sha1 = Sha1::new();
    let mut md5 = Md5::new();
    let mut outputs: Vec<Option<File>> = Vec::with_capacity(output_paths.len());
    let mut copied = vec![None; output_paths.len().max(1)];
    let mut errors = vec![None; output_paths.len().max(1)];

    for (index, path) in output_paths.iter().enumerate() {
        let Some(path) = path else {
            outputs.push(None);
            continue;
        };
        if let Some(parent) = path.parent()
            && let Err(err) = fs::create_dir_all(parent)
        {
            errors[index] = Some(err.to_string());
            outputs.push(None);
            continue;
        }
        if source.symlink {
            outputs.push(None);
        } else {
            match OpenOptions::new().write(true).create_new(true).open(path) {
                Ok(file) => {
                    set_mode(path, source.mode)?;
                    copied[index] = Some(path.clone());
                    outputs.push(Some(file));
                }
                Err(err) => {
                    errors[index] = Some(err.to_string());
                    outputs.push(None);
                }
            }
        }
    }

    if source.symlink {
        let target = fs::read_link(&source.path)?;
        let bytes = path_os_bytes(&target);
        sha1.update(&bytes);
        md5.update(&bytes);
        for (index, path) in output_paths.iter().enumerate() {
            let Some(path) = path else { continue };
            if errors[index].is_none() {
                match create_symlink(&target, path) {
                    Ok(()) => copied[index] = Some(path.clone()),
                    Err(err) => errors[index] = Some(err.to_string()),
                }
            }
        }
        return Ok(Processed {
            source: source.clone(),
            sha1: sha1.finalize().to_vec(),
            md5: md5.finalize().to_vec(),
            copied,
            errors,
            bytes: bytes.len() as u64,
        });
    }

    let mut input = File::open(&source.path)?;
    let mut buffer = vec![0; 512 * 1024];
    let mut total = 0;
    loop {
        if cancellation.is_cancelled() {
            return Err(io::Error::new(
                io::ErrorKind::Interrupted,
                format!("Reading of '{}' cancelled", source.path.display()),
            ));
        }
        let count = input.read(&mut buffer)?;
        if count == 0 {
            break;
        }
        let chunk = &buffer[..count];
        let (_, write_errors) = rayon::join(
            || rayon::join(|| sha1.update(chunk), || md5.update(chunk)),
            || {
                outputs
                    .par_iter_mut()
                    .enumerate()
                    .filter_map(|(index, output)| {
                        if let Some(file) = output.as_mut() {
                            if let Err(err) = file.write_all(chunk) {
                                let message = err.to_string();
                                *output = None;
                                return Some((index, message));
                            }
                        }
                        None
                    })
                    .collect::<Vec<_>>()
            },
        );
        for (index, error) in write_errors {
            errors[index] = Some(error);
        }
        total += count as u64;
    }
    if total != source.size as u64 {
        return Err(io::Error::new(
            io::ErrorKind::UnexpectedEof,
            format!(
                "Filesize of '{}' reported as {}, yet {} bytes were read",
                source.path.display(),
                source.size,
                total
            ),
        ));
    }
    for (index, output) in outputs.iter_mut().enumerate() {
        if let Some(file) = output
            && let Err(err) = file.flush()
        {
            errors[index] = Some(err.to_string());
        }
    }
    Ok(Processed {
        source: source.clone(),
        sha1: sha1.finalize().to_vec(),
        md5: md5.finalize().to_vec(),
        copied,
        errors,
        bytes: total,
    })
}

fn parse_sources(items: &[PathBuf]) -> Result<Vec<PathBuf>, GodiError> {
    let mut paths = Vec::new();
    for item in items {
        if !item.exists() {
            return Err(GodiError(format!("'{}' does not exist", item.display())));
        }
        let absolute = if item.is_absolute() {
            clean(item)
        } else {
            clean(&std::env::current_dir()?.join(item))
        };
        if !paths.contains(&absolute) {
            paths.push(absolute);
        }
    }
    paths.sort_by_key(|path| path.components().count());
    let mut result = Vec::new();
    for path in paths {
        if !result.iter().any(|root: &PathBuf| path.starts_with(root)) {
            result.push(path);
        }
    }
    Ok(result)
}

fn parse_destinations(
    destinations: &[PathBuf],
    sources: &[PathBuf],
) -> Result<Vec<PathBuf>, GodiError> {
    let destinations = parse_sources(destinations)?;
    for destination in &destinations {
        if !destination.is_dir() {
            return Err(GodiError(format!(
                "'{}' is not a destination directory",
                destination.display()
            )));
        }
        for source in sources {
            let source = if source.is_absolute() {
                clean(source)
            } else {
                clean(&std::env::current_dir()?.join(source))
            };
            if destination.starts_with(&source) {
                return Err(GodiError(format!(
                    "Cannot copy '{}' into itself at '{}'",
                    source.display(),
                    destination.display()
                )));
            }
        }
    }
    Ok(destinations)
}

fn clean(path: &Path) -> PathBuf {
    let mut result = PathBuf::new();
    for component in path.components() {
        match component {
            Component::CurDir => {}
            Component::ParentDir => {
                result.pop();
            }
            other => result.push(other.as_os_str()),
        }
    }
    result
}

fn index_path(root: &Path, format: SealFormat) -> PathBuf {
    root.join(format!(
        "godi_{}.{}",
        Local::now().format("%Y-%m-%d_%H%M%S"),
        format.extension()
    ))
}

fn rollback(
    processed: &[Processed],
    destination: usize,
    root: &Path,
    stats: &mut Statistics,
    emit: &mut impl FnMut(&Event),
) {
    let mut paths: Vec<_> = processed
        .iter()
        .filter_map(|p| p.copied.get(destination).and_then(Clone::clone))
        .collect();
    paths.sort_by_key(|path| std::cmp::Reverse(path.components().count()));
    for path in paths {
        if fs::remove_file(&path).is_ok() {
            stats.undone += 1;
            event(
                emit,
                &format!("Removed '{}'", path.display()),
                Importance::Info,
                false,
            );
            let mut parent = path.parent();
            while let Some(directory) = parent {
                if directory == root || fs::remove_dir(directory).is_err() {
                    break;
                }
                parent = directory.parent();
            }
        }
    }
}

fn fill_stats(
    processed: &[Processed],
    stats: &mut Statistics,
    started: Instant,
    cancellation: &CancellationToken,
) {
    stats.files_read = processed.len() as u64;
    stats.bytes_read = processed.iter().map(|p| p.bytes).sum();
    stats.files_written = processed
        .iter()
        .flat_map(|p| &p.copied)
        .filter(|p| p.is_some())
        .count() as u64;
    stats.bytes_written = processed
        .iter()
        .map(|p| p.bytes * p.copied.iter().filter(|v| v.is_some()).count() as u64)
        .sum();
    stats.cancelled = cancellation.is_cancelled();
    stats.elapsed = started.elapsed();
}

fn finish_event(emit: &mut impl FnMut(&Event), operation: &str, stats: &Statistics) {
    event(
        emit,
        &format!(
            "{operation} {}: {} files, {} bytes in {:.2?}{}",
            if stats.errors == 0 && !stats.cancelled {
                "SUCCESS"
            } else {
                "FAILED"
            },
            stats.files_read,
            stats.bytes_read,
            stats.elapsed,
            if stats.cancelled { " (cancelled)" } else { "" }
        ),
        Importance::Result,
        stats.errors > 0 || stats.cancelled,
    );
}

fn event(emit: &mut impl FnMut(&Event), message: &str, importance: Importance, is_error: bool) {
    emit(&Event {
        message: message.into(),
        importance,
        error: is_error.then(|| message.into()),
    });
}

fn is_seal_name(name: &str) -> bool {
    let Some(rest) = name.strip_prefix("godi_") else {
        return false;
    };
    let bytes = rest.as_bytes();
    bytes.len() >= 22
        && bytes[4] == b'-'
        && bytes[7] == b'-'
        && bytes[10] == b'_'
        && bytes[17] == b'.'
        && bytes[..17]
            .iter()
            .enumerate()
            .all(|(i, b)| matches!(i, 4 | 7 | 10) || b.is_ascii_digit())
}

#[cfg(unix)]
fn mode(metadata: &fs::Metadata) -> u32 {
    use std::os::unix::fs::MetadataExt;
    metadata.mode()
}

#[cfg(not(unix))]
fn mode(metadata: &fs::Metadata) -> u32 {
    if metadata.permissions().readonly() {
        0o444
    } else {
        0o666
    }
}

#[cfg(unix)]
fn set_mode(path: &Path, mode: u32) -> io::Result<()> {
    use std::os::unix::fs::PermissionsExt;
    fs::set_permissions(path, fs::Permissions::from_mode(mode))
}

#[cfg(not(unix))]
fn set_mode(_path: &Path, _mode: u32) -> io::Result<()> {
    Ok(())
}

#[cfg(unix)]
fn path_os_bytes(path: &Path) -> Vec<u8> {
    use std::os::unix::ffi::OsStrExt;
    path.as_os_str().as_bytes().to_vec()
}

#[cfg(not(unix))]
fn path_os_bytes(path: &Path) -> Vec<u8> {
    path.to_string_lossy().as_bytes().to_vec()
}

#[cfg(unix)]
fn create_symlink(target: &Path, path: &Path) -> io::Result<()> {
    std::os::unix::fs::symlink(target, path)
}

#[cfg(windows)]
fn create_symlink(target: &Path, path: &Path) -> io::Result<()> {
    if target.is_dir() {
        std::os::windows::fs::symlink_dir(target, path)
    } else {
        std::os::windows::fs::symlink_file(target, path)
    }
}
