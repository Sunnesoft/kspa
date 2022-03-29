package graph_shortest_paths

import (
	"path"
	"reflect"
	"testing"
)

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

	basePath := "./examples"

	smallGraphEntities := FromJsonFile(path.Join(basePath, "small.json"))
	largeGraphEntities := FromJsonFile(path.Join(basePath, "pools.json"))

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
				ent: smallGraphEntities,
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
				ent: largeGraphEntities,
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
