# SaveVector

> A vector database built from scratch in Go, with persistent vector storage and HNSW-based approximate nearest-neighbor search.

**SaveVector** is an experimental vector database written from scratch in **Go**.  
The project is focused on understanding and implementing the core building blocks behind modern vector databases — from storing vectors on disk to building an HNSW graph and performing similarity search.

The goal isn't to wrap an existing vector database. The goal is to understand **how a vector database actually works internally**.

---

## ✨ Features

- 🦫 **Built in Go**
- 💾 **Persistent vector storage**
- 🔎 **K-Nearest Neighbor (KNN) search**
- 🧠 **HNSW approximate nearest-neighbor indexing**
- 🔗 **External ID → internal ID mapping**
- 📦 **Protocol Buffer based database metadata**
- 📐 **Configurable vector dimensions**
- ⚡ **Cosine similarity search**
- 💽 **Persistent HNSW graph**
- 🗂️ **Database → Index architecture**
- 🧩 Modular storage and indexing components

---

## 🏗️ Architecture

SaveVector separates the database into a few core components:

```text
                    ┌─────────────────────┐
                    │      SaveVector     │
                    └──────────┬──────────┘
                               │
                     ┌─────────▼─────────┐
                     │      Database     │
                     │                   │
                     │  - metadata       │
                     │  - indexes        │
                     └─────────┬─────────┘
                               │
                     ┌─────────▼─────────┐
                     │       Index       │
                     │                   │
                     │  - dimension      │
                     │  - vector count   │
                     │  - vector storage │
                     └──────┬───────┬────┘
                            │       │
              ┌─────────────▼─┐   ┌─▼──────────────┐
              │ Vector Storage │   │  HNSW Index    │
              │                │   │                │
              │ vector.db     │   │ hnsw.db        │
              │ metadata.db   │   │ graph          │
              └────────────────┘   └────────────────┘
```

### Vector Storage

Vectors are flattened and converted into binary data before being written to disk.

Each vector has a fixed dimensionality:

```text
Vector
[0.12, 0.43, 0.91, ...]
        │
        ▼
  float32 → bytes
        │
        ▼
   vector.db
```

The internal ID is derived from the vector's position, while the external ID is stored separately as metadata.

---

## 🧠 HNSW Index

SaveVector implements **Hierarchical Navigable Small World (HNSW)** indexing from scratch.

The HNSW implementation maintains:

- Multiple graph layers
- An entry point
- Nodes
- Layer-specific neighbors
- Configurable `M`
- Configurable `efConstruction`
- Configurable `efSearch`

A simplified representation:

```text
Layer 2:

          A
         / \
        /   \
       C     D


Layer 1:

      A ─── B
     / \   / \
    C   D─E   F


Layer 0:

 A ─ B ─ C ─ D ─ E ─ F ─ G
  \   \ / \   \   /
   ────┴───┴───┴──
```

Higher layers provide long-range navigation, while the bottom layer contains the denser graph used for final search.

### Search

A query starts at the HNSW entry point and moves through progressively lower layers:

```text
Query
  │
  ▼
Entry Point
  │
  ▼
Highest Layer
  │
  ▼
Greedy Search
  │
  ▼
Lower Layer
  │
  ▼
Beam Search
  │
  ▼
Top-K Results
```

SaveVector currently uses **cosine similarity** for comparing vectors.

---

## 🔧 HNSW Parameters

HNSW can be configured using three parameters:

```go
h := indexing.NewHNSW(
    16,  // M
    200, // efConstruction
    50,  // efSearch
)
```

### `M`

Controls the maximum number of connections a node can have.

Higher values generally create a denser graph and can improve recall at the cost of memory and construction time.

### `efConstruction`

Controls the search width while building the graph.

Higher values can produce a better graph but increase insertion cost.

### `efSearch`

Controls the search width during querying.

Higher values generally improve search quality at the cost of query latency.

---

## 🚀 Quick Start

### Clone

```bash
git clone https://github.com/Shreyankthehacker/save-vector.git
cd save-vector
```

### Install dependencies

```bash
go mod download
```

### Run

```bash
go run .
```

---

## 💻 Example

Create a database and an index:

```go
db, err := models.CreateDatabase("mydb")
if err != nil {
    log.Fatal(err)
}

idx, err := db.CreateIndex("products", 128)
if err != nil {
    log.Fatal(err)
}
```

Create an HNSW index:

```go
h := indexing.NewHNSW(
    16,
    200,
    50,
)
```

Insert a vector:

```go
vector := []float32{
    0.12,
    0.42,
    0.87,
    // ...
}

if err := idx.InsertVector("doc-001", vector); err != nil {
    log.Fatal(err)
}
```

Then register the vector in the HNSW graph:

```go
internalID := idx.Count() - 1

if err := h.InsertVector(
    idx,
    internalID,
    vector,
); err != nil {
    log.Fatal(err)
}
```

Search for similar vectors:

```go
results, err := h.KNNSearch(
    idx,
    queryVector,
    5,
)

if err != nil {
    log.Fatal(err)
}

for _, result := range results {
    fmt.Printf(
        "ID=%d Score=%f\n",
        result.InternalID,
        result.Score,
    )
}
```

You can also search using an existing external ID:

```go
results, err := h.KNNSearchByExternalID(
    idx,
    "doc-001",
    5,
)
```

---

## 📁 Project Structure

```text
save-vector/
│
├── config/
│   └── ...
│
├── indexing/
│   ├── heap.go
│   └── hnsw.go
│
├── models/
│   ├── database.go
│   ├── databaseCreate.go
│   ├── databaseFetch.go
│   ├── index.go
│   ├── indexAccessors.go
│   ├── indexCreate.go
│   ├── indexFetch.go
│   ├── indexServices.go
│   ├── indexVectorInsert.go
│   └── vector.go
│
├── proto_models/
│   ├── db_meta.proto
│   └── db_meta.pb.go
│
├── utils/
│   ├── index.go
│   ├── logger.go
│   └── vector.go
│
├── main.go
├── go.mod
└── go.sum
```

### `models/`

Contains the database and index abstractions, vector storage, metadata management, and CRUD-style operations.

### `indexing/`

Contains the HNSW implementation and nearest-neighbor search logic.

### `proto_models/`

Contains Protocol Buffer definitions and generated Go code used for database metadata.

### `utils/`

Contains supporting functionality such as vector conversion, indexing utilities, and logging.

---

## 🔄 Insert Flow

A vector insertion follows approximately this pipeline:

```text
              Input Vector
                   │
                   ▼
            Validate Dimension
                   │
                   ▼
             Flatten Vector
                   │
                   ▼
            float32 → bytes
                   │
                   ▼
            Write to Disk
                   │
                   ▼
          Assign Internal ID
                   │
                   ▼
          Insert into HNSW
                   │
                   ▼
        Connect to Neighbors
                   │
                   ▼
          Persist HNSW Graph
```

The implementation keeps the raw vector storage and HNSW graph as separate concerns.

---

## 🔍 Search Flow

```text
                 Query Vector
                      │
                      ▼
                Entry Point
                      │
                      ▼
              Highest HNSW Layer
                      │
                      ▼
                Greedy Search
                      │
                      ▼
                  Layer 1
                      │
                      ▼
                  Layer 0
                      │
                      ▼
                Beam Search
                      │
                      ▼
                  Top-K
                      │
                      ▼
             Similarity Results
```

The HNSW implementation performs greedy traversal through higher layers and beam search at the lower layer before returning the nearest candidates.

---

## 💾 Persistence

SaveVector persists both the vectors and the HNSW graph.

The vector storage is written as binary data, while the HNSW graph is serialized separately.

This allows the graph to be loaded again rather than rebuilding the entire index every time the database starts.

---

## 🎯 Why Build a Vector Database From Scratch?

Modern vector databases can make similarity search look deceptively simple:

```python
results = db.search(query, top_k=10)
```

But behind that API are several interesting systems problems:

- How should high-dimensional vectors be stored?
- How do we map external document IDs to internal IDs?
- How can nearest-neighbor search avoid comparing against every vector?
- How does HNSW organize its graph?
- How should the graph be persisted?
- How do insertion and search interact with the index?
- How do we balance recall, memory, and latency?

SaveVector is my attempt to explore those questions by implementing the core pieces myself rather than hiding them behind an existing vector database.

---

## 🛣️ Roadmap

The project is actively evolving. Planned improvements include:

- [ ] More distance metrics
- [ ] Improved HNSW neighbor selection
- [ ] Vector deletion
- [ ] Vector updates
- [ ] Metadata / payload filtering
- [ ] Batch insertion
- [ ] Better persistence and recovery
- [ ] Concurrency support
- [ ] Benchmarks
- [ ] Recall evaluation against brute-force search
- [ ] REST API
- [ ] Client SDK
- [ ] Better memory management
- [ ] Comprehensive test suite
- [ ] Performance optimizations

---

## 📊 Benchmarking

Benchmarking is intentionally kept separate from implementation claims.

Future benchmarks will compare:

- Insert throughput
- Search latency
- Recall@K
- Memory consumption
- Disk usage
- HNSW construction time

against a brute-force baseline.

This will make it possible to measure the actual trade-offs introduced by the approximate-nearest-neighbor index.

---

## 🧪 Project Status

> **Experimental / Educational**

SaveVector is primarily a learning and systems-engineering project.

It is **not currently intended to be a production replacement for databases such as Qdrant, Milvus, Weaviate, or Pinecone**.

The interesting part of the project is the implementation itself: understanding how vector storage, indexing, graph traversal, persistence, and similarity search fit together.

---

## 🤝 Contributing

Contributions, ideas, benchmarks, and discussions are welcome.

If you're interested in vector databases, information retrieval, ANN algorithms, or database internals, feel free to explore the code and open an issue or pull request.

---

## 📚 What This Project Demonstrates

SaveVector is a practical exploration of:

```text
Vector Databases
      │
      ├── Vector Storage
      │
      ├── Binary Serialization
      │
      ├── Metadata Management
      │
      ├── ID Mapping
      │
      ├── Similarity Functions
      │
      ├── HNSW
      │    ├── Graph Construction
      │    ├── Layering
      │    ├── Greedy Search
      │    └── Beam Search
      │
      └── Persistent Indexes
```



---

## 📄 License

See the `LICENSE` file for license information.
