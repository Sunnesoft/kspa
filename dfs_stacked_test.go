package kspa

import (
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

	cycles := make([]map[uint64]MEdgeSeq, 0)
	cyclesq := make([]PriorityQueue, 0)

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
