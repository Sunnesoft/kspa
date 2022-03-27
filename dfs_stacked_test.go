package graph_shortest_paths

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDfsStacked_TopKCycles(t *testing.T) {
	type fields struct {
		deepLimit int
		cycles    []map[uint64]MEdgeSeq
	}
	type args struct {
		topK int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []PriorityQueue
	}{
		{
			name: "test_5_10_cycles",
			fields: fields{
				deepLimit: 5,
				cycles:    cycles,
			},
			args: args{
				topK: 10,
			},
			wantRes: cyclesq,
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsStacked{
				deepLimit: tt.fields.deepLimit,
				cycles:    tt.fields.cycles,
			}
			if gotRes := st.TopKCycles(tt.args.topK); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("DfsStacked.TopKCycles() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDfsStacked_BestCycles(t *testing.T) {
	type fields struct {
		deepLimit int
		cycles    []map[uint64]MEdgeSeq
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []EdgeSeq
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsStacked{
				deepLimit: tt.fields.deepLimit,
				cycles:    tt.fields.cycles,
			}
			if gotRes := st.BestCycles(); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("DfsStacked.BestCycles() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDfsStacked_TopKShortestPaths(t *testing.T) {
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
				deepLimit: 5,
			},
			args: args{
				topK:  1000,
				srcId: 10,
				g:     graphLarge,
			},
			wantRes: pathsGraphLargeTop100Lim5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &DfsStacked{
				deepLimit: tt.fields.deepLimit,
			}

			paths := st.TopKShortestPaths(tt.args.g, tt.args.srcId, tt.args.topK)
			paths = PriorityQueue2SortedArray(paths, false)
			// jsonText, _ := json.Marshal(paths)
			// outfn := fmt.Sprintf("./examples/%s.json", tt.name)
			// WriteText(outfn, jsonText)

			pathsb, _ := json.Marshal(paths)
			pathss := string(pathsb)

			if !reflect.DeepEqual(pathss, tt.wantRes) {
				t.Errorf("Dfs.TopKShortestPaths() = %v, want %v", pathss, tt.wantRes)
			}
		})
	}
}
