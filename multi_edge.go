package kspa

import (
	"fmt"
	"math"
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
