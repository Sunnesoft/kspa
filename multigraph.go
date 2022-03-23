package graph_shortest_paths

import (
	"container/list"
	"fmt"
	"math"
)

type Edge interface {
	Update()
	Weight() float64
	U() int
	V() int
}

type SingleEdge struct {
	data   *Entity
	weight float64
}

func (e *SingleEdge) Update() {
	e.weight = Weight(e.data)
}

func (e *SingleEdge) Weight() float64 {
	return e.weight
}

func (e *SingleEdge) U() int {
	return e.data.Id1
}

func (e *SingleEdge) V() int {
	return e.data.Id2
}

type MultiEdge struct {
	SingleEdge
	edges []*SingleEdge
}

func (e *MultiEdge) Update() {
	e.weight = math.MaxFloat64
	for _, n := range e.edges {
		if e.weight > n.weight {
			e.data = n.data
			e.weight = n.weight
		}
	}
}

func (e *MultiEdge) Weight() float64 {
	return e.weight
}

func (e *MultiEdge) U() int {
	return e.data.Id1
}

func (e *MultiEdge) V() int {
	return e.data.Id2
}

func Weight(a *Entity) float64 {
	return -math.Log(a.Relation)
}

type MultiGraph struct {
	Edges    []*MultiEdge
	Verteces map[int]int
	entities []*Entity

	predecessors [][]*MultiEdge
	successors   [][]*MultiEdge

	marked    []bool
	p         [][]*MultiEdge
	pathCount int
}

func (g *MultiGraph) buildVerteces() {
	g.Verteces = make(map[int]int)

	for _, v := range g.entities {
		g.Verteces[v.Id1] = -1
		g.Verteces[v.Id2] = -1
	}

	j := 0
	for i, v := range g.entities {
		if g.Verteces[v.Id1] == -1 {
			g.Verteces[v.Id1] = j
			j++
		}
		g.entities[i].Id1 = g.Verteces[v.Id1]

		if g.Verteces[v.Id2] == -1 {
			g.Verteces[v.Id2] = j
			j++
		}
		g.entities[i].Id2 = g.Verteces[v.Id2]
	}
}

func (g *MultiGraph) getGroupedEdgesById1(bufferSize int) (res [][]*Entity) {
	n := len(g.Verteces)
	res = make([][]*Entity, n)

	for i, v := range g.entities {
		lab := v.Id1
		if res[lab] == nil {
			res[lab] = make([]*Entity, 0, bufferSize)
		}
		res[lab] = append(res[lab], g.entities[i])
	}

	return
}

func (g *MultiGraph) buildEdges(bufferSize int) {
	groupedEdgesById1 := g.getGroupedEdgesById1(bufferSize)
	g.Edges = make([]*MultiEdge, 0, len(groupedEdgesById1))

	for _, d := range groupedEdgesById1 {
		groupedById2 := make(map[int][]*Entity)
		for _, e := range d {
			if _, ok := groupedById2[e.Id2]; !ok {
				groupedById2[e.Id2] = make([]*Entity, 0)
			}
			groupedById2[e.Id2] = append(groupedById2[e.Id2], e)
		}

		for _, edges := range groupedById2 {
			edge := &MultiEdge{edges: make([]*SingleEdge, len(edges))}

			for i, val := range edges {
				singleEdge := &SingleEdge{data: val}
				singleEdge.Update()
				edge.edges[i] = singleEdge
			}

			edge.Update()

			g.Edges = append(g.Edges, edge)
		}
	}
}

func (g *MultiGraph) setAdjacent() {
	n := len(g.Verteces)
	g.predecessors = make([][]*MultiEdge, n)
	g.successors = make([][]*MultiEdge, n)

	for _, v := range g.Edges {
		if g.predecessors[v.V()] == nil {
			g.predecessors[v.V()] = make([]*MultiEdge, 0, 1)
		}

		if g.successors[v.U()] == nil {
			g.successors[v.U()] = make([]*MultiEdge, 0, 1)
		}

		g.predecessors[v.V()] = append(g.predecessors[v.V()], v)
		g.successors[v.U()] = append(g.successors[v.U()], v)
	}
}

func (g *MultiGraph) BuildGraph(ent []Entity) {
	g.entities = make([]*Entity, len(ent))
	for i := range ent {
		g.entities[i] = &ent[i]
	}

	g.buildVerteces()
	g.buildEdges(0)
	g.setAdjacent()

	// for i, v := range g.Edges {
	// 	fmt.Println(i, v.data, v.weight, v.edges)
	// }
}

func (g *MultiGraph) pred(v int) []*MultiEdge {
	return g.predecessors[v]
}

func (g *MultiGraph) succ(u int) []*MultiEdge {
	return g.successors[u]
}

func (g *MultiGraph) dfs(src int, target int, counter int) {
	if counter > 5 {
		return
	}

	g.marked[src] = true

	for _, edge := range g.p[src] {
		if target == edge.data.Id1 {
			// fmt.Println("founded")
			g.pathCount++
			continue
		}

		if !g.marked[edge.data.Id1] {
			g.dfs(edge.data.Id1, target, counter+1)
		}
	}

	g.marked[src] = false
}

func firstIndexOf(vert int, path []int) int {
	for i, v := range path {
		if vert == v {
			return i
		}
	}
	return -1
}

func printCycle(vert int, path []int) {
	print := false
	for _, v := range path {
		if vert == v {
			print = true
		}

		if print {
			fmt.Printf("%d ", v)
		}
	}

	fmt.Println(vert)
}

func (g *MultiGraph) dfsStacked(src int, target int) {
	type vm struct {
		Vert   int
		Deep   int
		Parent int
		Index  int
	}

	vertexCount := len(g.Verteces)

	visited := make([]bool, vertexCount)

	deepLimit := 5
	// bufferFactor := 2

	s := make([]vm, 0) // vertexCount*bufferFactor)
	counter := 0

	s = append(s, vm{src, 1, -1, 0})

	// si := &s[counter]
	// si.Vert = src
	// si.Deep = 1
	// si.Parent = -1
	// si.Index = 0
	counter++

	pathCount := 0

	path := make([]int, deepLimit+1)
	psa := make([]float64, deepLimit+2)
	edges := make([]*MultiEdge, deepLimit+1)
	cycles := make([][][]*MultiEdge, vertexCount)

	var level_m1 int = 0
	var level int = 0
	var node vm

	topK := 100
	pq := NewPriorityQueue(0, topK)
	maxWeight := math.Inf(-1)

	for counter > 0 {
		counter--
		node = s[counter]
		level = node.Deep
		level_m1 = level - 1
		s = s[:counter]

		for i := level_m1; i < deepLimit+1; i++ {
			visited[path[i]] = false
			path[i] = 0
			edges[i] = nil
		}

		if visited[node.Vert] {
			// printCycle(node.Vert, path[:level-1])

			indx := firstIndexOf(node.Vert, path[:level-1])
			edge := g.p[node.Parent][node.Index]
			weight := psa[level-1] + edge.weight - psa[indx+1]

			if weight >= 0.0 {
				continue
			}

			if cycles[path[indx]] == nil {
				cycles[path[indx]] = make([][]*MultiEdge, 0)
			}

			cedges := make([]*MultiEdge, deepLimit)
			edges[level_m1] = edge
			copy(cedges, edges[1:level_m1])

			cycles[path[indx]] = append(cycles[path[indx]], cedges)
			continue
		}

		path[level_m1] = node.Vert
		psa[level] = psa[level_m1]

		if node.Parent != -1 {
			edges[level_m1] = g.p[node.Parent][node.Index]
			psa[level] += edges[level_m1].weight
		}

		visited[node.Vert] = true

		if level > deepLimit {
			continue
		}

		for i, edge := range g.p[node.Vert] {
			u := edge.data.Id1
			if target == u {
				pathCount++

				edges[level] = edge
				weight := psa[level] + edge.weight

				if weight >= 0 {
					continue
				}

				if pq.Len() < topK {
					cedges := make([]*MultiEdge, deepLimit)
					copy(cedges, edges[1:])
					pq.Append(cedges, weight)

					if pq.Len() == topK {
						pq.Init()
						maxWeight = pq[0].priority
					}
					continue
				}

				if weight < maxWeight {
					pq.Update(pq[0], edges[1:], weight)
					maxWeight = pq[0].priority
				}

				continue
			}

			s = append(s, vm{u, level + 1, node.Vert, i})

			// si := &s[counter]
			// si.Vert = u
			// si.Deep = level + 1
			// si.Parent = node.Vert
			// si.Index = i
			counter++
		}
	}

	fmt.Println(pathCount)
	fmt.Println(*(pq[0]), *(pq[topK-1]))
	// fmt.Println(cycles[1])
}

func (g *MultiGraph) Bfs(srcId int) {
	vertexCount := len(g.Verteces)
	// edges := g.Edges
	src := g.Verteces[srcId]

	visited := make([]bool, vertexCount)
	adj := make([][]int, vertexCount)
	pred := make([][]*MultiEdge, vertexCount)

	q := list.New()
	q.PushBack(src)
	visited[src] = true

	for q.Len() != 0 {
		e := q.Front()
		u, _ := e.Value.(int)
		q.Remove(e)

		for _, edge := range g.succ(u) {
			v := edge.V()

			if adj[v] == nil {
				adj[v] = make([]int, vertexCount)
			}

			if pred[v] == nil {
				pred[v] = make([]*MultiEdge, 0)
			}

			if adj[v][u] == 0 {
				adj[v][u]++
				pred[v] = append(pred[v], edge)
			}

			if visited[v] {
				continue
			}

			visited[v] = true
			q.PushBack(v)
		}
	}

	// fmt.Println(pred)

	g.p = pred
	g.marked = make([]bool, vertexCount)
	g.pathCount = 0
	// g.dfs(src, src, 0)
	g.dfsStacked(src, src)

	fmt.Println(g.pathCount)
}

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
			u, v, cost := edge.data.Id1, edge.data.Id2, edge.weight

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
		u, v, cost := edge.data.Id1, edge.data.Id2, edge.weight
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
