use crate::query::QueryMode;
use crate::query::SearchResult;
use crate::segment::reader::SegmentReaderCache;
use crate::segment::search::SegmentSearcher;

use std::io;

pub(crate) fn search_with_cache(
    cache: &SegmentReaderCache,
    query: &str,
    mode: QueryMode,
) -> io::Result<Vec<SearchResult>> {
    let terms: Vec<String> = crate::index::parser::tokenize(query)
        .into_iter()
        .map(|(term, _)| term)
        .collect();

    let mut all_results = Vec::new();

    for reader in cache.readers() {
        let searcher = SegmentSearcher::new(reader);

        match mode {
            QueryMode::All => {
                all_results.extend(searcher.search_all(&terms)?);
            }
            QueryMode::Any => {
                all_results.extend(searcher.search_any(&terms)?);
            }
            QueryMode::Phrase => {
                all_results.extend(searcher.search_phrase(&terms)?);
            }
        }
    }

    all_results.sort_by(|a, b| {
        b.score
            .total_cmp(&a.score)
            .then_with(|| a.path.cmp(&b.path))
    });

    Ok(all_results)
}

pub(crate) fn print_results(results: Vec<SearchResult>, limit: usize) {
    for result in results.into_iter().take(limit) {
        println!("{} score={}", result.path, result.score);
    }
}

pub(crate) fn print_repl_help() {
    eprintln!("commands:");
    eprintln!("  :help             show this help");
    eprintln!("  :limit <n>        set result limit");
    eprintln!("  :mode and         use AND search");
    eprintln!("  :mode or          use OR search");
    eprintln!("  :mode phrase      use phrase search");
    eprintln!("  :stats            show index and REPL stats");
    eprintln!("  :q, :quit         exit");
    eprintln!();
    eprintln!("queries:");
    eprintln!("  rust memory       search with current mode");
    eprintln!("  \"white whale\"     force phrase search");
}

pub(crate) fn print_repl_stats(cache: &SegmentReaderCache, mode: QueryMode, limit: usize) {
    let mut total_docs = 0usize;
    let mut total_terms = 0usize;
    let mut total_postings = 0usize;
    let mut total_positions = 0usize;

    for reader in cache.readers() {
        total_docs += reader.doc_count();
        total_terms += reader.term_count();
        total_postings += reader.posting_count();
        total_positions += reader.position_count();
    }

    let avg_doc_len = if total_docs == 0 {
        0.0
    } else {
        total_positions as f64 / total_docs as f64
    };

    eprintln!("segments={}", cache.readers().len());
    eprintln!("docs={total_docs}");
    eprintln!("terms={total_terms}");
    eprintln!("postings={total_postings}");
    eprintln!("positions={total_positions}");
    eprintln!("avg_doc_len={avg_doc_len:.2}");
    eprintln!("mode={}", mode.as_str());
    eprintln!("limit={limit}");
}
