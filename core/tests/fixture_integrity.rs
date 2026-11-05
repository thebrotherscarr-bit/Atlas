//! Golden-master integrity guards: the cut fixtures are the reference every
//! future Rust assertion meets, so their bytes are pinned here. A failure
//! means a fixture was edited after cutting — fold law applies: re-cut with
//! tools/cut_fixtures.py and witness it, never silently patch.

use std::fs;
use std::path::PathBuf;

use atlas_core::sha256::hex_digest;

/// The Archive-cut fixtures left the repository 2026-11-05 and stay local
/// only. Without them this proof skips, and says why, rather than failing.
fn fixtures_absent() -> bool {
    let d = std::path::PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join("..")
        .join("tests")
        .join("fixtures");
    if d.is_dir() {
        return false;
    }
    eprintln!("SKIP: tests/fixtures absent (Archive-cut, local only since 2026-11-05)");
    true
}

fn fixture(rel: &str) -> Vec<u8> {
    fs::read(
        PathBuf::from(env!("CARGO_MANIFEST_DIR"))
            .join("..")
            .join("tests")
            .join("fixtures")
            .join(rel),
    )
    .expect("fixture present")
}

#[test]
fn steward_chain_bytes_pinned() {
    if fixtures_absent() {
        return;
    }
    assert_eq!(
        hex_digest(&fixture("chains/steward_chain.jsonl")),
        "fcbbc39e9fdda22222f453c5fcc78e641f6329ecc21b1cf5eaf50f14244c34bb"
    );
}

#[test]
fn agents_seatlog_bytes_pinned() {
    if fixtures_absent() {
        return;
    }
    assert_eq!(
        hex_digest(&fixture("chains/agents_seatlog.jsonl")),
        "fea6dd33e7d1ceb6a2f17348a1e7536b0eee328afe5473fcb8c5c1af7f9f001f"
    );
}
