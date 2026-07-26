//! Library interface for sealing, copying, and verifying immutable data.

mod codec;
mod engine;

pub use codec::{SealEntry, SealFormat, read_seal, write_seal};
pub use engine::{
    CancellationToken, CommonOptions, CopyOptions, Event, FileFilter, GodiError, Importance,
    OperationReport, SealOptions, Statistics, VerifyOptions, seal, sealed_copy, verify,
};

#[cfg(feature = "web")]
pub mod web;
