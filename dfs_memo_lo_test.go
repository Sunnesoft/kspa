package kspa

import (
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
			name: "old-fn-arb",
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
			pathsb, err := PathsToJson(paths)
			if err != nil {
				panic(err)
			}

			resPaths, err := PathsToChainView(pathsb)
			if err != nil {
				panic(err)
			}

			// WriteText(tt.wantResFn, pathsb)

			wantResb, err := LoadText(tt.wantResFn)
			if err != nil {
				panic(err)
			}

			wantPaths, err := PathsToChainView(wantResb)
			if err != nil {
				panic(err)
			}

			if !IsChainViewsEquals(resPaths, wantPaths, 3) {
				t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
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
			name: "lo-arb",
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
			_, err := st.AddLimitOrders(tt.args.lo)
			if err != nil {
				panic(err)
			}

			paths := st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsb, err := PathsToJson(paths)
			if err != nil {
				panic(err)
			}

			resPaths, err := PathsToChainView(pathsb)
			if err != nil {
				panic(err)
			}

			// WriteText(tt.wantResFn, pathsb)

			wantResb, err := LoadText(tt.wantResFn)
			if err != nil {
				panic(err)
			}

			wantPaths, err := PathsToChainView(wantResb)
			if err != nil {
				panic(err)
			}

			if !IsChainViewsEquals(resPaths, wantPaths, 3) {
				t.Errorf("DfsMemo.Arbitrage() = %v, want %v", string(pathsb), string(wantResb))
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
			name: "lo-rand-arb",
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
			_, err := st.AddLimitOrders(tt.args.lo)
			if err != nil {
				panic(err)
			}

			paths := st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsb, err := PathsToJson(paths)
			if err != nil {
				panic(err)
			}

			err = WriteText(tt.wantResFn, pathsb)
			if err != nil {
				panic(err)
			}
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
				_, err := st.AddLimitOrders(lo)
				if err != nil {
					panic(err)
				}

				paths := st.Arbitrage(bb.args.srcIds, bb.args.topK)
				_ = paths
			}
		})
	}

	os.RemoveAll(sourceConfig.path)
}

func TestLimitOrdersUpdating(t *testing.T) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		dataPath       string
		srcIds         []int
		topK           int
		loPath         string
		loUpdationPath string
	}
	type testCase struct {
		name                     string
		fields                   fields
		args                     args
		loFilePathOutput         string
		loUpdationFilePathOutput string
		loRemovalsFilePathOutput string
		loRecoverFilePathOutput  string
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	basePath := "./examples"
	dataPath := path.Join(basePath, "data.txt")

	tests := []testCase{
		{
			name: "lo-add-remove",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:           100,
				srcIds:         []int{9, 12, 15},
				dataPath:       dataPath,
				loPath:         path.Join(basePath, "limit_orders.txt"),
				loUpdationPath: path.Join(basePath, "limit_orders_update.txt"),
			},
			loFilePathOutput:         path.Join(basePath, "lo_5_100_9_12_15_init.json"),
			loUpdationFilePathOutput: path.Join(basePath, "lo_5_100_9_12_15_upd.json"),
			loRemovalsFilePathOutput: path.Join(basePath, "lo_5_100_9_12_15_rem.json"),
			loRecoverFilePathOutput:  path.Join(basePath, "lo_5_100_9_12_15_rec.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := new(MultiGraph)
			source, _ := FromCsvFile(tt.args.dataPath)
			graph.Build(source)

			lo, _ := FromCsvFile(tt.args.loPath)
			loUpdation, _ := FromCsvFile(tt.args.loUpdationPath)

			st := &DfsMemo{}
			st.Init()
			st.SetDeepLimit(tt.fields.deepLimit)
			st.SetFnMode(FN_LO_ONLY)
			st.SetGraph(graph)

			// Adding limit orders as virtual edges into the Graph
			medges, err := st.AddLimitOrders(lo)
			if err != nil {
				panic(err)
			}

			paths := st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsb, err := PathsToJson(paths)
			if err != nil {
				panic(err)
			}
			err = WriteText(tt.loFilePathOutput, pathsb)
			if err != nil {
				panic(err)
			}

			wantChains, err := PathsToChainView(pathsb)
			if err != nil {
				panic(err)
			}

			// Adding another limit orders as virtual edges into the Graph
			medgesUpdation, err := st.AddLimitOrders(loUpdation)
			if err != nil {
				panic(err)
			}

			medges = append(medges, medgesUpdation...)

			paths = st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsb, err = PathsToJson(paths)
			if err != nil {
				panic(err)
			}
			err = WriteText(tt.loUpdationFilePathOutput, pathsb)
			if err != nil {
				panic(err)
			}

			// Remove last limit order from the Graph
			loRemovals, medges := medges[3:], medges[0:3]
			st.RemoveLimitOrders(loRemovals)

			paths = st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsb, err = PathsToJson(paths)
			if err != nil {
				panic(err)
			}
			err = WriteText(tt.loRemovalsFilePathOutput, pathsb)
			if err != nil {
				panic(err)
			}

			// Remove last limit order from the Graph
			loRemovals, medges = medges[2:], medges[0:2]
			st.RemoveLimitOrders(loRemovals)

			paths = st.Arbitrage(tt.args.srcIds, tt.args.topK)
			pathsb, err = PathsToJson(paths)
			if err != nil {
				panic(err)
			}
			err = WriteText(tt.loRecoverFilePathOutput, pathsb)
			if err != nil {
				panic(err)
			}

			resChains, err := PathsToChainView(pathsb)
			if err != nil {
				panic(err)
			}

			// Result paths must be the same as paths with initial limit orders
			if !IsChainViewsEquals(resChains, wantChains, 3) {
				t.Errorf("DfsMemo.Arbitrage() watch in %v, want in %v", tt.loRecoverFilePathOutput, tt.loFilePathOutput)
			}
		})
	}
}
