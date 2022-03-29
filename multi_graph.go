package graph_shortest_paths

type MultiGraph struct {
	Edges    MEdgeSeq
	Verteces map[int]int
	entities EntitySeq

	predecessors []MEdgeSeq
	successors   []MEdgeSeq
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
		g.entities[i].Id1i = g.Verteces[v.Id1]

		if g.Verteces[v.Id2] == -1 {
			g.Verteces[v.Id2] = j
			j++
		}
		g.entities[i].Id2i = g.Verteces[v.Id2]
	}
}

func (g *MultiGraph) getGroupedEdgesById1(bufferSize int) (res []EntitySeq) {
	n := len(g.Verteces)
	res = make([]EntitySeq, n)

	for i, v := range g.entities {
		lab := v.Id1i
		if res[lab] == nil {
			res[lab] = make(EntitySeq, 0, bufferSize)
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
			if _, ok := groupedById2[e.Id2i]; !ok {
				groupedById2[e.Id2i] = make([]*Entity, 0)
			}
			groupedById2[e.Id2i] = append(groupedById2[e.Id2i], e)
		}

		for _, edges := range groupedById2 {
			edge := &MultiEdge{edges: make(EdgeSeq, len(edges))}

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
	g.predecessors = make([]MEdgeSeq, n)
	g.successors = make([]MEdgeSeq, n)

	for _, v := range g.Edges {
		if g.predecessors[v.data.Id2i] == nil {
			g.predecessors[v.data.Id2i] = make(MEdgeSeq, 0, 1)
		}

		if g.successors[v.data.Id1i] == nil {
			g.successors[v.data.Id1i] = make(MEdgeSeq, 0, 1)
		}

		g.predecessors[v.data.Id2i] = append(g.predecessors[v.data.Id2i], v)
		g.successors[v.data.Id1i] = append(g.successors[v.data.Id1i], v)
	}
}

func (g *MultiGraph) BuildGraph(ent EntitySeq) {
	g.entities = ent
	g.buildVerteces()
	g.buildEdges(0)
	g.setAdjacent()
}

func (g *MultiGraph) Pred(v int) MEdgeSeq {
	return g.predecessors[v]
}

func (g *MultiGraph) Succ(u int) MEdgeSeq {
	return g.successors[u]
}
