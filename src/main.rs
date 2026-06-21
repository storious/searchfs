use std::env;
use std::io;
use std::path::Path;
use std::time::Instant;

use searchfs::engine::SearchEngine;
use searchfs::query::QueryMode;
use searchfs::snapshot;

fn main() -> io::Result<()> {
    let mut args = env::args().skip(1);

    match args.next().as_deref() {
        Some("build") => {
            let docs = args.next().unwrap_or_else(|| "./docs".to_string());
            let index = args.next().unwrap_or_else(|| "searchfs.idx".to_string());
            run_build(&docs, &index)
        }
        Some("search") => {
            let index = args.next().unwrap_or_else(|| "searchfs.idx".to_string());
            let query = args.next().unwrap_or_else(|| "rust".to_string());
            let limit = args
                .next()
                .and_then(|s| s.parse::<usize>().ok())
                .unwrap_or(10);
            let mode = args.next().unwrap_or_else(|| "and".to_string());

            run_search(&index, &query, limit, &mode)
        }
        _ => {
            eprintln!("usage:");
            eprintln!("  searchfs build <docs> <index>");
            eprintln!("  searchfs search <index> <query> [limit] [and|or|phrase]");
            Ok(())
        }
    }
}

fn run_build(docs: &str, index_path: &str) -> io::Result<()> {
    let mut engine = SearchEngine::new();

    let start = Instant::now();
    engine.index_dir(Path::new(docs))?;
    let elapsed = start.elapsed();

    let stats = engine.stats();

    eprintln!(
        "indexed docs={} terms={} postings={} positions={} index_time={:.2?}",
        engine.doc_count(),
        stats.terms,
        stats.postings,
        stats.total_positions,
        elapsed,
    );

    let snapshot = engine.into_snapshot();
    snapshot::save(Path::new(index_path), &snapshot)?;

    eprintln!("saved index={index_path}");

    Ok(())
}

fn run_search(index_path: &str, query: &str, limit: usize, mode_arg: &str) -> io::Result<()> {
    let load_start = Instant::now();
    let snapshot = snapshot::load(Path::new(index_path))?;
    let engine = SearchEngine::from_snapshot(snapshot);
    let load_elapsed = load_start.elapsed();

    let mode = if query.starts_with('"') && query.ends_with('"') {
        QueryMode::Phrase
    } else {
        QueryMode::try_from(mode_arg)
            .map_err(|msg| io::Error::new(io::ErrorKind::InvalidInput, msg))?
    };

    let query = query.trim_matches('"');

    let search_start = Instant::now();
    let results = engine.search(query, mode);
    let search_elapsed = search_start.elapsed();

    eprintln!("load_time={:.2?}", load_elapsed);
    eprintln!("search_time={:.2?}", search_elapsed);

    for result in results.into_iter().take(limit) {
        println!("{} score={}", result.path, result.score);
    }

    Ok(())
}
