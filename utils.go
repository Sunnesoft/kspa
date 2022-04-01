package kspa

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

func Weight(a *Entity) float64 {
	return -math.Log(a.Relation)
}

func IdsHash(id1, id2 int) (uint64, error) {
	if id1 > 1<<32 || id2 > 1<<32 {
		return 0, fmt.Errorf("IdsHash: ids must be non-negative number less than 2**32+1")
	}
	return uint64(id1)<<32 + uint64(id2), nil
}

func LoadText(fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ioutil.ReadAll(file)
}

func WriteText(fn string, data []byte) error {
	return os.WriteFile(fn, data, 0666)
}

func ProcessOutsideEdges(pq PriorityQueue, deepLimit int, topK int, reverseEdgeSeq bool) (res PriorityQueue) {
	mask := make([]int, deepLimit)
	limits := make([]int, deepLimit)
	path := make(EdgeSeq, deepLimit)

	res = NewPriorityQueue(0, topK)

	if pq.Len() == 0 {
		return
	}

	maxWeight := pq[0].priority

	for _, edges := range pq {
		medges := edges.value.(MEdgeSeq)

		for i := 0; i < deepLimit; i++ {
			path[i] = nil
		}

		weight := 0.0
		seqSize := len(medges)
		for i := 0; i < seqSize; i++ {
			if medges[i] == nil {
				seqSize = i
				break
			}
			limits[i] = len(medges[i].edges)
			path[i] = medges[i].edges[0]
			weight += path[i].weight
		}

		rem := 0

		for {
			for i := 1; i < seqSize && rem > 0; i++ {
				curEdges := medges[i].edges
				weight -= curEdges[mask[i]].weight
				mask[i] += rem
				mask[i], rem = mask[i]%limits[i], mask[i]/limits[i]
				path[i] = curEdges[mask[i]]
				weight += path[i].weight
			}

			if rem > 0 {
				break
			}

			if weight <= maxWeight {
				if res.Len() < topK {
					cpath := make(EdgeSeq, deepLimit)
					copy(cpath, path)

					if reverseEdgeSeq {
						cpath[:seqSize].ReverseEdgeSeq()
					}

					res.Append(cpath, weight)

					if res.Len() == topK {
						res.Init()
						maxWeight = res[0].priority
					}
				} else {
					ms, _ := res[0].value.(EdgeSeq)
					copy(ms, path)

					if reverseEdgeSeq {
						ms[:seqSize].ReverseEdgeSeq()
					}

					res.Update(res[0], res[0].value, weight)
					maxWeight = res[0].priority
				}
			}

			curEdges := medges[0].edges
			weight -= curEdges[mask[0]].weight
			mask[0] += 1
			mask[0], rem = mask[0]%limits[0], mask[0]/limits[0]
			path[0] = curEdges[mask[0]]
			weight += path[0].weight
		}
	}

	if res.Len() < topK {
		res.Init()
	}

	return
}

func firstIndexOf(vert int, path []int) int {
	for i, v := range path {
		if vert == v {
			return i
		}
	}
	return -1
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	return func(t *testing.T) {
	}
}

func reverse(numbers []int) []int {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func traceNegativeCycle(start int, predecessors []int, deepLimit int, uniquePaths bool, visited []bool) []int {
	path := make([]int, 0, deepLimit)
	path = append(path, start)

	prior := start
	for {
		prior = predecessors[prior]
		for i := 0; i < len(path); i++ {
			if prior == path[i] {
				path = path[i:]
				path = append(path, prior)
				path = reverse(path)
				return path
			}
		}

		if uniquePaths && visited[prior] {
			return nil
		}

		path = append(path, prior)
		visited[prior] = true
	}
}
