package graph_shortest_paths

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDfsMemo_TopKShortestPaths(t *testing.T) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g     *MultiGraph
		srcId int
		topK  int
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
	}{
		{
			name: "small_top5_lim5",
			fields: fields{
				deepLimit: 5,
			},
			args: args{
				topK:  5,
				srcId: 1,
				g:     graphSmall,
			},
			wantRes: pathsGraphSmallTop5Lim5,
		},
		{
			name: "large_top100_lim5",
			fields: fields{
				deepLimit: 7,
			},
			args: args{
				topK:  100,
				srcId: 10,
				g:     graphLarge,
			},
			wantRes: pathsGraphLargeTop100Lim5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsMemo{
				deepLimit: tt.fields.deepLimit,
			}

			paths := st.TopKShortestPaths(tt.args.g, tt.args.srcId, tt.args.topK)
			paths = PriorityQueue2SortedArray(paths, false)
			jsonText, _ := json.Marshal(paths)
			WriteText("./examples/small111.json", jsonText)

			pathsb, _ := json.Marshal(paths)
			pathss := string(pathsb)

			if !reflect.DeepEqual(pathss, tt.wantRes) {
				t.Errorf("Dfs.TopKShortestPaths() = %v, want %v", pathss, tt.wantRes)
			}
		})
	}
}
