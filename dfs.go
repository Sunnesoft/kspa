package graph_shortest_paths

import (
	"math"
)

type Dfs struct {
	deepLimit int
	topK      int

	visited []bool
	psa     []float64

	g     *MultiGraph
	edges MEdgeSeq

	pq        PriorityQueue
	maxWeight float64
}

func (st *Dfs) TopKShortestPaths(g *MultiGraph, srcId int, topK int) (res PriorityQueue) {
	st.psa = make([]float64, st.deepLimit+2)
	st.edges = make(MEdgeSeq, st.deepLimit+1)
	st.visited = make([]bool, len(g.Verteces))
	st.g = g
	st.pq = NewPriorityQueue(0, topK)
	st.topK = topK
	st.maxWeight = math.Inf(-1)

	src := g.Verteces[srcId]
	st.processOptEdges(src, src, 0)
	st.initResIfNot()
	res = processOutsideEdges(st.pq, st.deepLimit, topK, false)
	return
}

func (st *Dfs) processOptEdges(src int, target int, level int) {
	if level >= st.deepLimit {
		return
	}

	st.visited[src] = true

	for _, edge := range st.g.Succ(src) {
		st.edges[level] = edge
		weight := st.psa[level] + edge.weight
		st.psa[level+1] = weight

		if target == edge.data.Id2 {
			if weight >= 0 {
				continue
			}

			for i := level + 1; i < st.deepLimit; i++ {
				st.edges[i] = nil
			}

			if st.pq.Len() < st.topK {
				cedges := make(MEdgeSeq, st.deepLimit)
				copy(cedges, st.edges)
				st.pq.Append(cedges, weight)

				if st.pq.Len() == st.topK {
					st.pq.Init()
					st.maxWeight = st.pq[0].priority
				}
				continue
			}

			if weight < st.maxWeight {
				ms, _ := st.pq[0].value.(MEdgeSeq)
				copy(ms, st.edges)
				st.pq.Update(st.pq[0], st.pq[0].value, weight)
				st.maxWeight = st.pq[0].priority
			}

			continue
		}

		if !st.visited[edge.data.Id2] {
			st.processOptEdges(edge.data.Id2, target, level+1)
		}
	}

	st.visited[src] = false
}

func (st *Dfs) initResIfNot() {
	if st.pq.Len() < st.topK {
		st.pq.Init()
	}
}
