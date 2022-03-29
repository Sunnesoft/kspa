package graph_shortest_paths

import "math"

var MIN_WEIGHT float64 = math.Inf(-1)
var MAX_WEIGHT float64 = math.Inf(1)

type Dfs struct {
	g *MultiGraph

	deepLimit int
	edges     MEdgeSeq
	psa       []float64

	topK      int
	pq        PriorityQueue
	maxWeight float64
}

func (st *Dfs) SetDeepLimit(v int) {
	st.deepLimit = v
	st.psa = nil
	st.edges = nil
	st.psa = make([]float64, st.deepLimit+2)
	st.edges = make(MEdgeSeq, st.deepLimit+1)
}

func (st *Dfs) SetTopKValue(v int) {
	st.topK = v
	st.pq = nil
	st.pq = NewPriorityQueue(0, v)
	st.maxWeight = MIN_WEIGHT
}
