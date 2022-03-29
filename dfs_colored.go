package graph_shortest_paths

type DfsColored struct {
	Dfs
	visited []bool
}

func (st *DfsColored) TopK(g *MultiGraph, srcId int, targetId int, topK int) (res PriorityQueue) {
	st.visited = make([]bool, len(g.Verteces))
	st.g = g
	st.SetTopKValue(topK)

	src := g.Verteces[srcId]
	target := g.Verteces[targetId]
	st.processOptEdges(src, target, 0)
	st.initResIfNot()
	res = ProcessOutsideEdges(st.pq, st.deepLimit, topK, false)
	return
}

func (st *DfsColored) TopKOneToOne(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	st.visited = make([]bool, len(g.Verteces))
	st.g = g
	st.SetTopKValue(topK)
	res = make([]PriorityQueue, 0)

	//TODO

	return
}

func (st *DfsColored) TopKOneToMany(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	st.visited = make([]bool, len(g.Verteces))
	st.g = g
	st.SetTopKValue(topK)
	res = make([]PriorityQueue, 0)

	//TODO

	return
}

func (st *DfsColored) processOptEdges(src int, target int, level int) {
	if level >= st.deepLimit {
		return
	}

	st.visited[src] = true

	for _, edge := range st.g.Succ(src) {
		st.edges[level] = edge
		weight := st.psa[level] + edge.weight
		st.psa[level+1] = weight
		v := edge.data.Id2i

		if target == v {
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

		if !st.visited[v] {
			st.processOptEdges(v, target, level+1)
		}
	}

	st.visited[src] = false
}

func (st *DfsColored) initResIfNot() {
	if st.pq.Len() < st.topK {
		st.pq.Init()
	}
}
