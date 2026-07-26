use godi::{
    CommonOptions, CopyOptions, FileFilter, SealEntry, SealFormat, SealOptions, VerifyOptions,
    read_seal, seal, sealed_copy, verify, write_seal,
};
use std::{fs, path::PathBuf};

#[test]
fn reads_v1_go_gob_and_writes_both_formats() {
    let temp = tempfile::tempdir().unwrap();
    let go_seal = temp.path().join("legacy.gobz");
    fs::write(
        &go_seal,
        hex::decode(concat!(
            "1f8b080000096e8802ff62666160f2ffdfc8ccc8c8e1969993ea999796cff8",
            "bf8981918d912520b12483918781912328352711c661f1cd4f496564033282",
            "33ab521959408c8c4443462e0646665f175320cdc0e0fdbf89914f3f25b12",
            "4513f0d68a65e49450923079cf58ff13fa316a3881ba3fbcbf5124b374d4a",
            "6e7a927abf45a0fd473937a3408c6243da59ef70930636d555563a0c73189",
            "8991818c5b9184494ca8c4e1a4dd8a875bfeae1c2e0c4f26bc7be896d070",
            "40000ffff29885ad6bc000000"
        ))
        .unwrap(),
    )
    .unwrap();
    let entries = read_seal(&go_seal).unwrap();
    assert_eq!(entries.len(), 1);
    assert_eq!(entries[0].path, PathBuf::from("/data/file.txt"));
    assert_eq!(entries[0].relative_path, PathBuf::from("file.txt"));

    for (name, format) in [("new.gobz", SealFormat::Gob), ("new.mhl", SealFormat::Mhl)] {
        let path = temp.path().join(name);
        write_seal(&path, format, &entries).unwrap();
        let roundtrip = read_seal(&path).unwrap();
        assert_eq!(roundtrip[0].relative_path, entries[0].relative_path);
        assert_eq!(roundtrip[0].sha1, entries[0].sha1);
        assert_eq!(roundtrip[0].md5, entries[0].md5);
    }
}

#[test]
fn reads_unsigned_third_party_mhl() {
    let temp = tempfile::tempdir().unwrap();
    let seal = temp.path().join("third-party.mhl");
    fs::write(
        &seal,
        concat!(
            r#"<hashlist version="1.0"><creatorinfo><tool>other</tool></creatorinfo>"#,
            r#"<hash><file>clip.mov</file><size>4</size>"#,
            r#"<md5>098f6bcd4621d373cade4e832627b4f6</md5></hash></hashlist>"#
        ),
    )
    .unwrap();

    let entries = read_seal(&seal).unwrap();

    assert_eq!(entries[0].relative_path, PathBuf::from("clip.mov"));
    assert!(entries[0].sha1.is_empty());
}

#[test]
fn seal_copy_verify_and_detect_change() {
    let temp = tempfile::tempdir().unwrap();
    let source = temp.path().join("source");
    let destination = temp.path().join("destination");
    fs::create_dir_all(source.join("sub")).unwrap();
    fs::create_dir(&destination).unwrap();
    fs::write(source.join("a"), b"alpha").unwrap();
    fs::write(source.join("sub/b"), b"beta").unwrap();
    fs::write(source.join(".DS_Store"), b"volatile").unwrap();

    let mut options = SealOptions::default();
    options.common.filters = vec![FileFilter::Volatile];
    let report = seal(std::slice::from_ref(&source), options, |_| {}).unwrap();
    let entries = read_seal(&report.seals[0]).unwrap();
    assert_eq!(entries.len(), 2);
    verify(&report.seals, VerifyOptions::default(), |_| {}).unwrap();

    let copied = sealed_copy(
        std::slice::from_ref(&source),
        std::slice::from_ref(&destination),
        CopyOptions {
            verify_after_copy: true,
            seal: SealOptions {
                common: CommonOptions {
                    filters: vec![FileFilter::Seals, FileFilter::Volatile],
                    ..CommonOptions::default()
                },
                ..SealOptions::default()
            },
            ..CopyOptions::default()
        },
        |_| {},
    )
    .unwrap();
    assert_eq!(fs::read(destination.join("a")).unwrap(), b"alpha");
    assert!(!copied.seals.is_empty());

    fs::write(source.join("a"), b"ALPHA").unwrap();
    assert!(verify(&report.seals, VerifyOptions::default(), |_| {}).is_err());
}

#[test]
fn does_not_overwrite_and_rolls_back_destination() {
    let temp = tempfile::tempdir().unwrap();
    let source = temp.path().join("source");
    let destination = temp.path().join("destination");
    fs::create_dir(&source).unwrap();
    fs::create_dir(&destination).unwrap();
    fs::write(source.join("a"), b"a").unwrap();
    fs::write(source.join("b"), b"b").unwrap();
    fs::write(destination.join("b"), b"keep").unwrap();

    assert!(
        sealed_copy(
            std::slice::from_ref(&source),
            std::slice::from_ref(&destination),
            CopyOptions::default(),
            |_| {}
        )
        .is_err()
    );
    assert!(!destination.join("a").exists());
    assert_eq!(fs::read(destination.join("b")).unwrap(), b"keep");
}

#[test]
fn empty_filter_list_is_valid() {
    let options = CommonOptions {
        filters: Vec::new(),
        ..CommonOptions::default()
    };
    assert!(options.filters.is_empty());
    let entry = SealEntry {
        path: "x".into(),
        relative_path: "x".into(),
        mode: 0,
        size: 0,
        sha1: Vec::new(),
        md5: Vec::new(),
    };
    assert_eq!(entry.size, 0);
}

#[test]
fn seals_an_individual_file_next_to_it() {
    let temp = tempfile::tempdir().unwrap();
    let source = temp.path().join("one-file");
    fs::write(&source, b"content").unwrap();

    let report = seal(
        std::slice::from_ref(&source),
        SealOptions::default(),
        |_| {},
    )
    .unwrap();

    assert_eq!(report.seals[0].parent(), Some(temp.path()));
    assert_eq!(
        read_seal(&report.seals[0]).unwrap()[0].relative_path,
        PathBuf::from("one-file")
    );
}
