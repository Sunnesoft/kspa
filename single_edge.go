package graph_shortest_paths

import (
	"encoding/json"
	"strings"
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

type EdgeSeq []*SingleEdge

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

func (seq EdgeSeq) ReverseEdgeSeq() {
	for i := 0; i < len(seq)/2; i++ {
		j := len(seq) - i - 1
		seq[i], seq[j] = seq[j], seq[i]
	}
}

func (seq EdgeSeq) GetRelation() float64 {
	rel := 1.0
	for i := 0; i < len(seq); i++ {
		if seq[i] == nil {
			break
		}
		rel *= seq[i].data.Relation
	}
	return rel
}

func (seq EdgeSeq) MarshalJSON() ([]byte, error) {
	s := make([]string, 0, len(seq))
	value := 1.0
	var id int64 = -1

	if len(seq) > 0 {
		id = int64(seq[0].data.Id1)
	}

	for _, edge := range seq {
		if edge == nil {
			break
		}
		s = append(s, edge.data.EntityId)
		value *= edge.data.Relation
	}

	chain := strings.Join(s, " -> ")

	return json.Marshal(&struct {
		Id    int64   `json:"id"`
		Chain string  `json:"chain"`
		Value float64 `json:"value"`
	}{
		Id:    id,
		Chain: chain,
		Value: value,
	})
}
