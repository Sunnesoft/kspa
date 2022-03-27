package graph_shortest_paths

func InitSmallGraph() *MultiGraph {
	ent := FromJsonFile("./examples/small.json")

	g := new(MultiGraph)
	g.BuildGraph(ent)

	return g
}

func InitSmallGraphTop5Lim5() string {
	data, _ := LoadText("./examples/small_top5_lim5.json")
	return string(data)
}

func InitLargeGraphTop100Lim5() string {
	data, _ := LoadText("./examples/large_top100_lim5.json")
	return string(data)
}

func InitLargeGraph() *MultiGraph {
	ent := FromJsonFile("./examples/pools.json")

	g := new(MultiGraph)
	g.BuildGraph(ent)

	return g
}

func InitLargeGraphResult() string {
	data, _ := LoadText("./examples/pools_res.json")
	return string(data)
}
