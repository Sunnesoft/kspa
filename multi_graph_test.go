package kspa

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

func TestMultiGraph_BuildGraph(t *testing.T) {
	type fields struct {
		Edges        MEdgeSeq
		VertexIndex  map[int]int
		entities     EdgeSeq
		predecessors []MEdgeSeq
		successors   []MEdgeSeq
	}
	type args struct {
		ent EdgeSeq
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

	smallGraphEntities, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "small.json"))
	largeGraphEntities, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "pools.json"))

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
				VertexIndex:  tt.fields.VertexIndex,
				entities:     tt.fields.entities,
				predecessors: tt.fields.predecessors,
				successors:   tt.fields.successors,
			}
			g.Build(tt.args.ent)

			gotRes := result{
				vertCount:     len(g.VertexIndex),
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

func TestMultiGraph_UpdateRelation(t *testing.T) {
	type fields struct {
		entities EdgeSeq
	}
	type args struct {
		ent EdgeSeq
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	basePath := "./examples"

	smallGraphEntities, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "small.json"))
	smallGraphUpdate, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "small_update.json"))

	largeGraphEntities, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "pools.json"))
	largeGraphUpdate, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "pools_update.json"))
	largeGraphUpdateFail, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "pools_update_fail.json"))
	largeGraphUpdateInvIds, _ := FromJsonFile[*SingleEdge](path.Join(basePath, "pools_update_invalid_ids.json"))

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		indexCheck int
	}{
		{
			name: "small-update",
			fields: fields{
				entities: smallGraphEntities,
			},
			args: args{
				ent: smallGraphUpdate,
			},
			wantErr:    false,
			indexCheck: 0,
		},
		{
			name: "large-update",
			fields: fields{
				entities: largeGraphEntities,
			},
			args: args{
				ent: largeGraphUpdate,
			},
			wantErr:    false,
			indexCheck: 4,
		},
		{
			name: "large-update-fail",
			fields: fields{
				entities: largeGraphEntities,
			},
			args: args{
				ent: largeGraphUpdateFail,
			},
			wantErr: true,
		},
		{
			name: "large-update-invalid-ids",
			fields: fields{
				entities: largeGraphEntities,
			},
			args: args{
				ent: largeGraphUpdateInvIds,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &MultiGraph{}
			g.Build(tt.fields.entities)

			te := tt.fields.entities[tt.indexCheck]
			relationBefore := 0.0
			index, getIndexOk := g.GetEdgeIndex(te.Data.Id1, te.Data.Id2)

			if !getIndexOk {
				t.Errorf("MultiGraph.GetEdgeIndex() return %d, %t", index, getIndexOk)
			}

			if getIndexOk {
				medge := g.Edges[index]
				relationBefore = medge.edges[medge.index[te.Data.EntityId]].Data.Relation
			}

			err := g.UpdateRelation(tt.args.ent)
			relationAfter := 0.0

			if getIndexOk {
				medge := g.Edges[index]
				relationAfter = medge.edges[medge.index[te.Data.EntityId]].Data.Relation
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("MultiGraph.UpdateRelation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if relationAfter == relationBefore {
					t.Errorf("MultiGraph.UpdateRelation() relation not updated after %f, before %f", relationAfter, relationBefore)
				}
			}
		})
	}
}

func BenchmarkMultiGraph_UpdateRelation(b *testing.B) {
	type fields struct {
		entities EdgeSeq
	}
	type args struct {
		ent EdgeSeq
	}
	type testCase struct {
		name   string
		fields fields
		args   args
	}

	rand.Seed(time.Now().UnixNano())

	sourceConfig := struct {
		path      string
		count     int
		removeOld bool
		c         RandomEntitySeqInfo
	}{
		path:      "./benchmark/v5000_e20000",
		count:     10,
		removeOld: true,
		c: RandomEntitySeqInfo{
			VertexCount:     5000,
			VertexStdFactor: 50,
			EdgesCount:      20000,
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

		entities, _ := FromJsonFile[*SingleEdge](path.Join(sourceConfig.path, fn.Name()))
		bchmConfig = append(bchmConfig, testCase{
			name: fn.Name(),
			fields: fields{
				entities: entities,
			},
			args: args{
				ent: entities,
			},
		})
	}

	b.ResetTimer()

	for _, bb := range bchmConfig {
		b.Run(bb.name, func(b *testing.B) {
			g := &MultiGraph{}
			g.Build(bb.fields.entities)

			err := g.UpdateRelation(bb.args.ent)

			if err != nil {
				b.Errorf("MultiGraph.UpdateRelation() error = %v, wantErr %v", err, nil)
			}
		})
	}

	os.RemoveAll(sourceConfig.path)
}
