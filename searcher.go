package graph_shortest_paths

type Searcher interface {
	TopKShortestPaths(*MultiGraph, int, int) PriorityQueue
}
