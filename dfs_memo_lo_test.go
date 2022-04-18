package kspa

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

func TestOldPathsArbitrage(t *testing.T) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g      *MultiGraph
		srcIds []int
		topK   int
	}
	type testCase struct {
		name      string
		fields    fields
		args      args
		wantResFn string
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	basePath := "./examples"

	smallGraph := new(MultiGraph)
	source, _ := FromCsvFile(path.Join(basePath, "data.txt"))
	smallGraph.Build(source)

	tests := []testCase{
		{
			name: "data-1m",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:   100,
				srcIds: []int{9, 12, 15},
				g:      smallGraph,
			},
			wantResFn: path.Join(basePath, "data_5_100_9_12_15.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsMemo{}
			st.Init()
			st.SetDeepLimit(tt.fields.deepLimit)
			st.SetFnMode(FN_ALL_PATHS)
			st.SetGraph(tt.args.g)

			paths := st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsr := make([]PriorityQueue, len(paths))
			for i, path := range paths {
				pathsr[i] = PriorityQueue2SortedArray(path, false)
			}

			pathsb, err := json.MarshalIndent(pathsr, "", "\t")
			if err != nil {
				panic(err)
			}

			resPaths := make([][]ChainView, 0)
			err = json.Unmarshal(pathsb, &resPaths)
			if err != nil {
				panic(err)
			}

			// WriteText(tt.wantResFn, pathsb)

			wantResb, err := LoadText(tt.wantResFn)
			if err != nil {
				panic(err)
			}

			wantPaths := make([][]ChainView, 0)
			err = json.Unmarshal(wantResb, &wantPaths)
			if err != nil {
				panic(err)
			}

			if len(resPaths) != len(wantPaths) {
				t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
				return
			}

			for i := 0; i < len(resPaths); i++ {
				freq := make(map[string]int)

				if len(resPaths[i]) != len(wantPaths[i]) {
					t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
					return
				}

				for j := 0; j < len(resPaths[i]); j++ {
					freq[resPaths[i][j].Chain]++
					freq[wantPaths[i][j].Chain]++
				}

				for _, v := range freq {
					if v != 2 {
						t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
						return
					}
				}
			}
		})
	}
}

func TestLimitOrdersArbitrage(t *testing.T) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g      *MultiGraph
		srcIds []int
		topK   int
		lo     []*SingleEdge
	}
	type testCase struct {
		name      string
		fields    fields
		args      args
		wantResFn string
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	basePath := "./examples"
	loFp := path.Join(basePath, "data.txt")

	graph := new(MultiGraph)
	source, _ := FromCsvFile(loFp)
	graph.Build(source)

	lo, _ := FromCsvFile(path.Join(basePath, "limit_orders.txt"))

	tests := []testCase{
		{
			name: "data-1m",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:   100,
				srcIds: []int{9, 12, 15},
				g:      graph,
				lo:     lo,
			},
			wantResFn: path.Join(basePath, "lo_5_100_9_12_15.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsMemo{}
			st.Init()
			st.SetDeepLimit(tt.fields.deepLimit)
			st.SetFnMode(FN_LO_ONLY)
			st.SetGraph(tt.args.g)
			st.AddLimitOrders(tt.args.lo)

			paths := st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsr := make([]PriorityQueue, len(paths))
			for i, path := range paths {
				pathsr[i] = PriorityQueue2SortedArray(path, false)
			}

			pathsb, err := json.MarshalIndent(pathsr, "", "\t")
			if err != nil {
				panic(err)
			}

			resPaths := make([][]ChainView, 0)
			err = json.Unmarshal(pathsb, &resPaths)
			if err != nil {
				panic(err)
			}

			// WriteText(tt.wantResFn, pathsb)

			wantResb, err := LoadText(tt.wantResFn)
			if err != nil {
				panic(err)
			}

			wantPaths := make([][]ChainView, 0)
			err = json.Unmarshal(wantResb, &wantPaths)
			if err != nil {
				panic(err)
			}

			if len(resPaths) != len(wantPaths) {
				t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
				return
			}

			for i := 0; i < len(resPaths); i++ {
				freq := make(map[string]int)

				if len(resPaths[i]) != len(wantPaths[i]) {
					t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
					return
				}

				for j := 0; j < len(resPaths[i]); j++ {
					freq[resPaths[i][j].Chain]++
					freq[wantPaths[i][j].Chain]++
				}

				for _, v := range freq {
					if v != 2 {
						t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
						return
					}
				}
			}
		})
	}
}

func TestRandomLimitOrdersArbitrage(t *testing.T) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g      *MultiGraph
		srcIds []int
		topK   int
		lo     []*SingleEdge
	}
	type testCase struct {
		name      string
		fields    fields
		args      args
		wantResFn string
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	basePath := "./examples"
	loFp := path.Join(basePath, "data.txt")

	graph := new(MultiGraph)
	source, _ := FromCsvFile(loFp)
	graph.Build(source)

	lo := GenerateRandomLimitOrders(loFp, 300, 5.0)
	err := ToCsvFile(path.Join(basePath, "limit_orders_random.txt"), lo)
	if err != nil {
		panic(err)
	}

	tests := []testCase{
		{
			name: "data-1m",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:   100,
				srcIds: []int{9, 12, 15},
				g:      graph,
				lo:     lo,
			},
			wantResFn: path.Join(basePath, "lo_rand_5_100_9_12_15.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsMemo{}
			st.Init()
			st.SetDeepLimit(tt.fields.deepLimit)
			st.SetFnMode(FN_LO_ONLY)
			st.SetGraph(tt.args.g)
			st.AddLimitOrders(tt.args.lo)

			paths := st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsr := make([]PriorityQueue, len(paths))
			for i, path := range paths {
				pathsr[i] = PriorityQueue2SortedArray(path, false)
			}

			pathsb, err := json.MarshalIndent(pathsr, "", "\t")
			if err != nil {
				panic(err)
			}

			WriteText(tt.wantResFn, pathsb)
		})
	}
}

func BenchmarkRandomLimitOrdersArbitrage(b *testing.B) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g      string
		srcIds []int
		topK   int
		lo     string
	}
	type testCase struct {
		name   string
		fields fields
		args   args
	}

	rand.Seed(time.Now().UnixNano())

	basePath := "./examples"
	dataPath := path.Join(basePath, "data.txt")

	sourceConfig := struct {
		path      string
		data      string
		count     int
		removeOld bool
		c         RandomEdgeSeqInfo
	}{
		path:      "./benchmark/lo300",
		data:      dataPath,
		count:     5,
		removeOld: true,
		c: RandomEdgeSeqInfo{
			Count:    300,
			PercDiff: 5,
		},
	}

	GenerateRandomLimitOrdersCsv(sourceConfig.path, sourceConfig.data, sourceConfig.count, sourceConfig.removeOld, sourceConfig.c)

	files, _ := ioutil.ReadDir(sourceConfig.path)

	bchmConfig := make([]testCase, 0)

	for _, fn := range files {
		if fn.IsDir() {
			continue
		}

		bchmConfig = append(bchmConfig, testCase{
			name: fn.Name(),
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:   100,
				srcIds: []int{9, 12, 15},
				g:      dataPath,
				lo:     path.Join(sourceConfig.path, fn.Name()),
			},
		})
	}

	b.ResetTimer()

	for _, bb := range bchmConfig {

		b.Run("Arbitrage", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				graph := new(MultiGraph)
				source, _ := FromCsvFile(bb.args.g)
				graph.Build(source)

				lo, _ := FromCsvFile(bb.args.lo)
				b.StartTimer()

				st := &DfsMemo{}
				st.Init()
				st.SetDeepLimit(bb.fields.deepLimit)
				st.SetFnMode(FN_LO_ONLY)
				st.SetGraph(graph)
				st.AddLimitOrders(lo)

				paths := st.Arbitrage(bb.args.srcIds, bb.args.topK)
				_ = paths
			}
		})
	}

	os.RemoveAll(sourceConfig.path)
}
