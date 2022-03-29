package graph_shortest_paths

import (
	"container/list"
	"fmt"
	"math"
)

func (g *MultiGraph) BellmanFord(srcId int, deepLimit int, uniquePaths bool) *list.List {
	vertexCount := len(g.Verteces)
	edges := g.Edges
	src := g.Verteces[srcId]

	d := make([][]float64, vertexCount)
	p := make([][]int, vertexCount)

	for i := 0; i < vertexCount; i++ {
		d[i] = make([]float64, vertexCount)
		p[i] = make([]int, vertexCount)

		for j := 0; j < vertexCount; j++ {
			p[i][j] = -1
			d[i][j] = math.Inf(0)
		}
	}

	d[0][src] = 0

	if deepLimit <= 0 || deepLimit > vertexCount-1 {
		deepLimit = vertexCount - 1
	}

	for i := 1; i < deepLimit+1; i++ {
		for _, edge := range edges {
			fmt.Println(edge.V())
			if d[i-1][edge.U()]+edge.Weight() < d[i][edge.V()] {
				d[i][edge.V()] = d[i-1][edge.U()] + edge.Weight()
				p[i][edge.V()] = edge.U()
			}
		}
	}

	res := list.New()
	// visited := make([]bool, vertexCount)

	for _, edge := range edges {
		// if visited[edge.V()] {
		// 	continue
		// }

		if d[deepLimit][edge.U()]+edge.Weight() < d[deepLimit][edge.V()] {
			v := edge.V()
			for i := 0; i < deepLimit; i++ {
				v = p[deepLimit][v]
			}

			u := v
			path := make([]int, 0, deepLimit)

			for u != p[deepLimit][v] {
				path = append(path, v)
				v = p[deepLimit][v]
			}
			path = reverse(path)

			fmt.Println(edge.V(), path)

			// cycle := traceNegativeCycle(edge.V(), p, deepLimit, uniquePaths, visited)

			// if cycle == nil {
			// 	continue
			// }

			// if cycle[0] == src {
			// 	fmt.Println(cycle)
			// }

			//res.PushBack(cycle)
		}
	}

	return res
}

func (g *MultiGraph) Custom(srcId int) {
	vertexCount := len(g.Verteces)
	edges := g.Edges
	src := g.Verteces[srcId]

	du := make([]float64, vertexCount)
	dv := make([]float64, vertexCount)
	pu := make([]int, vertexCount)
	pv := make([]int, vertexCount)

	infinity := math.Inf(0)

	for i := 0; i < vertexCount; i++ {
		du[i] = infinity
		dv[i] = infinity

		pu[i] = -1
		pv[i] = -1
	}

	du[src] = 0
	dv[src] = 0

	n := len(edges)

	deepLimit := 3 // vertexCount

	fmt.Println(pu, du)

	for i := 1; i < deepLimit; i++ {
		for j := 0; j < n; j++ {
			edge := edges[j]
			u, v, cost := edge.data.Id1i, edge.data.Id2i, edge.weight

			if du[u] != infinity {
				if duc := du[u] + cost; duc < du[v] {
					du[v] = duc
					pu[v] = u
				}
			}

			if du[v] != infinity {
				if duv := dv[v] + cost; duv < dv[u] {
					dv[u] = duv
					pv[u] = v
				}
			}
		}

		fmt.Println(pu, du)
	}

	d := infinity
	var cedge *MultiEdge

	for _, edge := range edges {
		u, v, cost := edge.data.Id1i, edge.data.Id2i, edge.weight
		if du[u] != infinity && dv[v] != infinity && du[u]+cost+dv[v] < d {
			d = du[u] + cost + dv[v]
			cedge = edge
		}
	}

	fmt.Println(pu)
	fmt.Println(pv)

	if cedge != nil {
		left := tracePath(cedge.U(), deepLimit, pu, true)
		right := tracePath(cedge.V(), deepLimit, pv, false)

		fmt.Println(left, right, cedge.U(), cedge.V())
	}
}

func tracePath(v int, deepLimit int, p []int, rev bool) []int {
	path := make([]int, 0, deepLimit)
	path = append(path, v)

	prior := v
	for {
		prior = p[prior]
		for i := 0; i < len(path); i++ {
			if prior == path[i] {
				path = path[i:]
				path = append(path, prior)

				if rev {
					path = reverse(path)
				}

				return path
			}
		}

		if prior == -1 {
			if rev {
				path = reverse(path)
			}

			return path
		}

		path = append(path, prior)
	}
}

func reverse(numbers []int) []int {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func traceNegativeCycle(start int, predecessors []int, deepLimit int, uniquePaths bool, visited []bool) []int {
	path := make([]int, 0, deepLimit)
	path = append(path, start)

	prior := start
	for {
		prior = predecessors[prior]
		for i := 0; i < len(path); i++ {
			if prior == path[i] {
				path = path[i:]
				path = append(path, prior)
				path = reverse(path)
				return path
			}
		}

		// if uniquePaths && visited[prior] {
		// 	return nil
		// }

		path = append(path, prior)
		// visited[prior] = true
	}
}
