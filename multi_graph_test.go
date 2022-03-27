package graph_shortest_paths

import (
	"reflect"
	"testing"
)

// func TestBellmanFord(t *testing.T) {
// 	var ent EntitySeq
// 	ent.FromJsonFile("./esamples/small.json")

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	cycles := graph.BellmanFord(3, 0, false)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	fmt.Println(cycles.Len())

// 	for e := cycles.Front(); e != nil; e = e.Next() {
// 		fmt.Println(e.Value)
// 	}

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

// func TestFloydWarshall(t *testing.T) {
// 	var ent EntitySeq
// 	ent.FromJsonFile("./esamples/pools.json")

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	graph.FloydWarshall(0)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

// func TestCustom(t *testing.T) {
// 	var ent EntitySeq
// 	ent.FromJsonFile("./esamples/pools.json")

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	graph.Custom(9)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

// func TestCustomTrunc(t *testing.T) {
// 	var ent EntitySeq
// 	ent.FromJsonFile("./esamples/small.json")

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	graph.Custom(1)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

// func TestFloydWarshallTrunc(t *testing.T) {
// 	var ent EntitySeq
// 	ent.FromJsonFile("./esamples/small.json")

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	graph.FloydWarshall(0)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

// func TestBfsTrunc(t *testing.T) {
// 	ent := make([]Entity, 9)
// 	ent[0] = Entity{"entityId1", 1, 2, 0.2, 0}
// 	ent[1] = Entity{"entityId2", 2, 3, 12.5, 1}
// 	ent[2] = Entity{"entityId3", 3, 1, 0.41, 2}
// 	ent[3] = Entity{"entityId4", 2, 1, 0.49, 3}
// 	ent[4] = Entity{"entityId5", 1, 4, 2050, 4}
// 	ent[5] = Entity{"entityId6", 4, 1, 0.00049, 5}
// 	ent[6] = Entity{"entityId7", 2, 4, 1.1, 6}
// 	ent[7] = Entity{"entityId8", 4, 5, 17400, 7}
// 	ent[8] = Entity{"entityId9", 5, 1, 0.0003, 8}

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	graph.Bfs(1)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

// func BenchmarkBfs(b *testing.B) {
// 	ent := LoadFromJson("./pools.json")

// 	graph := new(MultiGraph)
// 	graph.BuildGraph(ent)

// 	ts := time.Now()
// 	graph.Bfs(10)
// 	duration := time.Since(ts)
// 	fmt.Println(duration)

// 	// if amount0.Cmp(refRes) != 0 {
// 	// 	t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
// 	// }
// }

func TestMultiGraph_BuildGraph(t *testing.T) {
	type fields struct {
		Edges        MEdgeSeq
		Verteces     map[int]int
		entities     EntitySeq
		predecessors []MEdgeSeq
		successors   []MEdgeSeq
	}
	type args struct {
		ent EntitySeq
	}
	type result struct {
		vertCount     int
		edgesCount    int
		entitiesCount int
		predCount     int
		successCount  int
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes result
	}{
		{
			name:   "Init_small_graph",
			fields: fields{},
			args: args{
				ent: graphSmall.entities,
			},
			wantRes: result{
				vertCount:     5,
				edgesCount:    9,
				entitiesCount: 9,
				predCount:     5,
				successCount:  5,
			},
		},
		{
			name:   "Init_large_graph",
			fields: fields{},
			args: args{
				ent: graphLarge.entities,
			},
			wantRes: result{
				vertCount:     5393,
				edgesCount:    11934,
				entitiesCount: 13402,
				predCount:     5393,
				successCount:  5393,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &MultiGraph{
				Edges:        tt.fields.Edges,
				Verteces:     tt.fields.Verteces,
				entities:     tt.fields.entities,
				predecessors: tt.fields.predecessors,
				successors:   tt.fields.successors,
			}
			g.BuildGraph(tt.args.ent)

			gotRes := result{
				vertCount:     len(g.Verteces),
				edgesCount:    len(g.Edges),
				entitiesCount: len(g.entities),
				predCount:     len(g.predecessors),
				successCount:  len(g.successors),
			}

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Dfs.TopKShortestPaths() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
