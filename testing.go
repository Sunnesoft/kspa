package graph_shortest_paths

import "testing"

var cycles []map[uint64]MEdgeSeq
var cyclesq []PriorityQueue
var graphSmall *MultiGraph
var graphLarge *MultiGraph
var pathsGraphSmallTop5Lim5 string
var pathsGraphLargeTop100Lim5 string

func setupTestCase(t *testing.T) func(t *testing.T) {
	graphSmall = InitSmallGraph()
	graphLarge = InitLargeGraph()
	pathsGraphSmallTop5Lim5 = InitSmallGraphTop5Lim5()
	pathsGraphLargeTop100Lim5 = InitLargeGraphTop100Lim5()

	return func(t *testing.T) {
	}
}
