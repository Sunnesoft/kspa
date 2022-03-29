package graph_shortest_paths

import (
	"container/list"
	"math"
)

type DfsStacked struct {
	deepLimit int
	cycles    []map[uint64]MEdgeSeq
}

func (st *DfsStacked) SetDeepLimit(v int) {
	st.deepLimit = v
}

func (st *DfsStacked) TopK(g *MultiGraph, srcId int, targetId int, topK int) (res PriorityQueue) {
	var pq PriorityQueue

	src := g.Verteces[srcId]
	predecessors := removeNotConnectedParts(g, src)
	pq, st.cycles = dfs(g, predecessors, src, src, st.deepLimit, topK)
	res = ProcessOutsideEdges(pq, st.deepLimit, topK, true)
	return
}

func (st *DfsStacked) TopKOneToOne(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	var pq []PriorityQueue

	n := len(srcIds)

	res = make([]PriorityQueue, n)
	srcs := make([]int, n)
	targets := make([]int, n)

	for i := 0; i < n; i++ {
		srcs[i] = g.Verteces[srcIds[i]]
		targets[i] = g.Verteces[targetIds[i]]
	}

	predecessors := removeNotConnectedPartsForMultiSources(g, srcs)
	pq, st.cycles = dfsm(g, predecessors, srcs, targets, st.deepLimit, topK)

	for i := 0; i < n; i++ {
		res[i] = ProcessOutsideEdges(pq[i], st.deepLimit, topK, true)
	}

	return
}

func (st *DfsStacked) TopKOneToMany(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	res = make([]PriorityQueue, 0)

	//TODO

	return
}

func (st *DfsStacked) TopKCycles(topK int) (res []PriorityQueue) {
	return getTopkCyclesByVertex(st.cycles, st.deepLimit, topK)
}

func (st *DfsStacked) BestCycles() (res []EdgeSeq) {
	return getBestCycleByVertex(st.cycles, st.deepLimit)
}

func removeNotConnectedParts(g *MultiGraph, src int) (pred []MEdgeSeq) {
	vertexCount := len(g.Verteces)

	visited := make([]bool, vertexCount)
	adj := make([][]int, vertexCount)
	pred = make([]MEdgeSeq, vertexCount)

	q := list.New()
	q.PushBack(src)
	visited[src] = true

	for q.Len() != 0 {
		e := q.Front()
		u, _ := e.Value.(int)
		q.Remove(e)

		for _, edge := range g.Succ(u) {
			v := edge.V()

			if adj[v] == nil {
				adj[v] = make([]int, vertexCount)
			}

			if pred[v] == nil {
				pred[v] = make(MEdgeSeq, 0)
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
	return
}

func removeNotConnectedPartsForMultiSources(g *MultiGraph, srcs []int) (pred []MEdgeSeq) {
	vertexCount := len(g.Verteces)

	visited := make([]bool, vertexCount)
	adj := make([][]int, vertexCount)
	pred = make([]MEdgeSeq, vertexCount)

	q := list.New()

	for _, src := range srcs {
		q.PushBack(src)
		visited[src] = true
	}

	for q.Len() != 0 {
		e := q.Front()
		u, _ := e.Value.(int)
		q.Remove(e)

		for _, edge := range g.Succ(u) {
			v := edge.V()

			if adj[v] == nil {
				adj[v] = make([]int, vertexCount)
			}

			if pred[v] == nil {
				pred[v] = make(MEdgeSeq, 0)
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
	return
}

func dfs(g *MultiGraph, pred []MEdgeSeq, src int, target int, deepLimit int, topK int) (pq PriorityQueue, cycles []map[uint64]MEdgeSeq) {
	type vm struct {
		Vert   int
		Deep   int
		Parent int
		Index  int
	}

	var level_m1 int = 0
	var level int = 0
	var node vm

	vertexCount := len(g.Verteces)
	visited := make([]bool, vertexCount)
	s := make([]vm, 0)
	path := make([]int, deepLimit+1)
	psa := make([]float64, deepLimit+2)
	edges := make(MEdgeSeq, deepLimit+1)
	cycles = make([]map[uint64]MEdgeSeq, vertexCount)
	pq = NewPriorityQueue(0, topK)

	maxWeight := math.Inf(-1)
	pathCount := 0
	counter := 0

	s = append(s, vm{src, 1, -1, 0})
	counter++

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
			indx := firstIndexOf(node.Vert, path[:level-1])
			edge := pred[node.Parent][node.Index]
			weight := psa[level-1] + edge.weight - psa[indx+1]

			bits := math.Float64bits(weight) >> 2

			if weight >= 0.0 {
				continue
			}

			if cycles[path[indx]] == nil {
				cycles[path[indx]] = make(map[uint64]MEdgeSeq)
			}

			if _, ok := cycles[path[indx]][bits]; !ok {
				cedges := make(MEdgeSeq, level-1)
				edges[level_m1] = edge
				copy(cedges, edges[1:level])

				cycles[path[indx]][bits] = cedges
			}

			continue
		}

		path[level_m1] = node.Vert
		psa[level] = psa[level_m1]

		if node.Parent != -1 {
			edges[level_m1] = pred[node.Parent][node.Index]
			psa[level] += edges[level_m1].weight
		}

		visited[node.Vert] = true

		if level > deepLimit {
			continue
		}

		for i, edge := range pred[node.Vert] {
			u := edge.data.Id1i
			if target == u {
				pathCount++

				edges[level] = edge
				weight := psa[level] + edge.weight

				if weight >= 0 {
					continue
				}

				if pq.Len() < topK {
					cedges := make(MEdgeSeq, deepLimit)
					copy(cedges, edges[1:])
					pq.Append(cedges, weight)

					if pq.Len() == topK {
						pq.Init()
						maxWeight = pq[0].priority
					}
					continue
				}

				if weight < maxWeight {
					ms, _ := pq[0].value.(MEdgeSeq)
					copy(ms, edges[1:])
					pq.Update(pq[0], pq[0].value, weight)
					maxWeight = pq[0].priority
				}

				continue
			}

			s = append(s, vm{u, level + 1, node.Vert, i})
			counter++
		}
	}

	if pq.Len() < topK {
		pq.Init()
	}

	return
}

func dfsm(g *MultiGraph, pred []MEdgeSeq, srcs []int, targets []int, deepLimit int, topK int) (pq []PriorityQueue, cycles []map[uint64]MEdgeSeq) {
	type vm struct {
		Vert   int
		Deep   int
		Parent int
		Index  int
	}

	var level_m1 int = 0
	var level int = 0
	var node vm

	maxWeight := math.Inf(-1)
	pathCount := 0
	counter := 0

	vertexCount := len(g.Verteces)
	visited := make([]bool, vertexCount)
	s := make([]vm, 0)
	path := make([]int, deepLimit+1)
	psa := make([]float64, deepLimit+2)
	edges := make(MEdgeSeq, deepLimit+1)
	cycles = make([]map[uint64]MEdgeSeq, vertexCount)
	n := len(srcs)

	pq = make([]PriorityQueue, n)
	for i := 0; i < n; i++ {
		pq[i] = NewPriorityQueue(0, topK)
		s = append(s, vm{targets[i], 1, -1, 0})
		counter++
	}

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
			indx := firstIndexOf(node.Vert, path[:level-1])
			edge := pred[node.Parent][node.Index]
			weight := psa[level-1] + edge.weight - psa[indx+1]

			bits := math.Float64bits(weight) >> 2

			if weight >= 0.0 {
				continue
			}

			if cycles[path[indx]] == nil {
				cycles[path[indx]] = make(map[uint64]MEdgeSeq)
			}

			if _, ok := cycles[path[indx]][bits]; !ok {
				cedges := make(MEdgeSeq, level-1)
				edges[level_m1] = edge
				copy(cedges, edges[1:level])

				cycles[path[indx]][bits] = cedges
			}

			continue
		}

		path[level_m1] = node.Vert
		psa[level] = psa[level_m1]

		if node.Parent != -1 {
			edges[level_m1] = pred[node.Parent][node.Index]
			psa[level] += edges[level_m1].weight
		}

		visited[node.Vert] = true

		if level > deepLimit {
			continue
		}

		for i, edge := range pred[node.Vert] {
			u := edge.data.Id1i
			curTarget := edges[1].data.Id2i
			srcIndex := -1

			for j, target := range targets {
				if target == curTarget {
					srcIndex = j
					break
				}
			}

			if srcIndex != -1 {
				pathCount++

				edges[level] = edge
				weight := psa[level] + edge.weight

				if weight >= 0 {
					continue
				}

				curpq := pq[srcIndex]

				if curpq.Len() < topK {
					cedges := make(MEdgeSeq, deepLimit)
					copy(cedges, edges[1:])
					curpq.Append(cedges, weight)

					if curpq.Len() == topK {
						curpq.Init()
						maxWeight = curpq[0].priority
					}
					continue
				}

				if weight < maxWeight {
					ms, _ := curpq[0].value.(MEdgeSeq)
					copy(ms, edges[1:])
					curpq.Update(curpq[0], curpq[0].value, weight)
					maxWeight = curpq[0].priority
				}

				continue
			}

			s = append(s, vm{u, level + 1, node.Vert, i})
			counter++
		}
	}

	for i := 0; i < n; i++ {
		if pq[i].Len() < topK {
			pq[i].Init()
		}
	}

	return
}

func getBestCycleByVertex(cycles []map[uint64]MEdgeSeq, deepLimit int) (res []EdgeSeq) {
	res = make([]EdgeSeq, len(cycles))

	for id, cyclePool := range cycles {
		var minWeight float64 = 0.0
		var minKey uint64 = 0

		for key, cycle := range cyclePool {
			weight := 0.0
			for i := 0; i < len(cycle); i++ {
				weight += cycle[i].weight
			}

			if minWeight > weight {
				minWeight = weight
				minKey = key
			}
		}

		if minKey != 0 {
			cycle := cyclePool[minKey]
			path := make(EdgeSeq, deepLimit)

			for i := 0; i < len(cycle); i++ {
				path[i] = &SingleEdge{data: cycle[i].data, weight: cycle[i].weight}
			}
			path[:len(cycle)].ReverseEdgeSeq()
			res[id] = path
		}
	}
	return
}

func getTopkCyclesByVertex(cycles []map[uint64]MEdgeSeq, deepLimit int, topK int) (res []PriorityQueue) {
	res = make([]PriorityQueue, len(cycles))

	mask := make([]int, deepLimit)
	limits := make([]int, deepLimit)
	path := make(EdgeSeq, deepLimit)

	for id, cyclePool := range cycles {
		maxWeight := 0.0
		topCycles := NewPriorityQueue(0, topK)

		for _, cycle := range cyclePool {
			for i := 0; i < deepLimit; i++ {
				path[i] = nil
			}

			weight := 0.0
			seqSize := len(cycle)
			for i := 0; i < seqSize; i++ {
				limits[i] = len(cycle[i].edges)
				path[i] = cycle[i].edges[0]
				weight += path[i].weight
			}

			rem := 0

			for {
				for i := 1; i < seqSize && rem > 0; i++ {
					curEdges := cycle[i].edges
					weight -= curEdges[mask[i]].weight
					mask[i] += rem
					mask[i], rem = mask[i]%limits[i], mask[i]/limits[i]
					path[i] = curEdges[mask[i]]
					weight += path[i].weight
				}

				if rem > 0 {
					break
				}

				if weight <= maxWeight {
					if topCycles.Len() < topK {
						cpath := make(EdgeSeq, deepLimit)
						copy(cpath, path)
						cpath[:seqSize].ReverseEdgeSeq()
						topCycles.Append(cpath, weight)

						if topCycles.Len() == topK {
							topCycles.Init()
							maxWeight = topCycles[0].priority
						}
					} else {
						ms, _ := topCycles[0].value.(EdgeSeq)
						copy(ms, path)
						ms[:seqSize].ReverseEdgeSeq()

						topCycles.Update(topCycles[0], topCycles[0].value, weight)
						maxWeight = topCycles[0].priority
					}
				}

				curEdges := cycle[0].edges
				weight -= curEdges[mask[0]].weight
				mask[0] += 1
				mask[0], rem = mask[0]%limits[0], mask[0]/limits[0]
				path[0] = curEdges[mask[0]]
				weight += path[0].weight
			}
		}

		if topCycles.Len() < topK {
			topCycles.Init()
		}

		res[id] = topCycles
	}

	return
}
