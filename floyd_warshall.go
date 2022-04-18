package kspa

import (
	"fmt"
	"math"
)

type FloydWarshall struct{}

func (st *FloydWarshall) TopK(g *MultiGraph, srcId int, targetId int, topK int) (res PriorityQueue) {
	if topK != 1 {
		panic(fmt.Errorf("FloydWarshall.TopK doesn't support several paths searching"))
	}

	vertexCount := len(g.VertexIndex)
	edges := g.Edges

	d := make([][]float64, vertexCount)
	p := make([][]int, vertexCount)

	infinity := math.Inf(0)

	for i := 0; i < vertexCount; i++ {
		d[i] = make([]float64, vertexCount)
		p[i] = make([]int, vertexCount)

		for j := 0; j < vertexCount; j++ {
			p[i][j] = -1

			if i == j {
				d[i][j] = 0.0
				p[i][j] = j
			} else {
				d[i][j] = infinity
			}
		}
	}

	for _, edge := range edges {
		d[edge.Data.Id1i][edge.Data.Id2i] = edge.Weight
		p[edge.Data.Id1i][edge.Data.Id2i] = edge.Data.Id2i
	}

	for k := 0; k < vertexCount; k++ {
		for i := 0; i < vertexCount; i++ {
			if d[i][k] < infinity {
				for j := 0; j < vertexCount; j++ {
					if d[k][j] < infinity {
						temp := d[i][k] + d[k][j]

						if d[i][j] == infinity || temp < d[i][j] {
							d[i][j] = temp
							p[i][j] = p[i][k]
						}
					}
				}
			}
		}
	}

	src := g.VertexIndex[srcId]

	path := make([]int, 0, vertexCount)
	path = append(path, src)
	prior := src
	prior = p[prior][src]

	for prior != src {
		path = append(path, prior)
		prior = p[prior][src]
	}

	fmt.Println(path)

	//TODO
	//filling result from path array

	res = NewPriorityQueue(0, 1)
	return
}

func (st *FloydWarshall) TopKOneToOne(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	panic(fmt.Errorf("FloydWarshall.TopKOneToOne is not provided"))
}

func (st *FloydWarshall) TopKOneToMany(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	panic(fmt.Errorf("FloydWarshall.TopKOneToMany is not provided"))
}
