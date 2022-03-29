package graph_shortest_paths

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

func TestDfsMemoOp(t *testing.T) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g         *MultiGraph
		srcIds    []int
		targetIds []int
		topK      int
	}
	type testCase struct {
		name      string
		fields    fields
		args      args
		wantResFn string
		function  string
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	basePath := "./examples"

	smallGraph := new(MultiGraph)
	smallGraph.BuildGraph(FromJsonFile(path.Join(basePath, "small.json")))

	largeGraph := new(MultiGraph)
	largeGraph.BuildGraph(FromJsonFile(path.Join(basePath, "pools.json")))

	tests := []testCase{
		{
			name: "small-1",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:      5,
				srcIds:    []int{1},
				targetIds: []int{1},
				g:         smallGraph,
			},
			function:  "TopK",
			wantResFn: path.Join(basePath, "small_5_5_1o.json"),
		},
		{
			name: "small-2",
			fields: fields{
				deepLimit: 4,
			},
			args: args{
				topK:      3,
				srcIds:    []int{1},
				targetIds: []int{1},
				g:         smallGraph,
			},
			function:  "TopK",
			wantResFn: path.Join(basePath, "small_3_4_1o.json"),
		},
		{
			name: "large-1",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:      100,
				srcIds:    []int{10},
				targetIds: []int{10},
				g:         largeGraph,
			},
			function:  "TopK",
			wantResFn: path.Join(basePath, "pools_100_5_10o.json"),
		},
		{
			name: "large-2",
			fields: fields{
				deepLimit: 6,
			},
			args: args{
				topK:      100,
				srcIds:    []int{10},
				targetIds: []int{10},
				g:         largeGraph,
			},
			function:  "TopK",
			wantResFn: path.Join(basePath, "pools_100_6_10o.json"),
		},
		{
			name: "large-3",
			fields: fields{
				deepLimit: 7,
			},
			args: args{
				topK:      100,
				srcIds:    []int{10},
				targetIds: []int{10},
				g:         largeGraph,
			},
			function:  "TopK",
			wantResFn: path.Join(basePath, "pools_100_7_10o.json"),
		},
		{
			name: "large-4",
			fields: fields{
				deepLimit: 8,
			},
			args: args{
				topK:      10,
				srcIds:    []int{10},
				targetIds: []int{10},
				g:         largeGraph,
			},
			function:  "TopK",
			wantResFn: path.Join(basePath, "pools_10_8_10o.json"),
		},
		{
			name: "large-5",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:      10,
				srcIds:    []int{10, 9, 15, 22, 3966, 450, 516, 2900, 70, 91},
				targetIds: []int{10, 9, 15, 22, 3966, 450, 516, 2900, 70, 91},
				g:         largeGraph,
			},
			function:  "TopKOneToOne",
			wantResFn: path.Join(basePath, "pools_10_5_10m.json"),
		},
		{
			name: "large-6",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:      50,
				srcIds:    []int{10, 22, 15, 9, 450},
				targetIds: []int{10, 22, 15, 9, 450},
				g:         largeGraph,
			},
			function:  "TopKOneToMany",
			wantResFn: path.Join(basePath, "pools_50_5_5v.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsMemo{}
			st.Init()
			st.SetDeepLimit(tt.fields.deepLimit)

			var pathsb []byte

			switch tt.function {
			case "TopK":
				paths := st.TopK(tt.args.g, tt.args.srcIds[0], tt.args.targetIds[0], tt.args.topK)
				pathsr := PriorityQueue2SortedArray(paths, false)
				pathsb, _ = json.MarshalIndent(pathsr, "", "\t")
			case "TopKOneToOne":
				paths := st.TopKOneToOne(tt.args.g, tt.args.srcIds, tt.args.targetIds, tt.args.topK)
				pathsr := make([]PriorityQueue, len(paths))
				for i, path := range paths {
					pathsr[i] = PriorityQueue2SortedArray(path, false)
				}
				pathsb, _ = json.MarshalIndent(pathsr, "", "\t")
			case "TopKOneToMany":
				paths := st.TopKOneToMany(tt.args.g, tt.args.srcIds, tt.args.targetIds, tt.args.topK)
				pathsr := make([]PriorityQueue, len(paths))
				for i, path := range paths {
					pathsr[i] = PriorityQueue2SortedArray(path, false)
				}
				pathsb, _ = json.MarshalIndent(pathsr, "", "\t")
			}

			// WriteText(tt.wantResFn, pathsb)

			wantResb, _ := LoadText(tt.wantResFn)
			wantRes := string(wantResb)

			pathss := string(pathsb)
			if !reflect.DeepEqual(pathss, wantRes) {
				t.Errorf("DfsMemo.%s() = %v, want %v", tt.function, pathss, wantRes)
			}
		})
	}
}

func BenchmarkDfsMemoOp(b *testing.B) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g    *MultiGraph
		topK int
	}
	type testCase struct {
		name         string
		fields       fields
		args         args
		oneCount     int
		oneToOneMax  int
		oneToManyMax int
	}

	rand.Seed(time.Now().UnixNano())

	sourceConfig := struct {
		path      string
		count     int
		removeOld bool
		c         RandomEntitySeqInfo
	}{
		path:      "./benchmark/v5000_e20000",
		count:     5,
		removeOld: true,
		c: RandomEntitySeqInfo{
			VertexCount:     5000,
			VertexStdFactor: 50,
			EdgesCount:      40000,
			RelationMin:     0.0,
			RelationMax:     100000.0,
			NoiseMean:       0.0,
			NoiseStdDev:     0.0001,
		},
	}

	GenerateRandomEntitiesJson(sourceConfig.path, sourceConfig.count, sourceConfig.removeOld, sourceConfig.c)

	files, _ := ioutil.ReadDir(sourceConfig.path)

	bchmConfig := make([]testCase, 0)

	for _, fn := range files {
		if fn.IsDir() {
			continue
		}

		entities := FromJsonFile(path.Join(sourceConfig.path, fn.Name()))
		graph := new(MultiGraph)
		graph.BuildGraph(entities)

		bchmConfig = append(bchmConfig, testCase{
			name: fn.Name(),
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK: 100,
				g:    graph,
			},
			oneCount:     30,
			oneToOneMax:  15,
			oneToManyMax: 5,
		})
	}

	b.ResetTimer()

	for _, bb := range bchmConfig {

		for j := 0; j < bb.oneCount; j++ {
			v := 0

			for vert := range bb.args.g.Verteces {
				if rand.Intn(2) > 0 {
					v = vert
					break
				}
			}

			b.Run("TopK", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					st := &DfsMemo{}
					st.Init()
					st.SetDeepLimit(bb.fields.deepLimit)

					paths := st.TopK(bb.args.g, v, v, bb.args.topK)
					_ = paths
				}
			})
		}

		for j := 1; j <= bb.oneToOneMax; j++ {
			ids := make([]int, 0, j)

			for vert := range bb.args.g.Verteces {
				if len(ids) == j {
					break
				}

				if rand.Intn(2) > 0 {
					ids = append(ids, vert)
				}
			}

			b.Run(fmt.Sprintf("TopKOneToOne_%d", j), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					st := &DfsMemo{}
					st.Init()
					st.SetDeepLimit(bb.fields.deepLimit)

					paths := st.TopKOneToOne(bb.args.g, ids, ids, bb.args.topK)
					_ = paths
				}
			})
		}

		for j := 1; j <= bb.oneToManyMax; j++ {
			ids := make([]int, 0, j)

			for vert := range bb.args.g.Verteces {
				if len(ids) == j {
					break
				}

				if rand.Intn(2) > 0 {
					ids = append(ids, vert)
				}
			}

			b.Run(fmt.Sprintf("TopKOneToMany_%d", j), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					st := &DfsMemo{}
					st.Init()
					st.SetDeepLimit(bb.fields.deepLimit)

					paths := st.TopKOneToMany(bb.args.g, ids, ids, bb.args.topK)
					_ = paths
				}
			})
		}
	}

	os.RemoveAll(sourceConfig.path)
}
