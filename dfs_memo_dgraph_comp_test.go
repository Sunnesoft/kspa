package kspa

import (
	"dgraph"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

func BenchmarkMemoDgraphCmp(b *testing.B) {
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
		count:     100,
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
				deepLimit: 4,
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

		b.Run("Arbitrage-KSPA", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				source, _ := FromCsvFile(bb.args.g)
				lo, _ := FromCsvFile(bb.args.lo)
				b.StartTimer()

				graph := new(MultiGraph)
				graph.Build(source)

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

		b.Run("Arbitrage-Dgraph", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				source, _ := FromCsvFile(bb.args.g)
				lo, _ := FromCsvFile(bb.args.lo)
				b.StartTimer()

				g := dgraph.NewGraph()
				for _, sedge := range source {
					ent := sedge.Data
					v1 := g.AddVertex(ent.Id1)
					v2 := g.AddVertex(ent.Id2)
					id, err := strconv.Atoi(ent.EntityId)
					if err != nil {
						continue
					}
					g.SetEdge(v1, v2, ent.Relation, id, "")

				}
				// g.MakeMatrix(true)

				for _, sedge := range lo {
					ent := sedge.Data
					v1 := g.Vertices[ent.Id1]
					v2 := g.Vertices[ent.Id1]
					g.SetEdge(v1, v2, ent.Relation, 0, ent.EntityId)
				}
				g.MakeMatrix(true)

				for _, src := range bb.args.srcIds {
					pathsd, total := g.AllPathsWithTemp([]int{src}, []int{src}, bb.fields.deepLimit, bb.args.topK)
					_ = pathsd
					_ = total
				}

				// for i, p := range pathsd {
				// 	fmt.Printf("%d. (%d) ", i+1, p.Chain[0].Source.Id)
				// 	for _, e := range p.Chain {
				// 		fmt.Printf("(%d) ", e.Target.Id)
				// 	}
				// 	fmt.Printf("  rate: %.4f\n", p.Rate)
				// }
			}
		})

		b.Run("Arbitrage-KSPA-add-single-lo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				source, _ := FromCsvFile(bb.args.g)
				lo, _ := FromCsvFile(bb.args.lo)
				b.StartTimer()

				graph := new(MultiGraph)
				graph.Build(source)

				st := &DfsMemo{}
				st.Init()
				st.SetDeepLimit(bb.fields.deepLimit)
				st.SetFnMode(FN_LO_ONLY)
				st.SetGraph(graph)

				for _, singleLo := range lo {
					_, err := st.AddLimitOrders(EdgeSeq{singleLo})
					if err != nil {
						panic(err)
					}

					paths := st.Arbitrage(bb.args.srcIds, bb.args.topK)
					_ = paths
				}

			}
		})
	}

	os.RemoveAll(sourceConfig.path)
}
