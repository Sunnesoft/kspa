package graph_shortest_paths

import (
	"encoding/json"
	"log"
	"strconv"
)

type Entity struct {
	EntityId string  `json:"EntityId"`
	Id1      int     `json:"Id1"`
	Id2      int     `json:"Id2"`
	Relation float64 `json:"Relation"`
}

type EntityRaw struct {
	EntityId string `json:"EntityId"`
	Id1      int    `json:"Id1"`
	Id2      int    `json:"Id2"`
	Relation string `json:"Relation"`
}

func (b *Entity) UnmarshalJSON(data []byte) error {
	var v EntityRaw
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	b.EntityId = v.EntityId
	b.Id1 = v.Id1 // fmt.Sprint(v.Id1)
	b.Id2 = v.Id2 // fmt.Sprint(v.Id2)

	if s, err := strconv.ParseFloat(v.Relation, 64); err == nil {
		b.Relation = s
	}

	return nil
}

type EntitySeq []*Entity

func FromJsonFile(fn string) (seq EntitySeq) {
	byteValue, err := LoadText(fn)

	if err != nil {
		log.Fatalln(err)
	}

	seq = make(EntitySeq, 0)
	json.Unmarshal(byteValue, &seq)
	return seq
}
