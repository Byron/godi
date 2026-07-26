use crate::{
    CancellationToken, CommonOptions, CopyOptions, Event, FileFilter, SealFormat, SealOptions,
    VerifyOptions, seal, sealed_copy, verify,
};
use axum::{
    Json, Router,
    body::{Body, to_bytes},
    extract::{
        ConnectInfo, Query, State,
        ws::{Message, WebSocket, WebSocketUpgrade},
    },
    http::{HeaderValue, Method, Request, StatusCode, header},
    response::{IntoResponse, Response},
    routing::{any, get},
};
use serde::{Deserialize, Serialize};
use std::{
    net::{IpAddr, SocketAddr},
    num::NonZeroUsize,
    path::{Path, PathBuf},
    str::FromStr,
    sync::{Arc, Mutex},
    time::{Duration, Instant},
};
use tokio::sync::broadcast;

const INDEX: &[u8] = include_bytes!("../web/index.html");
const APP_JS: &[u8] = include_bytes!("../web/app.js");
const APP_CSS: &[u8] = include_bytes!("../web/app.css");

#[derive(Clone)]
struct App {
    inner: Arc<Mutex<Inner>>,
    events: broadcast::Sender<String>,
}

struct Inner {
    state: WebState,
    owner: Option<IpAddr>,
    touched: Instant,
    cancellation: Option<CancellationToken>,
}

#[derive(Clone, Deserialize, Serialize)]
#[serde(default)]
#[serde(rename_all = "camelCase")]
struct WebState {
    mode: String,
    verbosity: String,
    spid: usize,
    spod: usize,
    fep: Vec<String>,
    sources: Vec<PathBuf>,
    destinations: Vec<PathBuf>,
    verify: String,
    format: String,
    socket_url: String,
    #[serde(rename = "status")]
    is_running: bool,
    execution_error: String,
}

impl Default for WebState {
    fn default() -> Self {
        Self {
            mode: "verify".into(),
            verbosity: "error".into(),
            spid: 1,
            spod: 1,
            fep: vec!["VOLATILE".into()],
            sources: Vec::new(),
            destinations: Vec::new(),
            verify: String::new(),
            format: "gob".into(),
            socket_url: "/api/v1/websocket".into(),
            is_running: false,
            execution_error: String::new(),
        }
    }
}

#[derive(Serialize)]
struct Defaults {
    modes: [&'static str; 3],
    verbosities: [&'static str; 2],
    feps: [&'static str; 4],
    formats: [&'static str; 2],
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
struct SocketEvent<'a> {
    message: &'a str,
    error: &'a str,
    importance: u8,
    client_id: &'a str,
    state: u8,
}

#[derive(Deserialize)]
struct DirQuery {
    path: PathBuf,
    #[serde(rename = "type")]
    kind: String,
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
struct DirItem {
    item: String,
    path: PathBuf,
    is_dir: bool,
}

pub async fn serve(address: &str, show: bool) -> Result<(), String> {
    let listener = tokio::net::TcpListener::bind(address)
        .await
        .map_err(|err| err.to_string())?;
    let (events, _) = broadcast::channel(128);
    let app = App {
        inner: Arc::new(Mutex::new(Inner {
            state: WebState::default(),
            owner: None,
            touched: Instant::now(),
            cancellation: None,
        })),
        events,
    };
    let router = Router::new()
        .route("/api/v1/state", any(state_handler))
        .route("/api/v1/dirlist", get(dirlist))
        .route("/api/v1/websocket", get(websocket))
        .fallback(get(static_file))
        .with_state(app);
    let url = format!("http://{address}");
    println!("About to listen on {url}");
    println!("Hit CTRL+C to close");
    if show {
        let _ = webbrowser::open(&url);
    }
    axum::serve(
        listener,
        router.into_make_service_with_connect_info::<SocketAddr>(),
    )
    .with_graceful_shutdown(async {
        let _ = tokio::signal::ctrl_c().await;
    })
    .await
    .map_err(|err| err.to_string())
}

async fn state_handler(
    State(app): State<App>,
    ConnectInfo(remote): ConnectInfo<SocketAddr>,
    request: Request<Body>,
) -> Response {
    let method = request.method().clone();
    let client_id = request
        .headers()
        .get("Client-ID")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("")
        .to_owned();

    if method.as_str() == "DEFAULTS" {
        return Json(Defaults {
            modes: ["verify", "seal", "sealed-copy"],
            verbosities: ["info", "error"],
            feps: ["HIDDEN", "SEALS", "SYMLINK", "VOLATILE"],
            formats: ["gob", "mhl"],
        })
        .into_response();
    }

    if method == Method::GET {
        let mut inner = app.inner.lock().expect("web state lock");
        expire_owner(&mut inner);
        let writable = inner.owner.is_none() || inner.owner == Some(remote.ip());
        let mut response = Json(inner.state.clone()).into_response();
        response.headers_mut().insert(
            "X-is-RW",
            HeaderValue::from_static(if writable { "true" } else { "false" }),
        );
        return response;
    }

    if method == Method::DELETE {
        let inner = app.inner.lock().expect("web state lock");
        return if let Some(token) = inner.cancellation.as_ref() {
            token.cancel();
            StatusCode::OK.into_response()
        } else {
            (
                StatusCode::PRECONDITION_FAILED,
                "No operation is currently in progress",
            )
                .into_response()
        };
    }

    if method != Method::PUT && method != Method::POST {
        return (StatusCode::BAD_REQUEST, "Unsupported Method").into_response();
    }

    let body = match to_bytes(request.into_body(), 1024 * 1024).await {
        Ok(body) => body,
        Err(err) => return (StatusCode::BAD_REQUEST, err.to_string()).into_response(),
    };
    let update = if body.is_empty() {
        None
    } else {
        match serde_json::from_slice::<WebState>(&body) {
            Ok(state) => Some(state),
            Err(err) => return (StatusCode::PRECONDITION_FAILED, err.to_string()).into_response(),
        }
    };

    let current = {
        let mut inner = app.inner.lock().expect("web state lock");
        expire_owner(&mut inner);
        if inner.state.is_running {
            return (
                StatusCode::PRECONDITION_FAILED,
                "Operation already in progress",
            )
                .into_response();
        }
        if inner.owner.is_some() && inner.owner != Some(remote.ip()) {
            return (StatusCode::UNAUTHORIZED, "Changing values not permitted").into_response();
        }
        if let Some(state) = update {
            if let Err(err) = validate(&state, false) {
                return (StatusCode::PRECONDITION_FAILED, err).into_response();
            }
            inner.state = state;
            inner.state.socket_url = "/api/v1/websocket".into();
        }
        inner.owner = Some(remote.ip());
        inner.touched = Instant::now();
        inner.state.clone()
    };

    if method == Method::PUT {
        broadcast_state(&app, 0, "", "", 0, &client_id);
        return Json(current).into_response();
    }
    if let Err(err) = validate(&current, true) {
        return (StatusCode::PRECONDITION_FAILED, err).into_response();
    }
    start_operation(app.clone(), current);
    Json(app.inner.lock().expect("web state lock").state.clone()).into_response()
}

fn start_operation(app: App, state: WebState) {
    let cancellation = CancellationToken::default();
    {
        let mut inner = app.inner.lock().expect("web state lock");
        inner.state.is_running = true;
        inner.state.execution_error.clear();
        inner.cancellation = Some(cancellation.clone());
    }
    broadcast_state(&app, 2, "", "", 0, "");
    tokio::task::spawn_blocking(move || {
        let streams = NonZeroUsize::new(state.spid).unwrap_or(NonZeroUsize::MIN);
        let filters = state
            .fep
            .iter()
            .filter_map(|value| FileFilter::parse(value).ok())
            .collect();
        let common = CommonOptions {
            input_streams: streams,
            filters,
            cancellation: cancellation.clone(),
        };
        let format = SealFormat::from_str(&state.format).unwrap_or_default();
        let mut emit = |event: &Event| {
            broadcast_state(
                &app,
                1,
                &event.message,
                event.error.as_deref().unwrap_or(""),
                event.importance as u8,
                "",
            );
        };
        let result = match state.mode.as_str() {
            "seal" => seal(&state.sources, SealOptions { common, format }, &mut emit),
            "sealed-copy" => sealed_copy(
                &state.sources,
                &state.destinations,
                CopyOptions {
                    seal: SealOptions { common, format },
                    output_streams: NonZeroUsize::new(state.spod).unwrap_or(NonZeroUsize::MIN),
                    verify_after_copy: !state.verify.is_empty(),
                },
                &mut emit,
            ),
            _ => verify(
                &state.sources,
                VerifyOptions {
                    input_streams: streams,
                    cancellation,
                },
                &mut emit,
            ),
        };
        {
            let mut inner = app.inner.lock().expect("web state lock");
            inner.state.is_running = false;
            inner.state.execution_error = result.err().map(|e| e.to_string()).unwrap_or_default();
            inner.cancellation = None;
            inner.owner = None;
        }
        broadcast_state(&app, 3, "", "", 0, "");
    });
}

async fn websocket(State(app): State<App>, upgrade: WebSocketUpgrade) -> impl IntoResponse {
    upgrade.on_upgrade(move |socket| websocket_client(socket, app.events.subscribe()))
}

async fn websocket_client(mut socket: WebSocket, mut events: broadcast::Receiver<String>) {
    while let Ok(event) = events.recv().await {
        if socket.send(Message::Text(event.into())).await.is_err() {
            break;
        }
    }
}

async fn dirlist(State(app): State<App>, Query(query): Query<DirQuery>) -> Response {
    if query.kind != "all" && query.kind != "sealOnly" {
        return (StatusCode::BAD_REQUEST, "Invalid request type").into_response();
    }
    let state = app.inner.lock().expect("web state lock").state.clone();
    let path = query.path;
    let (directory, prefix) = if path.to_string_lossy().ends_with(std::path::MAIN_SEPARATOR) {
        (path.clone(), String::new())
    } else {
        (
            path.parent()
                .unwrap_or_else(|| Path::new("."))
                .to_path_buf(),
            path.file_name()
                .unwrap_or_default()
                .to_string_lossy()
                .to_lowercase(),
        )
    };
    let mut items = Vec::new();
    let entries = match std::fs::read_dir(&directory) {
        Ok(entries) => entries,
        Err(err) => return (StatusCode::BAD_REQUEST, err.to_string()).into_response(),
    };
    for entry in entries.flatten() {
        let Ok(metadata) = std::fs::symlink_metadata(entry.path()) else {
            continue;
        };
        let name = entry.file_name().to_string_lossy().into_owned();
        if !name.to_lowercase().contains(&prefix) {
            continue;
        }
        if query.kind == "sealOnly"
            && !metadata.is_dir()
            && !(name.starts_with("godi_") && (name.ends_with(".gobz") || name.ends_with(".mhl")))
        {
            continue;
        }
        if query.kind == "all"
            && state
                .fep
                .iter()
                .filter_map(|v| FileFilter::parse(v).ok())
                .filter(|filter| !matches!(filter, FileFilter::Seals))
                .any(|filter| filter.matches(&entry.file_name(), &metadata))
        {
            continue;
        }
        items.push(DirItem {
            item: name,
            path: entry.path(),
            is_dir: metadata.is_dir(),
        });
    }
    Json(items).into_response()
}

async fn static_file(request: Request<Body>) -> Response {
    let (content_type, bytes) = match request.uri().path() {
        "/" | "/index.html" => ("text/html; charset=utf-8", INDEX),
        "/app.js" => ("text/javascript; charset=utf-8", APP_JS),
        "/app.css" => ("text/css; charset=utf-8", APP_CSS),
        _ => return StatusCode::NOT_FOUND.into_response(),
    };
    (
        [(header::CONTENT_TYPE, HeaderValue::from_static(content_type))],
        bytes,
    )
        .into_response()
}

fn validate(state: &WebState, complete: bool) -> Result<(), String> {
    if !matches!(state.mode.as_str(), "verify" | "seal" | "sealed-copy") {
        return Err(format!("Invalid mode: '{}'", state.mode));
    }
    if !matches!(state.verbosity.as_str(), "info" | "error") {
        return Err(format!("Invalid verbosity: '{}'", state.verbosity));
    }
    if state.spid == 0 || (state.mode == "sealed-copy" && state.spod == 0) {
        return Err("stream counts must be larger than 0".into());
    }
    if state.mode != "verify" {
        SealFormat::from_str(&state.format)?;
    }
    for filter in &state.fep {
        FileFilter::parse(filter).map_err(|err| err.to_string())?;
    }
    if complete && state.sources.is_empty() {
        return Err("Didn't provide a single source".into());
    }
    if complete && state.mode == "sealed-copy" && state.destinations.is_empty() {
        return Err("Need to provide at least one destination".into());
    }
    Ok(())
}

fn expire_owner(inner: &mut Inner) {
    if inner.touched.elapsed() >= Duration::from_secs(5 * 60) && !inner.state.is_running {
        inner.owner = None;
    }
}

fn broadcast_state(
    app: &App,
    state: u8,
    message: &str,
    error: &str,
    importance: u8,
    client_id: &str,
) {
    if let Ok(json) = serde_json::to_string(&SocketEvent {
        message,
        error,
        importance,
        client_id,
        state,
    }) {
        let _ = app.events.send(json);
    }
}
