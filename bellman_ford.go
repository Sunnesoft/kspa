package kspa

import (
	"fmt"
	"math"
)

type BellmanFord struct {
	deepLimit   int
	uniquePaths bool
}

func (st *BellmanFord) TopK(g *MultiGraph, srcId int, targetId int, topK int) (res PriorityQueue) {
	if topK != 1 {
		panic(fmt.Errorf("BellmanFord.TopK doesn't support several paths searching"))
	}

	vertexCount := len(g.VertexIndex)
	edges := g.Edges
	src := g.VertexIndex[srcId]

	d := make([]float64, vertexCount)
	p := make([]int, vertexCount)

	for i := 0; i < vertexCount; i++ {
		p[i] = -1
		d[i] = math.Inf(0)
	}

	d[src] = 0

	if st.deepLimit <= 0 || st.deepLimit > vertexCount-1 {
		st.deepLimit = vertexCount - 1
	}

	for i := 1; i < st.deepLimit+1; i++ {
		for _, edge := range edges {
			if d[edge.Data.Id1i]+edge.Weight < d[edge.Data.Id2i] {
				d[edge.Data.Id2i] = d[edge.Data.Id1i] + edge.Weight
				p[edge.Data.Id2i] = edge.Data.Id1i
			}
		}
	}

	res = NewPriorityQueue(0, topK)
	visited := make([]bool, vertexCount)

	for _, edge := range edges {
		if visited[edge.Data.Id2i] {
			continue
		}

		if d[edge.Data.Id1i]+edge.Weight < d[edge.Data.Id2i] {
			cycle := traceNegativeCycle(edge.Data.Id2i, p, st.deepLimit, st.uniquePaths, visited)

			if cycle == nil {
				continue
			}

			//TODO
			//replace path by MultiEdge sequence and process weight
			res.Append(cycle, 0.0)
		}
	}
	return
}

func (st *BellmanFord) TopKOneToOne(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	panic(fmt.Errorf("BellmanFord.TopKOneToOne is not provided"))
}

func (st *BellmanFord) TopKOneToMany(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	panic(fmt.Errorf("BellmanFord.TopKOneToMany is not provided"))
}
