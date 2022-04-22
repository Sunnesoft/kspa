package kspa

import (
	"dgraph"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func DgraphPathToChainView(seq *dgraph.Path) (*ChainView, bool) {
	s := make([]string, 0, len(seq.Chain))
	value := seq.Rate
	var in int64 = -1
	var out int64 = -1

	if len(seq.Chain) > 0 {
		in = int64(seq.Chain[0].Source.Id)
	}

	existTempId := false

	for _, edge := range seq.Chain {
		if edge == nil {
			break
		}

		if edge.Id == 0 && edge.TempId != "" {
			existTempId = true
			s = append(s, edge.TempId)
		} else {
			s = append(s, fmt.Sprint(edge.Id))
		}
		out = int64(edge.Target.Id)
	}

	chain := strings.Join(s, " -> ")

	return &ChainView{
		In:    in,
		Out:   out,
		Chain: chain,
		Value: value,
	}, existTempId
}

func BenchmarkMemoDgraphCmp(b *testing.B) {
	type fields struct {
		deepLimit int
	}
	type args struct {
		g      string
		srcIds []int
		topK   int
		lo     string
		baselo string
	}
	type testCase struct {
		name     string
		fields   fields
		args     args
		outputFn string
	}

	rand.Seed(time.Now().UnixNano())

	basePath := "./examples"
	dataPath := path.Join(basePath, "data.txt")
	outPath := "./benchmark/out"

	os.RemoveAll(outPath)
	err := os.MkdirAll(outPath, 0766)

	if err != nil {
		panic(err)
	}

	sourceConfig := struct {
		path      string
		data      string
		count     int
		removeOld bool
		c         RandomEdgeSeqInfo
	}{
		path:      "./benchmark/lo300",
		data:      dataPath,
		count:     101,
		removeOld: true,
		c: RandomEdgeSeqInfo{
			Count:    300,
			PercDiff: 5,
		},
	}

	GenerateRandomLimitOrdersCsv(sourceConfig.path, sourceConfig.data, sourceConfig.count, sourceConfig.removeOld, sourceConfig.c)

	files, _ := ioutil.ReadDir(sourceConfig.path)

	bchmConfig := make([]testCase, 0)

	baseLoPath := ""

	for _, fn := range files {
		if fn.IsDir() {
			continue
		}

		if baseLoPath == "" {
			baseLoPath = path.Join(sourceConfig.path, fn.Name())
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
				baselo: baseLoPath,
			},
			outputFn: fn.Name(),
		})
	}

	b.ResetTimer()

	for _, bb := range bchmConfig {

		b.Run("ArbitrageKSPA", func(b *testing.B) {
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
				medges, err := st.AddLimitOrders(lo)
				_ = medges
				if err != nil {
					panic(err)
				}

				paths := st.Arbitrage(bb.args.srcIds, bb.args.topK)

				b.StopTimer()
				pathsb, err := PathsToJson(paths)
				if err != nil {
					panic(err)
				}
				err = WriteText(path.Join(outPath, fmt.Sprintf("kspa_%s", bb.outputFn)), pathsb)
				if err != nil {
					panic(err)
				}
				b.StartTimer()
			}
		})

		b.Run("ArbitrageDgraph", func(b *testing.B) {
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
					v2 := g.Vertices[ent.Id2]
					g.SetEdge(v1, v2, ent.Relation, 0, ent.EntityId)
				}
				g.MakeMatrix(true)

				cv := make([][]*ChainView, len(bb.args.srcIds))

				for i, src := range bb.args.srcIds {
					pathsd, _ := g.AllPathsWithTemp([]int{src}, []int{src}, bb.fields.deepLimit, bb.args.topK)

					b.StopTimer()
					cv[i] = make([]*ChainView, 0)
					for _, val := range pathsd {
						if view, hasLo := DgraphPathToChainView(val); hasLo {
							cv[i] = append(cv[i], view)
						}
					}
					b.StartTimer()
				}

				b.StopTimer()
				jsonText, err := json.MarshalIndent(cv, "", "\t")
				if err != nil {
					panic(err)
				}
				err = WriteText(path.Join(outPath, fmt.Sprintf("dgraph_%s", bb.outputFn)), jsonText)
				if err != nil {
					panic(err)
				}
				b.StartTimer()
			}
		})

		b.Run("KspaAdd300Lo1LoAtTime", func(b *testing.B) {
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

		b.Run("KspaAdd300LoAsBatchAdd1LoRem1LoAtTime", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				source, _ := FromCsvFile(bb.args.g)
				lo, _ := FromCsvFile(bb.args.lo)
				baselo, _ := FromCsvFile(bb.args.baselo)
				b.StartTimer()

				graph := new(MultiGraph)
				graph.Build(source)

				st := &DfsMemo{}
				st.Init()
				st.SetDeepLimit(bb.fields.deepLimit)
				st.SetFnMode(FN_LO_ONLY)
				st.SetGraph(graph)

				medges := make(MEdgeSeq, 0)

				// kedges, _ := st.AddLimitOrders(append(baselo, lo...))
				// _ = kedges

				for _, singleLo := range baselo[:150] {
					ledges, err := st.AddLimitOrders(EdgeSeq{singleLo})
					if err != nil {
						panic(err)
					}

					medges = append(medges, ledges...)
				}

				ledges, err := st.AddLimitOrders(lo[:150])
				if err != nil {
					panic(err)
				}
				medges = append(medges, ledges...)

				for _, singleLo := range lo[150:] {
					ledges, err := st.AddLimitOrders(EdgeSeq{singleLo})
					if err != nil {
						panic(err)
					}

					medges = append(medges, ledges...)

					loRemovals := medges[0]
					medges = medges[1:]
					st.RemoveLimitOrders(MEdgeSeq{loRemovals})

					paths := st.Arbitrage(bb.args.srcIds, bb.args.topK)
					_ = paths
				}

			}
		})
	}

	os.RemoveAll(sourceConfig.path)
}
