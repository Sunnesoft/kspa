package kspa

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

func (e *SingleEdge) UpdateRelation(relation float64) {
	e.data.Relation = relation
	e.Update()
}

func (e *SingleEdge) Weight() float64 {
	return e.weight
}

func (e *SingleEdge) U() int {
	return e.data.Id1i
}

func (e *SingleEdge) V() int {
	return e.data.Id2i
}

func (b *SingleEdge) UnmarshalJSON(data []byte) error {
	var v Entity
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	b.data = &v
	b.Update()
	return nil
}

func (b *SingleEdge) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.data)
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

func (seq *EdgeSeq) BuildVertexIndex() map[int]int {
	vertexIndex := make(map[int]int)

	for _, v := range *seq {
		vertexIndex[v.data.Id1] = -1
		vertexIndex[v.data.Id2] = -1
	}

	j := 0
	for i, v := range *seq {
		if vertexIndex[v.data.Id1] == -1 {
			vertexIndex[v.data.Id1] = j
			j++
		}
		(*seq)[i].data.Id1i = vertexIndex[v.data.Id1]

		if vertexIndex[v.data.Id2] == -1 {
			vertexIndex[v.data.Id2] = j
			j++
		}
		(*seq)[i].data.Id2i = vertexIndex[v.data.Id2]
	}
	return vertexIndex
}

func (seq *EdgeSeq) SetVertexIndex(vertexIndex map[int]int) {
	for i, v := range *seq {
		(*seq)[i].data.Id1i = vertexIndex[v.data.Id1]
		(*seq)[i].data.Id2i = vertexIndex[v.data.Id2]
	}
}

func (seq EdgeSeq) MarshalJSON() ([]byte, error) {
	s := make([]string, 0, len(seq))
	value := 1.0
	var in int64 = -1
	var out int64 = -1

	if len(seq) > 0 {
		in = int64(seq[0].data.Id1)
	}

	for _, edge := range seq {
		if edge == nil {
			break
		}
		s = append(s, edge.data.EntityId)
		value *= edge.data.Relation
		out = int64(edge.data.Id2)
	}

	chain := strings.Join(s, " -> ")

	return json.Marshal(&struct {
		In    int64   `json:"in"`
		Out   int64   `json:"out"`
		Chain string  `json:"chain"`
		Value float64 `json:"value"`
	}{
		In:    in,
		Out:   out,
		Chain: chain,
		Value: value,
	})
}
