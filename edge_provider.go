package kspa

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

		if len(tokens) != 4 {
			return nil, fmt.Errorf("incorrect tokens size %d, wants 4", len(tokens))
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

		relation, err := strconv.ParseFloat(tokens[2], 64)

		if err != nil {
			return nil, err
		}

		edge := &SingleEdge{
			data: &Entity{
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
