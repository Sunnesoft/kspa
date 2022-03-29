package graph_shortest_paths

import (
	"math"
)

type MultiEdge struct {
	SingleEdge
	edges EdgeSeq
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

func Weight(a *Entity) float64 {
	return -math.Log(a.Relation)
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
