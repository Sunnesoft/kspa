package kspa

import (
	"fmt"
	"math"
	"sort"
)

type MultiEdge struct {
	SingleEdge
	edges EdgeSeq
	index map[string]int
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
	return e.data.Id1i
}

func (e *MultiEdge) V() int {
	return e.data.Id2i
}

func (e *MultiEdge) BuildIndex() {
	e.index = nil
	e.index = make(map[string]int)

	for i, edge := range e.edges {
		e.index[edge.data.EntityId] = i
	}
}

func (e *MultiEdge) Len() int {
	return len(e.edges)
}

func (e *MultiEdge) UpdateRelation(entityId string, relation float64) error {
	index, ok := e.index[entityId]

	if !ok {
		return fmt.Errorf("multi-edge structure was changed, please run MultiEdge.BuildIndex")
	}

	edge := e.edges[index]
	edge.UpdateRelation(relation)

	if e.weight > edge.weight {
		e.weight = edge.weight
		e.data = edge.data
	}
	return nil
}

func (e *MultiEdge) Add(s *SingleEdge) error {
	if _, ok := e.index[s.data.EntityId]; ok {
		return e.UpdateRelation(s.data.EntityId, s.data.Relation)
	}

	e.edges = append(e.edges, s)
	e.index[s.data.EntityId] = len(e.edges) - 1

	if e.weight > s.weight {
		e.weight = s.weight
		e.data = s.data
	}
	return nil
}

func (e *MultiEdge) AddMany(s EdgeSeq) error {
	for _, si := range s {
		if err := e.Add(si); err != nil {
			return err
		}
	}
	return nil
}

func (e *MultiEdge) AddManyWithoutUniqueChecking(s EdgeSeq) {
	nextIndex := len(e.edges)
	e.edges = append(e.edges, s...)

	for i := nextIndex; i < len(e.edges); i++ {
		si := s[i-nextIndex]
		e.index[si.data.EntityId] = i

		if e.weight > si.weight {
			e.weight = si.weight
			e.data = si.data
		}
	}
}

func (e *MultiEdge) MergeWithoutUniqueChecking(m *MultiEdge) {
	nextIndex := len(e.edges)
	e.edges = append(e.edges, m.edges...)

	if e.weight > m.weight {
		e.data = m.data
		e.weight = m.weight
	}

	for i, edge := range m.edges {
		e.index[edge.data.EntityId] = nextIndex + i
	}
}

func (e *MultiEdge) Remove(entityId string) error {
	index, ok := e.index[entityId]

	if !ok {
		return fmt.Errorf("multi-edge struct hasn't element with Id %s", entityId)
	}

	e.edges = append(e.edges[:index], e.edges[index+1:]...)
	e.index = nil
	e.index = make(map[string]int)
	e.weight = math.MaxFloat64

	for i, edge := range e.edges {
		if e.weight > edge.weight {
			e.data = edge.data
			e.weight = edge.weight
		}
		e.index[edge.data.EntityId] = i
	}

	return nil
}

func (e *MultiEdge) RemoveManyByIds(entityIds []string) error {
	mask := make([]bool, len(e.edges))

	for _, entityId := range entityIds {
		index, ok := e.index[entityId]

		if !ok {
			return fmt.Errorf("multi-edge struct hasn't element with Id %s", entityId)
		}

		mask[index] = true
	}

	newEdges := make([]*SingleEdge, len(e.edges)-len(entityIds))
	e.weight = math.MaxFloat64
	e.index = nil
	e.index = make(map[string]int)

	i := 0
	for j, cond := range mask {
		if cond {
			continue
		}

		edge := e.edges[j]

		if e.weight > edge.weight {
			e.data = edge.data
			e.weight = edge.weight
		}

		e.index[edge.data.EntityId] = i
		newEdges[i] = edge
		i++
	}

	e.edges = nil
	e.edges = newEdges

	return nil
}

func (e *MultiEdge) RemoveMany(m *MultiEdge) {
	indeces := make([]int, 0, len(m.edges))
	for _, edge := range m.edges {
		if index, ok := e.index[edge.data.EntityId]; ok {
			indeces = append(indeces, index)
		}
	}

	sort.Slice(indeces, func(i, j int) bool {
		return indeces[j] < indeces[i]
	})

	for _, index := range indeces {
		e.edges = append(e.edges[:index], e.edges[index+1:]...)
	}

	e.weight = math.MaxFloat64
	e.index = nil
	e.index = make(map[string]int)

	for i, edge := range e.edges {
		if e.weight > edge.weight {
			e.data = edge.data
			e.weight = edge.weight
		}
		e.index[edge.data.EntityId] = i
	}
}

type MEdgeSeq []*MultiEdge

func (seq MEdgeSeq) ReverseEdgeSeq() {
	for i := 0; i < len(seq)/2; i++ {
		j := len(seq) - i - 1
		seq[i], seq[j] = seq[j], seq[i]
	}
}

func (seq MEdgeSeq) GetWeight() float64 {
	rel := 0.0
	for i := 0; i < len(seq); i++ {
		if seq[i] == nil {
			break
		}
		rel += seq[i].weight
	}
	return rel
}

func SingleToMultiEdges(entities EdgeSeq) MEdgeSeq {
	groupedEdgesById1 := make(map[int]EdgeSeq)

	for i, v := range entities {
		lab := v.data.Id1i
		if groupedEdgesById1[lab] == nil {
			groupedEdgesById1[lab] = make(EdgeSeq, 0)
		}
		groupedEdgesById1[lab] = append(groupedEdgesById1[lab], entities[i])
	}

	res := make(MEdgeSeq, 0, len(groupedEdgesById1))

	for _, d := range groupedEdgesById1 {
		groupedById2 := make(map[int]EdgeSeq)
		for _, e := range d {
			if _, ok := groupedById2[e.data.Id2i]; !ok {
				groupedById2[e.data.Id2i] = make(EdgeSeq, 0)
			}
			groupedById2[e.data.Id2i] = append(groupedById2[e.data.Id2i], e)
		}

		for _, edges := range groupedById2 {
			edge := &MultiEdge{edges: edges}
			edge.Update()
			edge.BuildIndex()

			res = append(res, edge)
		}
	}
	return res
}
