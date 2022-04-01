package kspa

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

type Entity struct {
	EntityId string  `json:"EntityId"`
	Id1      int     `json:"Id1"`
	Id2      int     `json:"Id2"`
	Relation float64 `json:"Relation"`
	Id1i     int     `json:"Id1i"`
	Id2i     int     `json:"Id2i"`
}

type EntityRaw struct {
	EntityId string `json:"EntityId"`
	Id1      int    `json:"Id1"`
	Id2      int    `json:"Id2"`
	Relation string `json:"Relation"`
}

type RandomEntitySeqInfo struct {
	VertexCount     int
	EdgesCount      int
	VertexStdFactor int
	RelationMin     float64
	RelationMax     float64
	NoiseMean       float64
	NoiseStdDev     float64
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

func (b *Entity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&EntityRaw{
		EntityId: b.EntityId,
		Id1:      b.Id1,
		Id2:      b.Id2,
		Relation: fmt.Sprintf("%.15f", b.Relation),
	})
}

type EntitySeq []*Entity

func FromJsonFile(fn string) (seq EntitySeq) {
	byteValue, err := LoadText(fn)

	if err != nil {
		log.Fatalln(err)
	}

	seq = make(EntitySeq, 0)
	err = json.Unmarshal(byteValue, &seq)

	if err != nil {
		log.Fatalln(err)
	}

	return seq
}

func ToJsonFile(fn string, seq EntitySeq) {
	data, err := json.MarshalIndent(seq, "", "\t")

	if err != nil {
		log.Fatalln(err)
	}

	if err = WriteText(fn, data); err != nil {
		log.Fatalln(err)
	}
}

func GenerateRandomEntities(c RandomEntitySeqInfo) EntitySeq {
	edges := make(map[string]*Entity)

	mean := float64(rand.Intn(c.VertexCount))
	std := float64(c.VertexCount / c.VertexStdFactor)

	for i := 0; i < c.VertexCount; i++ {
		u := i
		v := u

		for u == v {
			v = int(rand.NormFloat64()*std + mean)
		}

		id := fmt.Sprintf("entity-%d-%d", u, v)
		relation := rand.Float64()*(c.RelationMax-c.RelationMin) + c.RelationMin
		edges[id] = &Entity{EntityId: id, Id1: u, Id2: v, Relation: relation}
	}

	for i := 0; i < c.EdgesCount-c.VertexCount; {
		u := int(rand.NormFloat64()*std + mean)
		v := u

		for u == v {
			v = int(rand.NormFloat64()*std + mean)
		}

		id := fmt.Sprintf("entity-%d-%d", u, v)

		if _, ok := edges[id]; ok {
			continue
		}

		revid := fmt.Sprintf("entity-%d-%d", v, u)
		relation := 0.0

		if revval, ok := edges[revid]; ok {
			sample := rand.NormFloat64()*c.NoiseStdDev + c.NoiseMean
			relation = math.Abs(1.0/revval.Relation + sample)
		} else {
			relation = rand.Float64()*(c.RelationMax-c.RelationMin) + c.RelationMin
		}

		edges[id] = &Entity{EntityId: id, Id1: u, Id2: v, Relation: relation}
		i++
	}

	res := make(EntitySeq, len(edges))

	i := 0
	for _, v := range edges {
		res[i] = v
		i++
	}
	return res
}

func GenerateRandomEntitiesJson(base string, count int, removeOld bool, c RandomEntitySeqInfo) {
	symbols := strings.Split("abcdefgkilepsdmdkslixua_dsa_DsaalDkFiwSkdmAkAlPDjwQmmCUIYJMheelf_sda", "")

	if removeOld {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, 0766)

	for i := 0; i < count; i++ {
		seq := GenerateRandomEntities(c)

		rand.Shuffle(len(symbols), func(i, j int) {
			symbols[i], symbols[j] = symbols[j], symbols[i]
		})

		buf := []byte(strings.Join(symbols, ""))
		h := make([]byte, 32)
		sha3.ShakeSum256(h, buf)
		fn := fmt.Sprintf("%x.json", h)

		p := path.Join(base, fn)
		ToJsonFile(p, seq)
	}
}
