package kspa

import (
	"fmt"
	"sort"
)

type MultiGraph struct {
	Edges       MEdgeSeq
	EdgeIndex   map[uint64]int
	VertexIndex map[int]int

	entities     EdgeSeq
	predecessors []MEdgeSeq
	successors   []MEdgeSeq
}

func (g *MultiGraph) Build(ent EdgeSeq) {
	g.entities = ent
	g.buildVertexIndex()
	g.buildEdges(0)
	g.setAdjacent()
	g.buildEdgeIndex()
}

func (g *MultiGraph) Pred(v int) MEdgeSeq {
	return g.predecessors[v]
}

func (g *MultiGraph) Succ(u int) MEdgeSeq {
	return g.successors[u]
}

func (g *MultiGraph) Add(edges EdgeSeq, status int) (MEdgeSeq, error) {
	edges.SetVertexIndex(g.VertexIndex)
	medges := SingleToMultiEdges(edges)

	for _, medge := range medges {
		if index, ok := g.GetEdgeIndex(medge.Data.Id1, medge.Data.Id2); ok {
			g.Edges[index].MergeWithoutUniqueChecking(medge)
			g.Edges[index].status = status
			continue
		}

		medge.status = status
		g.Edges = append(g.Edges, medge)
		label, err := IdsHash(medge.Data.Id1, medge.Data.Id2)

		if err != nil {
			return nil, err
		}

		g.EdgeIndex[label] = len(g.EdgeIndex)

		if g.successors[medge.Data.Id1i] == nil {
			g.successors[medge.Data.Id1i] = make(MEdgeSeq, 0, 1)
		}

		g.successors[medge.Data.Id1i] = append(g.successors[medge.Data.Id1i], medge)
	}
	return medges, nil
}

func (g *MultiGraph) Remove(medges MEdgeSeq, status int) {
	indeces := make([]int, 0, len(medges))

	for _, medge := range medges {
		if index, ok := g.GetEdgeIndex(medge.Data.Id1, medge.Data.Id2); ok {
			if !g.Edges[index].RemoveMany(medge) {
				g.Edges[index].status = status
			}

			if g.Edges[index].Len() == 0 {
				indeces = append(indeces, index)

				for i, suc := range g.successors[medge.Data.Id1i] {
					if suc == medge {
						g.successors = append(g.successors[:i], g.successors[i+1:]...)
						break
					}
				}
			}
		}
	}

	if len(indeces) > 0 {
		sort.Slice(indeces, func(i, j int) bool {
			return indeces[j] < indeces[i]
		})

		for _, index := range indeces {
			g.Edges = append(g.Edges[:index], g.Edges[index+1:]...)
		}

		g.buildEdgeIndex()
	}
}

func (g *MultiGraph) UpdateRelation(ent EdgeSeq) error {
	for _, entity := range ent {
		label, err := IdsHash(entity.Data.Id1, entity.Data.Id2)

		if err != nil {
			return err
		}

		index, ok := g.EdgeIndex[label]

		if !ok {
			return fmt.Errorf("graph structure was changed, please use MultiGraph.Build")
		}

		err = g.Edges[index].UpdateRelation(entity.Data.EntityId, entity.Data.Relation)

		if err != nil {
			return err
		}
	}
	return nil
}

func (g *MultiGraph) GetEdgeIndex(id1, id2 int) (int, bool) {
	label, err := IdsHash(id1, id2)

	if err != nil {
		return 0, false
	}

	index, ok := g.EdgeIndex[label]
	return index, ok
}

func (g *MultiGraph) buildVertexIndex() {
	g.VertexIndex = g.entities.BuildVertexIndex()
}

func (g *MultiGraph) buildEdgeIndex() {
	g.EdgeIndex = nil
	g.EdgeIndex = make(map[uint64]int)
	for i, edge := range g.Edges {
		label, err := IdsHash(edge.Data.Id1, edge.Data.Id2)

		if err != nil {
			panic(err)
		}

		g.EdgeIndex[label] = i
	}
}

func (g *MultiGraph) buildEdges(bufferSize int) {
	g.Edges = SingleToMultiEdges(g.entities)
}

func (g *MultiGraph) setAdjacent() {
	n := len(g.VertexIndex)
	g.predecessors = make([]MEdgeSeq, n)
	g.successors = make([]MEdgeSeq, n)

	for _, v := range g.Edges {
		if g.predecessors[v.Data.Id2i] == nil {
			g.predecessors[v.Data.Id2i] = make(MEdgeSeq, 0, 1)
		}

		if g.successors[v.Data.Id1i] == nil {
			g.successors[v.Data.Id1i] = make(MEdgeSeq, 0, 1)
		}

		g.predecessors[v.Data.Id2i] = append(g.predecessors[v.Data.Id2i], v)
		g.successors[v.Data.Id1i] = append(g.successors[v.Data.Id1i], v)
	}
}
