# godi

`godi` seals immutable data, verifies it later, and can copy it to multiple
destinations while hashing each input only once.

## Build

Rust 1.85 or newer is required. Go is neither required nor used.

```sh
cargo build --release
cargo test --all-features
```

The optional self-contained web frontend is included with:

```sh
cargo build --release --features web
cargo run --features web -- web
```

## Usage

```sh
# Seal files or directories. The default format is compressed .gobz.
godi seal ~/valuable-data

# Verify an existing .gobz or MHL seal.
godi verify ~/valuable-data/godi_2026-07-26_120000.gobz

# Copy once to one or more existing destination directories.
godi sealed-copy ~/valuable-data -- /backup/one /backup/two

# Re-read and verify successful copies.
godi sealed-copy --verify ~/valuable-data /backup/one
```

Global controls retain the v1.1 interface:

- `--streams-per-input-device`, alias `--spid` or legacy `-spid`
- `--verbosity=statistics|info|warn|error|result|off`
- `--file-exclude-patterns=VOLATILE,HIDDEN,SYMLINK,SEALS,*.tmp`

`sealed-copy` never overwrites existing files. If one destination fails, only
files created in that destination are removed; unaffected destinations finish.
Pressing Ctrl-C triggers the same cleanup.

## Seal compatibility

- `gob` writes gzip-compressed, signed `.gobz` files compatible with godi v1.1.
- `mhl` reads and writes Media Hash List 1.0 XML, including unsigned third-party
  MHL files.
- Every generated seal carries a SHA-1 signature over its entries.

The crate also exposes `seal`, `sealed_copy`, `verify`, `read_seal`, and
`write_seal` for embedding.

## License

LGPL-3.0-only. See [LICENSE.md](LICENSE.md).
