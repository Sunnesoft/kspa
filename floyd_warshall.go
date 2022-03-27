package graph_shortest_paths

import (
	"fmt"
	"math"
)

func (g *MultiGraph) FloydWarshall(srcId int) {
	vertexCount := len(g.Verteces)
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
		d[edge.U()][edge.V()] = edge.Weight()
		p[edge.U()][edge.V()] = edge.V()
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

	src := g.Verteces[srcId]

	path := make([]int, 0, vertexCount)
	path = append(path, src)
	prior := src
	prior = p[prior][src]

	for prior != src {
		path = append(path, prior)
		prior = p[prior][src]
	}

	fmt.Println(path)
}
