use flate2::{Compression, read::GzDecoder, write::GzEncoder};
use quick_xml::{de::from_reader, se::to_string};
use serde::{Deserialize, Serialize};
use sha1::{Digest, Sha1};
use std::{
    fmt,
    fs::{File, OpenOptions},
    io::{self, BufReader, BufWriter, Read, Write},
    path::{Path, PathBuf},
};

#[derive(Clone, Copy, Debug, Default, Eq, PartialEq)]
pub enum SealFormat {
    #[default]
    Gob,
    Mhl,
}

impl SealFormat {
    pub fn extension(self) -> &'static str {
        match self {
            Self::Gob => "gobz",
            Self::Mhl => "mhl",
        }
    }
}

impl std::str::FromStr for SealFormat {
    type Err = String;

    fn from_str(value: &str) -> Result<Self, Self::Err> {
        match value {
            "gob" => Ok(Self::Gob),
            "mhl" => Ok(Self::Mhl),
            _ => Err(format!(
                "invalid seal format '{value}', expected gob or mhl"
            )),
        }
    }
}

impl fmt::Display for SealFormat {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(match self {
            Self::Gob => "gob",
            Self::Mhl => "mhl",
        })
    }
}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct SealEntry {
    pub path: PathBuf,
    pub relative_path: PathBuf,
    pub mode: u32,
    pub size: i64,
    pub sha1: Vec<u8>,
    pub md5: Vec<u8>,
}

pub fn write_seal(path: &Path, format: SealFormat, entries: &[SealEntry]) -> io::Result<()> {
    let file = OpenOptions::new().write(true).create_new(true).open(path)?;
    let result = match format {
        SealFormat::Gob => write_gob(file, entries),
        SealFormat::Mhl => write_mhl(file, entries),
    };
    if result.is_err() {
        let _ = std::fs::remove_file(path);
    }
    result
}

pub fn read_seal(path: &Path) -> io::Result<Vec<SealEntry>> {
    let file = File::open(path)?;
    match path.extension().and_then(|v| v.to_str()) {
        Some("gobz") => read_gob(file),
        Some("mhl") => read_mhl(file),
        _ => Err(io::Error::new(
            io::ErrorKind::InvalidInput,
            format!("unknown seal file format: '{}'", path.display()),
        )),
    }
}

fn signature(entries: &[SealEntry], mhl: bool) -> Vec<u8> {
    let mut hash = Sha1::new();
    for entry in entries {
        hash.update(path_bytes(&entry.relative_path));
        hash.update(if mhl {
            path_bytes(&entry.relative_path)
        } else {
            path_bytes(&entry.path)
        });
        hash.update(&entry.sha1);
        hash.update(&entry.md5);
    }
    hash.finalize().to_vec()
}

fn write_gob(file: File, entries: &[SealEntry]) -> io::Result<()> {
    let mut out = GzEncoder::new(BufWriter::new(file), Compression::best());
    frame(&mut out, |v| {
        signed(v, 2);
        unsigned(v, 0);
        signed(v, 1);
    })?;
    write_file_info_type(&mut out)?;
    for entry in entries {
        frame(&mut out, |v| {
            signed(v, 65);
            let mut previous = -1i32;
            entry_field_bytes(v, &mut previous, 0, &path_bytes(&entry.path));
            entry_field_bytes(v, &mut previous, 1, &path_bytes(&entry.relative_path));
            if entry.mode != 0 {
                entry_field(v, &mut previous, 2);
                unsigned(v, entry.mode as u64);
            }
            if entry.size != 0 {
                entry_field(v, &mut previous, 3);
                signed(v, entry.size);
            }
            entry_field_bytes(v, &mut previous, 4, &entry.sha1);
            entry_field_bytes(v, &mut previous, 5, &entry.md5);
            unsigned(v, 0);
        })?;
    }
    frame(&mut out, |v| {
        signed(v, 1);
        unsigned(v, 0);
        unsigned(v, 1);
    })?;
    let seal_signature = signature(entries, false);
    frame(&mut out, |v| {
        signed(v, 5);
        unsigned(v, 0);
        bytes(v, &seal_signature);
    })?;
    out.finish()?.flush()
}

fn read_gob(file: File) -> io::Result<Vec<SealEntry>> {
    let mut input = GzDecoder::new(BufReader::new(file));
    let mut entries = Vec::new();
    let mut file_type = None;
    let mut saw_version = false;
    let mut stored_signature = None;
    while let Some(message) = read_frame(&mut input)? {
        let mut message = message.as_slice();
        let type_id = read_signed(&mut message)?;
        if type_id < 0 {
            file_type.get_or_insert(-type_id);
            continue;
        }
        match type_id {
            2 if !saw_version => {
                if read_unsigned(&mut message)? != 0 {
                    return invalid("invalid gobz version wrapper");
                }
                if read_signed(&mut message)? != 1 {
                    return invalid("unsupported gobz version");
                }
                saw_version = true;
            }
            id if Some(id) == file_type => entries.push(read_file_info(&mut message)?),
            1 => {
                let _ = read_unsigned(&mut message)?;
                let _ = read_unsigned(&mut message)?;
            }
            5 => {
                let _ = read_unsigned(&mut message)?;
                stored_signature = Some(read_bytes(&mut message)?);
            }
            _ => return invalid("unexpected value in gobz stream"),
        }
    }
    if !saw_version {
        return invalid("missing gobz version");
    }
    if stored_signature.as_deref() != Some(signature(&entries, false).as_slice()) {
        return invalid("Signature mismatch - seal was modified");
    }
    Ok(entries)
}

fn write_file_info_type(out: &mut impl Write) -> io::Result<()> {
    frame(out, |v| {
        signed(v, -65);
        unsigned(v, 3); // WireType.StructT
        unsigned(v, 1); // structType.CommonType
        unsigned(v, 1);
        bytes(v, b"FileInfo");
        unsigned(v, 1);
        signed(v, 65);
        unsigned(v, 0);
        unsigned(v, 1); // structType.Field
        unsigned(v, 6);
        for (name, id) in [
            ("Path", 6),
            ("RelaPath", 6),
            ("Mode", 3),
            ("Size", 2),
            ("Sha1", 5),
            ("MD5", 5),
        ] {
            unsigned(v, 1);
            bytes(v, name.as_bytes());
            unsigned(v, 1);
            signed(v, id);
            unsigned(v, 0);
        }
        unsigned(v, 0);
        unsigned(v, 0);
    })
}

fn read_file_info(input: &mut &[u8]) -> io::Result<SealEntry> {
    let mut entry = SealEntry {
        path: PathBuf::new(),
        relative_path: PathBuf::new(),
        mode: 0,
        size: 0,
        sha1: Vec::new(),
        md5: Vec::new(),
    };
    let mut field = -1i32;
    loop {
        let delta = read_unsigned(input)? as i32;
        if delta == 0 {
            break;
        }
        field += delta;
        match field {
            0 => entry.path = path_from_bytes(read_bytes(input)?)?,
            1 => entry.relative_path = path_from_bytes(read_bytes(input)?)?,
            2 => entry.mode = read_unsigned(input)? as u32,
            3 => entry.size = read_signed(input)?,
            4 => entry.sha1 = read_bytes(input)?,
            5 => entry.md5 = read_bytes(input)?,
            _ => return invalid("unsupported FileInfo field"),
        }
    }
    Ok(entry)
}

fn frame(out: &mut impl Write, fill: impl FnOnce(&mut Vec<u8>)) -> io::Result<()> {
    let mut payload = Vec::new();
    fill(&mut payload);
    let mut header = Vec::new();
    unsigned(&mut header, payload.len() as u64);
    out.write_all(&header)?;
    out.write_all(&payload)
}

fn read_frame(input: &mut impl Read) -> io::Result<Option<Vec<u8>>> {
    let Some(size) = read_unsigned_from(input)? else {
        return Ok(None);
    };
    let mut bytes = vec![0; size as usize];
    input.read_exact(&mut bytes)?;
    Ok(Some(bytes))
}

fn entry_field(out: &mut Vec<u8>, previous: &mut i32, field: i32) {
    unsigned(out, (field - *previous) as u64);
    *previous = field;
}

fn entry_field_bytes(out: &mut Vec<u8>, previous: &mut i32, field: i32, value: &[u8]) {
    if value.is_empty() {
        return;
    }
    entry_field(out, previous, field);
    bytes(out, value);
}

fn bytes(out: &mut Vec<u8>, value: &[u8]) {
    unsigned(out, value.len() as u64);
    out.extend_from_slice(value);
}

fn unsigned(out: &mut Vec<u8>, value: u64) {
    if value < 128 {
        out.push(value as u8);
        return;
    }
    let raw = value.to_be_bytes();
    let first = raw.iter().position(|v| *v != 0).unwrap_or(7);
    let len = 8 - first;
    out.push((-(len as i8)) as u8);
    out.extend_from_slice(&raw[first..]);
}

fn signed(out: &mut Vec<u8>, value: i64) {
    let encoded = if value < 0 {
        (!(value as u64) << 1) | 1
    } else {
        (value as u64) << 1
    };
    unsigned(out, encoded);
}

fn read_unsigned(input: &mut &[u8]) -> io::Result<u64> {
    let first = *input.first().ok_or_else(|| eof("gob integer"))?;
    *input = &input[1..];
    if first < 128 {
        return Ok(first as u64);
    }
    let len = (-(first as i8)) as usize;
    if len > 8 || input.len() < len {
        return invalid("invalid gob integer");
    }
    let mut value = 0;
    for byte in &input[..len] {
        value = (value << 8) | u64::from(*byte);
    }
    *input = &input[len..];
    Ok(value)
}

fn read_signed(input: &mut &[u8]) -> io::Result<i64> {
    let value = read_unsigned(input)?;
    Ok(if value & 1 == 1 {
        !(value >> 1) as i64
    } else {
        (value >> 1) as i64
    })
}

fn read_bytes(input: &mut &[u8]) -> io::Result<Vec<u8>> {
    let len = read_unsigned(input)? as usize;
    if input.len() < len {
        return Err(eof("gob bytes"));
    }
    let value = input[..len].to_vec();
    *input = &input[len..];
    Ok(value)
}

fn read_unsigned_from(input: &mut impl Read) -> io::Result<Option<u64>> {
    let mut first = [0];
    match input.read_exact(&mut first) {
        Ok(()) => {}
        Err(err) if err.kind() == io::ErrorKind::UnexpectedEof => return Ok(None),
        Err(err) => return Err(err),
    }
    if first[0] < 128 {
        return Ok(Some(first[0] as u64));
    }
    let len = (-(first[0] as i8)) as usize;
    if len > 8 {
        return invalid("invalid gob frame");
    }
    let mut raw = [0; 8];
    input.read_exact(&mut raw[8 - len..])?;
    Ok(Some(u64::from_be_bytes(raw)))
}

#[derive(Serialize, Deserialize)]
#[serde(rename = "hashlist")]
struct HashList {
    #[serde(rename = "@version")]
    version: String,
    #[serde(rename = "hash", default)]
    hashes: Vec<MhlHash>,
    #[serde(default)]
    signature: Option<MhlSignature>,
}

#[derive(Serialize, Deserialize)]
struct MhlHash {
    file: String,
    size: i64,
    #[serde(default)]
    sha1: String,
    #[serde(default)]
    md5: String,
}

#[derive(Serialize, Deserialize)]
struct MhlSignature {
    sha1: String,
}

fn write_mhl(file: File, entries: &[SealEntry]) -> io::Result<()> {
    let hashes = entries
        .iter()
        .map(|entry| MhlHash {
            file: entry.relative_path.to_string_lossy().into_owned(),
            size: entry.size,
            sha1: hex::encode(&entry.sha1),
            md5: hex::encode(&entry.md5),
        })
        .collect();
    let list = HashList {
        version: "1.0".into(),
        hashes,
        signature: Some(MhlSignature {
            sha1: hex::encode(signature(entries, true)),
        }),
    };
    let xml = to_string(&list).map_err(io::Error::other)?;
    let mut out = BufWriter::new(file);
    out.write_all(b"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")?;
    out.write_all(xml.as_bytes())?;
    out.flush()
}

fn read_mhl(file: File) -> io::Result<Vec<SealEntry>> {
    let list: HashList = from_reader(BufReader::new(file)).map_err(io::Error::other)?;
    if list.version != "1.0" || list.hashes.is_empty() {
        return invalid("unsupported or empty media hash list");
    }
    let mut entries = Vec::with_capacity(list.hashes.len());
    for hash in list.hashes {
        let sha1 = decode_hash("SHA1", &hash.file, &hash.sha1, 20)?;
        let md5 = decode_hash("MD5", &hash.file, &hash.md5, 16)?;
        if sha1.is_empty() && md5.is_empty() {
            return invalid("media hash entry has no hash");
        }
        entries.push(SealEntry {
            path: PathBuf::from(&hash.file),
            relative_path: PathBuf::from(hash.file),
            mode: 0,
            size: hash.size,
            sha1,
            md5,
        });
    }
    if let Some(expected) = list.signature {
        if hex::decode(expected.sha1).ok().as_deref() != Some(signature(&entries, true).as_slice())
        {
            return invalid("Signature mismatch - seal was modified");
        }
    }
    Ok(entries)
}

fn decode_hash(kind: &str, file: &str, value: &str, size: usize) -> io::Result<Vec<u8>> {
    if value.is_empty() {
        return Ok(Vec::new());
    }
    let bytes =
        hex::decode(value).map_err(|err| io::Error::new(io::ErrorKind::InvalidData, err))?;
    if bytes.len() != size {
        return invalid(&format!("invalid {kind} hash for '{file}'"));
    }
    Ok(bytes)
}

#[cfg(unix)]
fn path_bytes(path: &Path) -> Vec<u8> {
    use std::os::unix::ffi::OsStrExt;
    path.as_os_str().as_bytes().to_vec()
}

#[cfg(not(unix))]
fn path_bytes(path: &Path) -> Vec<u8> {
    path.to_string_lossy().as_bytes().to_vec()
}

#[cfg(unix)]
fn path_from_bytes(bytes: Vec<u8>) -> io::Result<PathBuf> {
    use std::os::unix::ffi::OsStringExt;
    Ok(std::ffi::OsString::from_vec(bytes).into())
}

#[cfg(not(unix))]
fn path_from_bytes(bytes: Vec<u8>) -> io::Result<PathBuf> {
    String::from_utf8(bytes)
        .map(PathBuf::from)
        .map_err(|err| io::Error::new(io::ErrorKind::InvalidData, err))
}

fn invalid<T>(message: &str) -> io::Result<T> {
    Err(io::Error::new(io::ErrorKind::InvalidData, message))
}

fn eof(message: &str) -> io::Error {
    io::Error::new(io::ErrorKind::UnexpectedEof, message)
}
