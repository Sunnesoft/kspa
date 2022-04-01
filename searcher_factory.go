package kspa

import (
	"encoding/json"
	"fmt"
)

func NewDfs(name string, deepLimit int) (Searcher, error) {
	switch name {
	case "colored":
		st := &DfsColored{}
		st.SetDeepLimit(deepLimit)
		return st, nil
	case "memo":
		st := &DfsMemo{}
		st.Init()
		st.SetDeepLimit(deepLimit)
		return st, nil
	case "stacked":
		st := &DfsStacked{}
		st.SetDeepLimit(deepLimit)
		return st, nil
	default:
		return nil, fmt.Errorf("NewDfs: invalid option name %s", name)
	}
}

func DfsDo(st Searcher, op string, g *MultiGraph, srcIds []int, targetIds []int, topK int) (pathsb []byte, err error) {
	switch op {
	case "TopK":
		paths := st.TopK(g, srcIds[0], targetIds[0], topK)
		pathsr := PriorityQueue2SortedArray(paths, false)
		pathsb, err = json.MarshalIndent(pathsr, "", "\t")
		return
	case "TopKOneToOne":
		paths := st.TopKOneToOne(g, srcIds, targetIds, topK)
		pathsr := make([]PriorityQueue, len(paths))
		for i, path := range paths {
			pathsr[i] = PriorityQueue2SortedArray(path, false)
		}
		pathsb, err = json.MarshalIndent(pathsr, "", "\t")
		return
	case "TopKOneToMany":
		paths := st.TopKOneToMany(g, srcIds, targetIds, topK)
		pathsr := make([]PriorityQueue, len(paths))
		for i, path := range paths {
			pathsr[i] = PriorityQueue2SortedArray(path, false)
		}
		pathsb, err = json.MarshalIndent(pathsr, "", "\t")
		return
	default:
		return nil, fmt.Errorf("DfsDo: invalid method %s", op)
	}
}

func NewSearcher(major string, minor string) (Searcher, error) {
	switch major {
	case "dfs":
		return NewDfs(minor, 5)
	default:
		return nil, fmt.Errorf("NewSearcher: invalid option major %s, minor %s", major, minor)
	}
}
