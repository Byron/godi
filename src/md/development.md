# Development

`godi` is a Rust 2024 crate. Go is not used or required.

```sh
cargo test --all-features
cargo fmt --all -- --check
cargo clippy --all-targets --all-features -- -D warnings
```

The library in `src/engine.rs` discovers inputs, hashes SHA-1 and MD5 in
parallel, fans each buffer out to copy destinations, and rolls back files it
created when a destination fails. `src/codec.rs` reads and writes MHL and the
gzip-compressed gob stream used by godi v1.1. The optional `web` feature embeds
the dependency-free frontend and HTTP API.

Release builds use:

```sh
cargo build --release --all-features
```
