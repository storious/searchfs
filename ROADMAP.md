# SearchFS Roadmap

SearchFS is a learning project that incrementally implements the core
building blocks of a modern search engine.

---

## v0.1.x — In-Memory Search

- In-memory inverted index
- Tokenization
- AND / OR / Phrase search
- TF-IDF ranking
- CLI

---

## v0.2.x — Persistence

- Snapshot serialization
- Snapshot loading
- Incremental indexing
- Persistent search

---

## v0.3.x — Segment Architecture

- Immutable segments
- Multi-segment search
- Segment merge / compaction
- Segment metadata
- Document metadata
- BM25 ranking
- Segment inspection
- Interactive REPL
- Reader / scorer refactoring

---

## v0.4.x — Query Engine

- Skip lists
- Posting list compression
- Binary term dictionary
- Query planner
- Top-K collector

---

## v0.5.x — Storage Engine

✓ Storage abstraction
✓ Local filesystem backend
□ In-memory storage backend
□ Segment cache
□ Memory-mapped local storage
□ Background merge scheduler
□ Parallel segment search

---

## v0.6.x — Search Service

- HTTP API
- REST search service
- JSON responses
- Concurrent query execution
- Search metrics

---

## v0.7.x — Distributed Search

- Metadata service
- Sharding
- Replication
- Distributed query execution
- Cluster management
