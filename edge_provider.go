package kspa

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

type EdgeConstraint interface {
	*Entity | *SingleEdge
}

func FromJsonFile[T EdgeConstraint](fn string) (seq []T, e error) {
	byteValue, err := LoadText(fn)

	if err != nil {
		return nil, err
	}

	seq = make([]T, 0)
	err = json.Unmarshal(byteValue, &seq)

	if err != nil {
		return nil, err
	}

	return seq, nil
}

func ToJsonFile[T EdgeConstraint](fn string, seq []T) error {
	data, err := json.MarshalIndent(seq, "", "\t")

	if err != nil {
		return err
	}

	if err = WriteText(fn, data); err != nil {
		return err
	}
	return nil
}

func FromCsvFile(fn string) (seq []*SingleEdge, e error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	seq = make([]*SingleEdge, 0)
	e = nil

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")

		if len(tokens) < 4 {
			return nil, fmt.Errorf("incorrect tokens size %d, wants more then 3", len(tokens))
		}

		// id, err := strconv.Atoi(tokens[0])

		// if err != nil {
		// 	return nil, err
		// }

		id := tokens[0]

		tokIn, err := strconv.Atoi(tokens[1])

		if err != nil {
			return nil, err
		}

		tokOut, err := strconv.Atoi(tokens[2])

		if err != nil {
			return nil, err
		}

		relation, err := strconv.ParseFloat(tokens[3], 64)

		if err != nil {
			return nil, err
		}

		edge := &SingleEdge{
			Data: &Entity{
				EntityId: id,
				Id1:      tokIn,
				Id2:      tokOut,
				Relation: relation,
			},
		}
		edge.Update()

		seq = append(seq, edge)
	}

	if err := scanner.Err(); err != nil {
		return seq, err
	}

	return
}

func ToCsvFile(fn string, seq []*SingleEdge) error {
	var b strings.Builder
	for _, p := range seq {
		fmt.Fprintf(&b, "%s,%d,%d,%.15f,\n", p.Data.EntityId, p.Data.Id1, p.Data.Id2, p.Data.Relation)
	}

	bb := []byte(b.String())

	return WriteText(fn, bb)
}

func PriorityQueues2BinaryFile(fp string, pqs []PriorityQueue) error {
	file, _ := os.Create(fp)
	defer file.Close()

	gob.Register(EdgeSeq{})

	enc := gob.NewEncoder(file)
	err := enc.Encode(&pqs)
	return err
}

func BinaryFile2PriorityQueues(fp string) ([]PriorityQueue, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var pqs []PriorityQueue

	enc := gob.NewDecoder(file)
	err = enc.Decode(&pqs)
	return pqs, err
}

func GenerateRandomLimitOrders(fp string, count int, percDiff float64) []*SingleEdge {
	source, _ := FromCsvFile(fp)
	n := len(source)
	res := make([]*SingleEdge, count)

	symbols := strings.Split("abcdefgkilepsdmdkslixua_dsa_DsaalDkFiwSkdmAkAlPDjwQmmCUIYJMheelf_sda", "")

	for i := 0; i < count; i++ {
		index := rand.Intn(n)

		lo := new(SingleEdge)
		*lo = *source[index]
		lo.Data = new(Entity)
		*(lo.Data) = *(source[index].Data)

		lo.UpdateRelation(lo.Data.Relation * (1.0 + 0.01*(2.0*rand.Float64()*percDiff-percDiff)))

		rand.Shuffle(len(symbols), func(i, j int) {
			symbols[i], symbols[j] = symbols[j], symbols[i]
		})

		lo.Data.EntityId = strings.Join(symbols[0:6], "") + lo.Data.EntityId
		res[i] = lo
	}

	return res
}

type RandomEdgeSeqInfo struct {
	Count    int
	PercDiff float64
}

func GenerateRandomLimitOrdersCsv(base string, tpl string, count int, removeOld bool, c RandomEdgeSeqInfo) {
	symbols := strings.Split("abcdefgkilepsdmdkslixua_dsa_DsaalDkFiwSkdmAkAlPDjwQmmCUIYJMheelf_sda", "")

	if removeOld {
		os.RemoveAll(base)
	}

	err := os.MkdirAll(base, 0766)

	if err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		seq := GenerateRandomLimitOrders(tpl, c.Count, c.PercDiff)

		rand.Shuffle(len(symbols), func(i, j int) {
			symbols[i], symbols[j] = symbols[j], symbols[i]
		})

		buf := []byte(strings.Join(symbols, ""))
		h := make([]byte, 32)
		sha3.ShakeSum256(h, buf)
		fn := fmt.Sprintf("%x.csv", h)

		p := path.Join(base, fn)
		err = ToCsvFile(p, seq)

		if err != nil {
			panic(err)
		}
	}
}
