package indexing













import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/Shreyankthehacker/savector/config"
	"github.com/Shreyankthehacker/savector/models"
)






type HNSW struct {
	M           int     
	EfConstruct int     
	EfSearch    int     
	ml          float64 
}



func NewHNSW(M, efConstruct, efSearch int) *HNSW {
	return &HNSW{
		M:           M,
		EfConstruct: efConstruct,
		EfSearch:    efSearch,
		ml:          1.0 / math.Log(float64(M)),
	}
}

func (h *HNSW) randomLevel() int {
	return int(math.Floor(-math.Log(rand.Float64()) * h.ml))
}

func (h *HNSW) maxNeighbours(layer int) int {
	if layer == 0 {
		return h.M * 2
	}
	return h.M
}





type hnswNode struct {
	id        int
	maxLayer  int
	neighbors [][]int 
}

type hnswGraph struct {
	nodes      []*hnswNode
	entryPoint int 
	maxLayer   int
}

func newGraph() *hnswGraph {
	return &hnswGraph{entryPoint: -1, maxLayer: 0}
}

func (g *hnswGraph) addNode(maxLayer int) int {
	id := len(g.nodes)
	nb := make([][]int, maxLayer+1)
	for i := range nb {
		nb[i] = []int{}
	}
	g.nodes = append(g.nodes, &hnswNode{id: id, maxLayer: maxLayer, neighbors: nb})
	return id
}


func (g *hnswGraph) connect(id1, id2, layer int) {
	g.nodes[id1].neighbors[layer] = append(g.nodes[id1].neighbors[layer], id2)
	g.nodes[id2].neighbors[layer] = append(g.nodes[id2].neighbors[layer], id1)
}



func (g *hnswGraph) neighboursAt(id, layer int) []int {
	n := g.nodes[id]
	if layer > n.maxLayer {
		return nil
	}
	return n.neighbors[layer]
}



func (g *hnswGraph) pruneNeighbours(id, layer, maxM int, vecs [][]float32) {
	n := g.nodes[id]
	if layer > n.maxLayer || len(n.neighbors[layer]) <= maxM {
		return
	}
	base := vecs[id]
	nb := n.neighbors[layer]

	type sc struct {
		nid   int
		score float32
	}
	scored := make([]sc, len(nb))
	for i, nid := range nb {
		scored[i] = sc{nid, cosineSim(base, vecs[nid])}
	}
	sort.Slice(scored, func(i, j int) bool { return scored[i].score > scored[j].score })

	trimmed := make([]int, maxM)
	for i := range trimmed {
		trimmed[i] = scored[i].nid
	}
	n.neighbors[layer] = trimmed
}
















func saveGraph(g *hnswGraph, path string) error {
	buf := new(bytes.Buffer)
	wi := func(v int32) { binary.Write(buf, binary.LittleEndian, v) }

	wi(int32(len(g.nodes)))
	wi(int32(g.entryPoint))
	wi(int32(g.maxLayer))

	for _, n := range g.nodes {
		wi(int32(n.id))
		wi(int32(n.maxLayer))
		for layer := 0; layer <= n.maxLayer; layer++ {
			nb := n.neighbors[layer]
			wi(int32(len(nb)))
			for _, nid := range nb {
				wi(int32(nid))
			}
		}
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func loadGraph(path string) (*hnswGraph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return newGraph(), nil
	}

	r := bytes.NewReader(data)
	ri := func() int32 {
		var v int32
		binary.Read(r, binary.LittleEndian, &v)
		return v
	}

	numNodes := int(ri())
	entryPoint := int(ri())
	maxLayer := int(ri())

	g := &hnswGraph{
		entryPoint: entryPoint,
		maxLayer:   maxLayer,
		nodes:      make([]*hnswNode, 0, numNodes),
	}

	for i := 0; i < numNodes; i++ {
		id := int(ri())
		ml := int(ri())
		nb := make([][]int, ml+1)
		for layer := 0; layer <= ml; layer++ {
			count := int(ri())
			nb[layer] = make([]int, count)
			for j := 0; j < count; j++ {
				nb[layer][j] = int(ri())
			}
		}
		g.nodes = append(g.nodes, &hnswNode{id: id, maxLayer: ml, neighbors: nb})
	}

	fmt.Printf("[hnsw] loaded graph: %d nodes, entry=%d, maxLayer=%d\n",
		len(g.nodes), g.entryPoint, g.maxLayer)
	return g, nil
}







func loadAllVectors(idx *models.Index) ([][]float32, error) {
	count := idx.Count()
	if count == 0 {
		return [][]float32{}, nil
	}
	vecs := make([][]float32, count)
	for i := 0; i < count; i++ {
		v, err := idx.FetchVectorByInternalId(i)
		if err != nil {
			return nil, fmt.Errorf("loadAllVectors id=%d: %w", i, err)
		}
		vecs[i] = v.Data()
	}
	return vecs, nil
}



func LoadAllVectorFromLayer(idx *models.Index, layer int) ([][]float32, error) {
	g, err := loadGraph(idx.Path() + "/" + config.INDEX_HNSW_FILE)
	if err != nil {
		return nil, err
	}
	all, err := loadAllVectors(idx)
	if err != nil {
		return nil, err
	}
	var result [][]float32
	for _, n := range g.nodes {
		if n.maxLayer >= layer {
			result = append(result, all[n.id])
		}
	}
	return result, nil
}







func (h *HNSW) InsertVector(idx *models.Index, internalID int, vector []float32) error {
	hnswPath := idx.Path() + "/" + config.INDEX_HNSW_FILE

	g, err := loadGraph(hnswPath)
	if err != nil {
		return fmt.Errorf("InsertVector: %w", err)
	}
	allVecs, err := loadAllVectors(idx)
	if err != nil {
		return fmt.Errorf("InsertVector: %w", err)
	}

	nodeLayer := h.randomLevel()
	_ = g.addNode(nodeLayer) 

	
	if g.entryPoint == -1 {
		g.entryPoint = internalID
		g.maxLayer = nodeLayer
		return saveGraph(g, hnswPath)
	}

	ep := g.entryPoint

	
	for lc := g.maxLayer; lc > nodeLayer; lc-- {
		ep = greedyStep(g, allVecs, vector, ep, lc)
	}

	
	for lc := intMin(nodeLayer, g.maxLayer); lc >= 0; lc-- {
		candidates := beamSearch(g, allVecs, vector, ep, lc, h.EfConstruct)

		mMax := h.maxNeighbours(lc)
		neighbours := pickBest(candidates, mMax)

		for _, nbID := range neighbours {
			g.connect(internalID, nbID, lc)
			g.pruneNeighbours(nbID, lc, mMax, allVecs)
		}

		if len(candidates) > 0 {
			ep = candidates[0].id 
		}
	}

	
	if nodeLayer > g.maxLayer {
		g.maxLayer = nodeLayer
		g.entryPoint = internalID
	}

	return saveGraph(g, hnswPath)
}






type KNNResult struct {
	InternalID int
	Score      float32 
}


func (h *HNSW) KNNSearch(idx *models.Index, query []float32, k int) ([]KNNResult, error) {
	hnswPath := idx.Path() + "/" + config.INDEX_HNSW_FILE

	g, err := loadGraph(hnswPath)
	if err != nil {
		return nil, fmt.Errorf("KNNSearch: %w", err)
	}
	if g.entryPoint == -1 {
		return nil, nil
	}

	allVecs, err := loadAllVectors(idx)
	if err != nil {
		return nil, fmt.Errorf("KNNSearch: %w", err)
	}

	ep := g.entryPoint

	
	for lc := g.maxLayer; lc >= 1; lc-- {
		ep = greedyStep(g, allVecs, query, ep, lc)
	}

	
	ef := h.EfSearch
	if k > ef {
		ef = k
	}
	candidates := beamSearch(g, allVecs, query, ep, 0, ef)

	if len(candidates) > k {
		candidates = candidates[:k]
	}

	results := make([]KNNResult, len(candidates))
	for i, c := range candidates {
		results[i] = KNNResult{InternalID: c.id, Score: c.score}
	}
	return results, nil
}



func (h *HNSW) KNNSearchByExternalID(idx *models.Index, externalID string, k int) ([]KNNResult, error) {
	internalID, err := idx.FetchInternalIdByExternalId(externalID)
	if err != nil || internalID == -1 {
		return nil, fmt.Errorf("KNNSearchByExternalID: %q not found", externalID)
	}
	v, err := idx.FetchVectorByInternalId(internalID)
	if err != nil {
		return nil, err
	}
	return h.KNNSearch(idx, v.Data(), k)
}





type cand struct {
	id    int
	score float32
}


func greedyStep(g *hnswGraph, vecs [][]float32, q []float32, entryID, layer int) int {
	cur := entryID
	best := cosineSim(q, vecs[cur])
	for {
		moved := false
		for _, nid := range g.neighboursAt(cur, layer) {
			s := cosineSim(q, vecs[nid])
			if s > best {
				best, cur, moved = s, nid, true
			}
		}
		if !moved {
			return cur
		}
	}
}


func beamSearch(g *hnswGraph, vecs [][]float32, q []float32, entryID, layer, ef int) []cand {
	visited := make(map[int]bool)
	visited[entryID] = true

	es := cosineSim(q, vecs[entryID])

	
	work := &maxCandHeap{{entryID, es}}
	heap.Init(work)

	
	res := []cand{{entryID, es}}
	worstRes := es

	for work.Len() > 0 {
		best := heap.Pop(work).(cand)

		if best.score < worstRes && len(res) >= ef {
			break
		}

		for _, nid := range g.neighboursAt(best.id, layer) {
			if visited[nid] {
				continue
			}
			visited[nid] = true
			s := cosineSim(q, vecs[nid])

			if len(res) < ef || s > worstRes {
				heap.Push(work, cand{nid, s})
				res = append(res, cand{nid, s})
				if len(res) > ef {
					
					minIdx := 0
					for i, c := range res {
						if c.score < res[minIdx].score {
							minIdx = i
						}
					}
					res[minIdx] = res[len(res)-1]
					res = res[:len(res)-1]
				}
				
				worstRes = res[0].score
				for _, c := range res {
					if c.score < worstRes {
						worstRes = c.score
					}
				}
			}
		}
	}

	sort.Slice(res, func(i, j int) bool { return res[i].score > res[j].score })
	return res
}


func pickBest(candidates []cand, m int) []int {
	n := len(candidates)
	if n > m {
		n = m
	}
	ids := make([]int, n)
	for i := range ids {
		ids[i] = candidates[i].id
	}
	return ids
}





type maxCandHeap []cand

func (h maxCandHeap) Len() int            { return len(h) }
func (h maxCandHeap) Less(i, j int) bool  { return h[i].score > h[j].score }
func (h maxCandHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxCandHeap) Push(x interface{}) { *h = append(*h, x.(cand)) }
func (h *maxCandHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}








func cosineSim(a, b []float32) float32 {
	var dot float32
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}