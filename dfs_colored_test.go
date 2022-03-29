package graph_shortest_paths

import (
	"encoding/json"
	"path"
	"reflect"
	"testing"
)

func TestDfsColoredOp(t *testing.T) {
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
			wantResFn: path.Join(basePath, "pools_100_5_10o_col.json"),
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
			wantResFn: path.Join(basePath, "pools_100_6_10o_col.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsColored{}
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
				t.Errorf("DfsColored.%s() = %v, want %v", tt.function, pathss, wantRes)
			}
		})
	}
}
