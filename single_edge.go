package kspa

import (
	"encoding/json"
	"strings"
)

type SingleEdge struct {
	Data   *Entity
	Weight float64
	Status int
}

type EdgeSeq []*SingleEdge

func (e *SingleEdge) Update() {
	e.Weight = Weight(e.Data)
}

func (e *SingleEdge) UpdateRelation(relation float64) {
	e.Data.Relation = relation
	e.Update()
}

func (b *SingleEdge) UnmarshalJSON(data []byte) error {
	var v Entity
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	b.Data = &v
	b.Update()
	return nil
}

func (b *SingleEdge) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Data)
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
		rel *= seq[i].Data.Relation
	}
	return rel
}

func (seq *EdgeSeq) BuildVertexIndex() map[int]int {
	vertexIndex := make(map[int]int)

	for _, v := range *seq {
		vertexIndex[v.Data.Id1] = -1
		vertexIndex[v.Data.Id2] = -1
	}

	j := 0
	for i, v := range *seq {
		if vertexIndex[v.Data.Id1] == -1 {
			vertexIndex[v.Data.Id1] = j
			j++
		}
		(*seq)[i].Data.Id1i = vertexIndex[v.Data.Id1]

		if vertexIndex[v.Data.Id2] == -1 {
			vertexIndex[v.Data.Id2] = j
			j++
		}
		(*seq)[i].Data.Id2i = vertexIndex[v.Data.Id2]
	}
	return vertexIndex
}

func (seq *EdgeSeq) SetVertexIndex(vertexIndex map[int]int) {
	for i, v := range *seq {
		(*seq)[i].Data.Id1i = vertexIndex[v.Data.Id1]
		(*seq)[i].Data.Id2i = vertexIndex[v.Data.Id2]
	}
}

func (seq EdgeSeq) MarshalJSON() ([]byte, error) {
	return json.Marshal(EdgeSeqToChainView(seq))
}

type ChainView struct {
	In    int64   `json:"in"`
	Out   int64   `json:"out"`
	Chain string  `json:"chain"`
	Value float64 `json:"value"`
}

func EdgeSeqToChainView(seq EdgeSeq) *ChainView {
	s := make([]string, 0, len(seq))
	value := 1.0
	var in int64 = -1
	var out int64 = -1

	if len(seq) > 0 {
		in = int64(seq[0].Data.Id1)
	}

	for _, edge := range seq {
		if edge == nil {
			break
		}
		s = append(s, edge.Data.EntityId)
		value *= edge.Data.Relation
		out = int64(edge.Data.Id2)
	}

	chain := strings.Join(s, " -> ")

	return &ChainView{
		In:    in,
		Out:   out,
		Chain: chain,
		Value: value,
	}
}
