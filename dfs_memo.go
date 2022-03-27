package graph_shortest_paths

import (
	"math"
)

type TreeNode struct {
	base   *MultiEdge
	src    int
	target int
	level  int
}

func NewTreeNode(base *MultiEdge, src int, target int, level int) *TreeNode {
	return &TreeNode{
		base:   base,
		src:    src,
		target: target,
		level:  level,
	}
}

type DfsMemo struct {
	deepLimit int
	topK      int

	visited []bool
	psa     []float64

	g     *MultiGraph
	edges MEdgeSeq

	memo [][][][]MEdgeSeq
}

func (st *DfsMemo) TopKShortestPaths(g *MultiGraph, srcId int, topK int) (res PriorityQueue) {
	st.psa = make([]float64, st.deepLimit+2)
	st.edges = make(MEdgeSeq, st.deepLimit+1)
	st.visited = make([]bool, len(g.Verteces))
	st.g = g
	st.topK = topK

	src := g.Verteces[srcId]
	st.initMemo()
	st.processOptEdges(src, src, 0)
	pq := st.setRes(src, src)
	res = processOutsideEdges(pq, st.deepLimit, topK, false)
	return
}

func (st *DfsMemo) processOptEdges(src int, target int, level int) (res []MEdgeSeq) {
	if level >= st.deepLimit {
		return
	}

	// if st.visited[src] {
	// 	return
	// }

	if st.memo[src][target][level] != nil {
		return st.memo[src][target][level]
	}

	res = make([]MEdgeSeq, 0)

	// st.visited[src] = true
	succ := st.g.Succ(src)

	for j, edge := range succ {
		// st.edges[level] = edge
		// weight := st.psa[level] + edge.weight
		// st.psa[level+1] = weight

		if target == edge.data.Id2 {
			// if weight >= 0 {
			// 	continue
			// }

			res = append(res, MEdgeSeq{succ[j]})

			continue
		}

		// if !st.visited[edge.data.Id2] {
		loc := st.processOptEdges(edge.data.Id2, target, level+1)

		if len(loc) > 0 {
			tmp := make([]MEdgeSeq, len(loc))

			for i, seq := range loc {
				size := len(seq) + 1
				seqi := make(MEdgeSeq, size)
				seqi[0] = succ[j]
				copy(seqi[1:size], seq)
				tmp[i] = seqi
			}

			res = append(res, tmp...)
		}
		// }
	}

	st.memo[src][target][level] = res
	// st.visited[src] = false
	return
}

func (st *DfsMemo) initMemo() {
	n := len(st.g.Verteces)
	st.memo = make([][][][]MEdgeSeq, n)
	for i := 0; i < n; i++ {
		st.memo[i] = make([][][]MEdgeSeq, n)
		for j := 0; j < n; j++ {
			st.memo[i][j] = make([][]MEdgeSeq, st.deepLimit)
		}
	}
}

func (st *DfsMemo) setRes(src int, target int) (pq PriorityQueue) {
	pq = NewPriorityQueue(0, st.topK)
	maxWeight := math.Inf(-1)

	for level := 0; level < st.deepLimit; level++ {
		paths := st.memo[src][target][level]
		if paths == nil {
			continue
		}

		for _, path := range paths {
			weight := path.GetWeight()

			if weight >= 0.0 {
				continue
			}

			if pq.Len() < st.topK {
				pq.Append(path, weight)

				if pq.Len() == st.topK {
					pq.Init()
					maxWeight = pq[0].priority
				}
				continue
			}

			if weight < maxWeight {
				pq.Update(pq[0], path, weight)
				maxWeight = pq[0].priority
			}
		}

	}

	if pq.Len() < st.topK {
		pq.Init()
	}
	return
}
