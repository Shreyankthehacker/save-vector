// example_usage.go — not a real file to compile, just a usage reference.
//
// This shows how to wire InsertVector and KNNSearch into your existing flow.

package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Shreyankthehacker/savector/indexing"
	"github.com/Shreyankthehacker/savector/models"

)

func main() {
	// ── 1. Create or open a database and index ─────────────────────────────────
	db, err := models.CreateDatabase("mydb")
	if err != nil {
		log.Fatal(err)
	}

	idx, err := db.CreateIndex("products", 128) // 128-dimensional vectors
	if err != nil {
		log.Fatal(err)
	}

	// ── 2. Create an HNSW index with sensible defaults ─────────────────────────
	//   M=16          → max 16 neighbours per node per layer (32 at layer 0)
	//   efConstruct=200 → beam width during insertion (higher = better graph quality)
	//   efSearch=50   → beam width during search (higher = more accurate, slower)
	h := indexing.NewHNSW(16, 200, 50)

	// ── 3. Insert vectors ──────────────────────────────────────────────────────
	// Step A: write the raw vector to disk via the existing models method.
	// Step B: register the vector in the HNSW graph.
	//
	// Both steps MUST happen together; wrap them in a helper like insertFull().

	vectors := []struct {
		id  string
		vec []float32
	}{
		{"doc-001", makeRandVector(128)},
		{"doc-002", makeRandVector(128)},
		{"doc-003", makeRandVector(128)},
	}

	for _, item := range vectors {
		if err := insertFull(idx, h, item.id, item.vec); err != nil {
			log.Fatalf("insert %s: %v", item.id, err)
		}
		fmt.Printf("Inserted %s (count=%d)\n", item.id, idx.Count())
	}

	// ── 4. Search ──────────────────────────────────────────────────────────────
	query := makeRandVector(128)
	k := 3

	results, err := h.KNNSearch(idx, query, k)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTop-%d results:\n", k)
	for rank, r := range results {
		fmt.Printf("  #%d  internalID=%d  score=%.4f\n", rank+1, r.InternalID, r.Score)
	}

	// ── 5. Search by external ID (find neighbours of an existing vector) ───────
	neighbours, err := h.KNNSearchByExternalID(idx, "doc-001", k)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nNeighbours of doc-001:\n")
	for rank, r := range neighbours {
		fmt.Printf("  #%d  internalID=%d  score=%.4f\n", rank+1, r.InternalID, r.Score)
	}
}

// insertFull writes the vector to disk AND inserts it into the HNSW graph.
// Always call these two steps together.
func insertFull(idx *models.Index, h *indexing.HNSW, externalID string, vec []float32) error {
	// 1. Write to vector.db and payload.db (existing logic in indexVectorInsert.go)
	if err := idx.InsertVector(externalID, vec); err != nil {
		return fmt.Errorf("raw insert: %w", err)
	}

	// 2. Register in HNSW graph (writes to hnsw.db)
	internalID := idx.Count() - 1 // InsertVector increments count before returning
	if err := h.InsertVector(idx, internalID, vec); err != nil {
		return fmt.Errorf("hnsw insert: %w", err)
	}

	return nil
}

// makeRandVector is just a stand-in; replace with real data.
func makeRandVector(dim int) []float32 {
	v := make([]float32, dim)
	for i := range v {
		v[i] = rand.Float32() // use whatever you have
	}
	return v
}