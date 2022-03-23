package graph_shortest_paths

import (
	"fmt"
	"testing"
	"time"
)

func TestBellmanFord(t *testing.T) {

	ent := make([]Entity, 9)
	ent[0] = Entity{"entityId1", 1, 2, 0.2, 0}
	ent[1] = Entity{"entityId2", 2, 3, 12.5, 1}
	ent[2] = Entity{"entityId3", 3, 1, 0.41, 2}
	ent[3] = Entity{"entityId4", 2, 1, 0.49, 3}
	ent[4] = Entity{"entityId5", 1, 4, 2050, 4}
	ent[5] = Entity{"entityId6", 4, 1, 0.00049, 5}
	ent[6] = Entity{"entityId7", 2, 4, 1.1, 6}
	ent[7] = Entity{"entityId8", 4, 5, 17400, 7}
	ent[8] = Entity{"entityId9", 5, 1, 0.0003, 8}

	// ent := loadFromJson("./pools.json")

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	cycles := graph.BellmanFord(3, 0, false)
	duration := time.Since(ts)
	fmt.Println(duration)

	fmt.Println(cycles.Len())

	for e := cycles.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}

func TestFloydWarshall(t *testing.T) {

	ent := LoadFromJson("./pools.json")

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	graph.FloydWarshall(0)
	duration := time.Since(ts)
	fmt.Println(duration)

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}

func TestCustom(t *testing.T) {

	ent := LoadFromJson("./pools.json")

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	graph.Custom(9)
	duration := time.Since(ts)
	fmt.Println(duration)

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}

func TestCustomTrunc(t *testing.T) {

	ent := make([]Entity, 9)
	ent[0] = Entity{"entityId1", 1, 2, 0.2, 0}
	ent[1] = Entity{"entityId2", 2, 3, 12.5, 1}
	ent[2] = Entity{"entityId3", 3, 1, 0.41, 2}
	ent[3] = Entity{"entityId4", 2, 1, 0.49, 3}
	ent[4] = Entity{"entityId5", 1, 4, 2050, 4}
	ent[5] = Entity{"entityId6", 4, 1, 0.00049, 5}
	ent[6] = Entity{"entityId7", 2, 4, 1.1, 6}
	ent[7] = Entity{"entityId8", 4, 5, 17400, 7}
	ent[8] = Entity{"entityId9", 5, 1, 0.0003, 8}

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	graph.Custom(1)
	duration := time.Since(ts)
	fmt.Println(duration)

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}

func TestFloydWarshallTrunc(t *testing.T) {

	ent := make([]Entity, 9)
	ent[0] = Entity{"entityId1", 1, 2, 0.2, 0}
	ent[1] = Entity{"entityId2", 2, 3, 12.5, 1}
	ent[2] = Entity{"entityId3", 3, 1, 0.41, 2}
	ent[3] = Entity{"entityId4", 2, 1, 0.49, 3}
	ent[4] = Entity{"entityId5", 1, 4, 2050, 4}
	ent[5] = Entity{"entityId6", 4, 1, 0.00049, 5}
	ent[6] = Entity{"entityId7", 2, 4, 1.1, 6}
	ent[7] = Entity{"entityId8", 4, 5, 17400, 7}
	ent[8] = Entity{"entityId9", 5, 1, 0.0003, 8}

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	graph.FloydWarshall(0)
	duration := time.Since(ts)
	fmt.Println(duration)

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}

func TestBfsTrunc(t *testing.T) {
	ent := make([]Entity, 9)
	ent[0] = Entity{"entityId1", 1, 2, 0.2, 0}
	ent[1] = Entity{"entityId2", 2, 3, 12.5, 1}
	ent[2] = Entity{"entityId3", 3, 1, 0.41, 2}
	ent[3] = Entity{"entityId4", 2, 1, 0.49, 3}
	ent[4] = Entity{"entityId5", 1, 4, 2050, 4}
	ent[5] = Entity{"entityId6", 4, 1, 0.00049, 5}
	ent[6] = Entity{"entityId7", 2, 4, 1.1, 6}
	ent[7] = Entity{"entityId8", 4, 5, 17400, 7}
	ent[8] = Entity{"entityId9", 5, 1, 0.0003, 8}

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	graph.Bfs(1)
	duration := time.Since(ts)
	fmt.Println(duration)

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}

func BenchmarkBfs(b *testing.B) {
	ent := LoadFromJson("./pools.json")

	graph := new(MultiGraph)
	graph.BuildGraph(ent)

	ts := time.Now()
	graph.Bfs(10)
	duration := time.Since(ts)
	fmt.Println(duration)

	// if amount0.Cmp(refRes) != 0 {
	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	// }
}
