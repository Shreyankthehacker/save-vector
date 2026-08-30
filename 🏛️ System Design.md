# 🏛️ System Design

SaveVector is designed as a lightweight, persistent vector database where **storage, indexing, metadata, and search are separated into independent components**.

The system can be viewed at two levels:

1. **High-Level Design (HLD)** — how the major components interact.
2. **Low-Level Design (LLD)** — how vectors, indexes, IDs, and HNSW nodes are represented and persisted.

---

# High-Level Design

At a high level, SaveVector follows this architecture:

```text
                         ┌───────────────────────┐
                         │       Client          │
                         │                       │
                         │  Insert / Search /    │
                         │  Fetch / Manage DB    │
                         └───────────┬───────────┘
                                     │
                                     ▼
                         ┌───────────────────────┐
                         │      Database API     │
                         │                       │
                         │  Database             │
                         │  Index                │
                         │  Vector Operations    │
                         └───────────┬───────────┘
                                     │
                  ┌──────────────────┼──────────────────┐
                  │                  │                  │
                  ▼                  ▼                  ▼
        ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
        │    Metadata     │ │ Vector Storage  │ │  HNSW Index     │
        │                 │ │                 │ │                 │
        │ DB metadata     │ │ Raw vectors     │ │ ANN graph       │
        │ Index metadata  │ │ Binary storage  │ │ Search          │
        │ Configuration   │ │                 │ │ Traversal       │
        └────────┬────────┘ └────────┬────────┘ └────────┬────────┘
                 │                   │                   │
                 ▼                   ▼                   ▼
        ┌─────────────────────────────────────────────────────┐
        │                    Persistent Storage                │
        │                                                     │
        │   metadata.pb      vectors.db       hnsw.db        │
        └─────────────────────────────────────────────────────┘
```

The major design principle is:

> **The vector store is responsible for storing vectors. The HNSW index is responsible for finding vectors. Metadata connects the two.**

---

## Core Components

### 1. Database

The `Database` is the top-level abstraction.

```text
Database
   │
   ├── Database Metadata
   │
   └── Indexes
          │
          ├── Index A
          ├── Index B
          └── Index C
```

A database can contain multiple indexes.

Each index represents a collection of vectors with a fixed dimensionality.

---

### 2. Index

An `Index` represents a vector collection.

Conceptually:

```text
Index
 │
 ├── Name
 ├── Dimension
 ├── Vector Count
 ├── Vector Storage
 ├── External → Internal ID mapping
 └── HNSW Index
```

For example:

```text
products
 ├── dimension = 768
 ├── vectors = 1,000,000
 └── HNSW
```

The fixed dimensionality allows the system to validate vectors before inserting them.

---

### 3. Vector Storage

The vector storage layer is responsible for storing the actual vector values.

```text
External ID
     │
     ▼
Internal ID
     │
     ▼
Vector Offset
     │
     ▼
┌──────────────────────────┐
│       vectors.db         │
│                          │
│ Vector 0                 │
│ Vector 1                 │
│ Vector 2                 │
│ ...                      │
│ Vector N                 │
└──────────────────────────┘
```

Vectors are stored in binary form rather than as JSON.

For a `float32` vector:

```text
[0.12, 0.45, 0.91, ...]
          │
          ▼
      float32 bytes
          │
          ▼
      Persistent file
```

This reduces storage overhead compared with text-based serialization.

---

### 4. HNSW Index

The HNSW index provides approximate nearest-neighbor search.

Instead of comparing the query against every vector:

```text
Query
  │
  ├── Vector 1
  ├── Vector 2
  ├── Vector 3
  ├── ...
  └── Vector N
```

SaveVector navigates an HNSW graph:

```text
                    Query
                      │
                      ▼
                Entry Point
                      │
                      ▼
                Layer 2
              ┌───────┴───────┐
              ▼               ▼
             A                 D
              \               /
               └──────┬──────┘
                      ▼
                    Layer 1
                 ┌────┼────┐
                 ▼    ▼    ▼
                 B    C    E
                      │
                      ▼
                    Layer 0
              ───────────────────
              Dense candidate graph
                      │
                      ▼
                   Top-K
```

The graph allows the search to avoid scanning the entire vector collection.

---

# Low-Level Design

The low-level design focuses on the internal representation of vectors, IDs, indexes, and HNSW nodes.

---

## 1. Database → Index Relationship

The relationship can be represented as:

```text
Database
│
├── metadata
│
└── indexes
      │
      ├── Index
      │    ├── metadata
      │    ├── vector storage
      │    └── HNSW
      │
      └── Index
           ├── metadata
           ├── vector storage
           └── HNSW
```

This allows different collections to exist independently.

For example:

```text
Database: ecommerce

├── products
│    ├── dimension = 768
│    └── HNSW
│
├── users
│    ├── dimension = 384
│    └── HNSW
│
└── documents
     ├── dimension = 1536
     └── HNSW
```

---

# 2. External ID vs Internal ID

SaveVector separates the identifier exposed to the application from the identifier used internally by the index.

```text
Application
     │
     │ "document-123"
     ▼
External ID
     │
     ▼
Internal ID
     │
     ▼
Vector Storage
     │
     ▼
Vector
```

Example:

```text
External ID        Internal ID

"doc-001"     ───►     0
"doc-002"     ───►     1
"doc-003"     ───►     2
"doc-004"     ───►     3
```

The HNSW graph operates on the internal IDs.

This avoids storing large or arbitrary external identifiers throughout the graph.

---

# 3. HNSW Node Design

Each HNSW node represents a vector in the collection.

Conceptually:

```text
HNSWNode
│
├── Internal ID
│
├── Level
│
└── Neighbors
      │
      ├── Layer 0
      │    ├── Node A
      │    ├── Node B
      │    └── Node C
      │
      ├── Layer 1
      │    ├── Node D
      │    └── Node E
      │
      └── Layer 2
           └── Node F
```

A node can therefore participate in multiple graph layers.

---

# 4. HNSW Graph

The HNSW graph consists of multiple layers.

```text
Layer 2

        A ───────── D
        │
        │
        G


Layer 1

     A ─── B ─── D
      \    │    /
       \   │   /
          C


Layer 0

 A ─ B ─ C ─ D ─ E ─ F ─ G ─ H
   \   \ / \   / \   /
    ────┴───┴───┴──
```

### Layer 0

Contains the densest representation of the graph.

It provides the final candidate search.

### Higher Layers

Contain fewer nodes and provide long-range connections.

These layers allow the search algorithm to quickly move toward the region of the graph containing the nearest vectors.

---

# 5. Vector Insertion

The complete insertion pipeline is:

```text
             Insert(id, vector)
                    │
                    ▼
          Validate vector dimension
                    │
                    ▼
             Generate Internal ID
                    │
                    ▼
          Persist vector to storage
                    │
                    ▼
        Update ID → Internal ID mapping
                    │
                    ▼
           Insert node into HNSW
                    │
                    ▼
       Select neighbors at each layer
                    │
                    ▼
          Update graph connections
                    │
                    ▼
          Persist HNSW graph
```

### HNSW insertion

A simplified insertion process:

```text
New Vector
     │
     ▼
Randomly select maximum level
     │
     ▼
Start from entry point
     │
     ▼
Search from highest layer
     │
     ▼
Move toward closer nodes
     │
     ▼
Drop to next layer
     │
     ▼
Repeat
     │
     ▼
Reach Layer 0
     │
     ▼
Find candidate neighbors
     │
     ▼
Connect new node
     │
     ▼
Update neighbor connections
```

---

# 6. KNN Search

Given:

```text
Query Vector Q
K = 5
```

the search process is:

```text
                    Query
                      │
                      ▼
                Entry Point
                      │
                      ▼
               Highest Layer
                      │
                Greedy Search
                      │
                      ▼
                 Next Layer
                      │
                Greedy Search
                      │
                      ▼
                  Layer 0
                      │
                Beam Search
                      │
                      ▼
             Candidate Neighbors
                      │
                      ▼
              Similarity Ranking
                      │
                      ▼
                   Top K
```

The search uses `efSearch` to control the number of candidates explored.

Conceptually:

```text
Small efSearch
     │
     ├── Faster search
     └── Potentially lower recall

Large efSearch
     │
     ├── More candidates
     ├── Better recall
     └── Higher latency
```

---

# 7. Similarity Calculation

SaveVector currently uses **cosine similarity**.

For vectors:

```text
A = [a₁, a₂, ..., aₙ]

B = [b₁, b₂, ..., bₙ]
```

cosine similarity is:

```text
                 A · B
similarity = ───────────────
             ||A|| × ||B||
```

The value represents the angular similarity between two vectors.

Conceptually:

```text
              B
             /
            /
           / θ
          /
---------/----------> A

small θ  → high similarity
large θ  → low similarity
```

---

# 8. Persistence Design

SaveVector separates persistent data into different files.

Conceptually:

```text
Database Directory
│
├── metadata
│     │
│     └── Database / Index configuration
│
├── vectors
│     │
│     └── Raw vector data
│
└── hnsw
      │
      └── HNSW graph
```

This separation has an important advantage:

```text
Vector Data
     ≠
Index Data
     ≠
Metadata
```

The vector itself doesn't need to contain graph information, while the graph doesn't need to duplicate the vector values.

---

# 9. Persistence / Recovery

When SaveVector starts again:

```text
             Application Restart
                     │
                     ▼
              Open Database
                     │
                     ▼
             Read Metadata
                     │
                     ▼
              Open Index
                     │
          ┌──────────┴──────────┐
          ▼                     ▼
   Load Vector Storage     Load HNSW Graph
          │                     │
          └──────────┬──────────┘
                     ▼
              Ready for Search
```

The objective is to avoid rebuilding the entire HNSW graph from the raw vectors every time the database is opened.

---

# 10. Memory vs Disk

One of the key design decisions in SaveVector is separating the **persistent representation** from the **in-memory representation**.

```text
                    SaveVector
                        │
              ┌─────────┴─────────┐
              │                   │
              ▼                   ▼
          In Memory             On Disk
              │                   │
              │                   ├── Metadata
              │                   ├── Vectors
              │                   └── HNSW
              │
              ├── Graph
              ├── Search state
              └── Runtime objects
```

Disk provides persistence, while memory provides the structures required for fast graph traversal and search.

---

# 🔄 Complete System Flow

Putting everything together:

```text
                       ┌──────────────┐
                       │    Client    │
                       └──────┬───────┘
                              │
                    Insert / Search
                              │
                              ▼
                    ┌─────────────────┐
                    │    Database     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │      Index      │
                    └───────┬─────────┘
                            │
             ┌──────────────┼──────────────┐
             │              │              │
             ▼              ▼              ▼
        Metadata      Vector Storage     HNSW
             │              │              │
             │              │              │
             ▼              ▼              ▼
       metadata.pb      vectors.db      hnsw.db
             │              │              │
             └──────────────┼──────────────┘
                            │
                            ▼
                       Persistent
                        Database
```

---

# ⚡ Search Path

The complete query path is:

```text
                    Query Vector
                         │
                         ▼
                     Database
                         │
                         ▼
                       Index
                         │
                         ▼
                       HNSW
                         │
                         ▼
                  Graph Traversal
                         │
                         ▼
                  Candidate IDs
                         │
                         ▼
                  Vector Storage
                         │
                         ▼
                Retrieve Candidates
                         │
                         ▼
                Similarity Evaluation
                         │
                         ▼
                     Ranking
                         │
                         ▼
                      Top-K
                         │
                         ▼
                    Client
```

This separation is important because **HNSW is an index, not the source of truth for the vector data**.

---

# 🧩 Design Principles

SaveVector follows several core design principles.

### Separation of concerns

```text
Storage       → Stores vectors
HNSW          → Finds candidate vectors
Metadata      → Describes the database
Database      → Coordinates everything
```

### Persistence first

Data structures are designed with persistence in mind rather than treating disk storage as an afterthought.

### IDs are separated

External application IDs are separated from internal graph IDs.

### Approximate search

HNSW trades exactness for significantly more efficient nearest-neighbor search compared with brute-force scanning.

### Configurable search quality

Parameters such as:

```text
M
efConstruction
efSearch
```

allow the user to control the trade-off between:

```text
Memory
   ↕
Build Time
   ↕
Search Latency
   ↕
Recall
```

---

# 🔮 Future System Design

The architecture can be extended toward a client/server vector database:

```text
                   ┌─────────────────┐
                   │      Client     │
                   └────────┬────────┘
                            │
                         HTTP/gRPC
                            │
                            ▼
                   ┌─────────────────┐
                   │   SaveVector    │
                   │     Server      │
                   └────────┬────────┘
                            │
             ┌──────────────┼──────────────┐
             ▼              ▼              ▼
        Query Engine    Index Manager   Storage
             │              │              │
             ▼              ▼              ▼
           HNSW          Collections     Disk
```

Potential future components include:

- REST / gRPC API
- Metadata filtering
- Batch ingestion
- Concurrent reads/writes
- WAL / crash recovery
- Background index building
- Sharding
- Replication
- Quantization
- SIMD-optimized distance calculations
- Distributed search

This would evolve SaveVector from an embedded vector database implementation into a standalone vector database service.