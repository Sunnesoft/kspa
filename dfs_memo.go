package graph_shortest_paths

import "math"

var MIN_WEIGHT float64 = math.Inf(-1)
var MAX_WEIGHT float64 = math.Inf(1)

const (
	THR_ZERO = iota
	THR_MEAN
	THR_MEAN_STDDEV
	THR_CUSTOM
)

type TreeNode struct {
	base   *MultiEdge
	src    int
	target int
	level  int
}

func NewTreeNode(base *MultiEdge, src int, target int, level int) *TreeNode {
	return &TreeNode{
		base:   base,
		src:    src,
		target: target,
		level:  level,
	}
}

type TreeItemStat struct {
	minWeight   float64
	maxWeight   float64
	meanWeight  float64
	mean2Weight float64
	pathsCount  float64
}

type TreeItem struct {
	stat TreeItemStat
	seq  []*TreeNode
}

type DfsMemo struct {
	deepLimit int
	topK      int

	g *MultiGraph

	psa   []float64
	edges MEdgeSeq

	memo map[int]map[int]map[int]*TreeItem

	pq        PriorityQueue
	maxWeight float64

	threshold     float64
	thresholdMode int
}

func (st *DfsMemo) TopKShortestPaths(g *MultiGraph, srcId int, topK int) (res PriorityQueue) {
	st.psa = make([]float64, st.deepLimit+2)
	st.edges = make(MEdgeSeq, st.deepLimit+1)
	st.g = g
	st.topK = topK
	st.maxWeight = MIN_WEIGHT
	st.pq = NewPriorityQueue(0, topK)

	src := g.Verteces[srcId]
	st.initMemo()
	st.processOptEdges(src, src, 0)
	st.traceMemo(src, src)
	res = processOutsideEdges(st.pq, st.deepLimit, topK, false)
	return
}

func (st *DfsMemo) SetTreshold(t float64) {
	st.threshold = t
	st.thresholdMode = THR_CUSTOM
}

func (st *DfsMemo) SetTresholdMode(mode int) {
	st.thresholdMode = mode
}

func (st *DfsMemo) prepareThreshold(src int, target int, level int) {
	switch st.thresholdMode {
	case THR_ZERO:
		st.threshold = 0
	case THR_MEAN:
		nodes := st.memo[src][target][level]
		st.threshold = nodes.stat.meanWeight
	case THR_MEAN_STDDEV:
		nodes := st.memo[src][target][level]
		st.threshold = nodes.stat.meanWeight - math.Sqrt((nodes.stat.mean2Weight-nodes.stat.meanWeight*nodes.stat.meanWeight)/nodes.stat.pathsCount)
	}
}

func (st *DfsMemo) processOptEdges(src int, target int, level int) (bool, TreeItemStat) {
	if level >= st.deepLimit {
		return false, TreeItemStat{
			minWeight:   MIN_WEIGHT,
			maxWeight:   MAX_WEIGHT,
			meanWeight:  0,
			mean2Weight: 0,
			pathsCount:  0}
	}

	if item := st.inMemo(src, target, level); item != nil {
		return true, item.stat
	}

	res := make([]*TreeNode, 0)
	stat := TreeItemStat{
		minWeight:   MAX_WEIGHT,
		maxWeight:   MIN_WEIGHT,
		meanWeight:  0,
		mean2Weight: 0,
		pathsCount:  0,
	}

	for _, edge := range st.g.Succ(src) {
		v := edge.data.Id2
		if target == v {
			res = append(res, NewTreeNode(edge, -1, -1, -1))

			weight := edge.weight
			stat.meanWeight += weight
			stat.mean2Weight += weight * weight
			stat.pathsCount += 1

			if weight < stat.minWeight {
				stat.minWeight = weight
			}

			if weight > stat.maxWeight {
				stat.maxWeight = weight
			}

			continue
		}

		ok, statw := st.processOptEdges(v, target, level+1)

		if ok {
			res = append(res, NewTreeNode(edge, v, target, level+1))

			minw := statw.minWeight + edge.weight
			maxw := statw.maxWeight + edge.weight

			if minw < stat.minWeight {
				stat.minWeight = minw
			}

			if maxw > stat.maxWeight {
				stat.maxWeight = maxw
			}

			stat.meanWeight += statw.meanWeight + statw.pathsCount*edge.weight
			stat.mean2Weight += statw.mean2Weight + statw.pathsCount*(edge.weight*edge.weight)
			stat.pathsCount += statw.pathsCount
		}
	}

	st.toMemo(&TreeItem{seq: res, stat: stat}, src, target, level)
	return true, stat
}

func (st *DfsMemo) initMemo() {
	st.memo = make(map[int]map[int]map[int]*TreeItem)
}

func (st *DfsMemo) inMemo(src int, target int, level int) *TreeItem {
	if vs, ok := st.memo[src]; ok {
		if vt, ok := vs[target]; ok {
			if vl, ok := vt[level]; ok {
				return vl
			}
		}
	}

	return nil
}

func (st *DfsMemo) toMemo(res *TreeItem, src int, target int, level int) {
	vs, ok := st.memo[src]

	if !ok {
		st.memo[src] = make(map[int]map[int]*TreeItem)
		vs = st.memo[src]
	}

	vt, ok := vs[target]

	if !ok {
		vs[target] = make(map[int]*TreeItem)
		vt = vs[target]
	}

	vt[level] = res
}

func (st *DfsMemo) traceMemo(src int, target int) {
	levels := st.memo[src][target]

	for level, _ := range levels {
		st.prepareThreshold(src, target, level)
		st.nextMemoItem(src, target, level)
	}
}

func (st *DfsMemo) nextMemoItem(src int, target int, level int) {
	if src < 0 || target < 0 || level < 0 {
		return
	}

	nodes := st.memo[src][target][level]

	if st.psa[level]+nodes.stat.minWeight >= st.threshold {
		return
	}

	if st.maxWeight != MIN_WEIGHT && st.psa[level]+nodes.stat.minWeight > st.maxWeight {
		return
	}

	for _, node := range nodes.seq {
		st.edges[level] = node.base
		weight := st.psa[level] + node.base.weight
		st.psa[level+1] = weight

		if node.src < 0 {
			if weight >= 0 {
				continue
			}

			for i := level + 1; i < st.deepLimit; i++ {
				st.edges[i] = nil
			}

			if st.pq.Len() < st.topK {
				cedges := make(MEdgeSeq, st.deepLimit)
				copy(cedges, st.edges)
				st.pq.Append(cedges, weight)

				if st.pq.Len() == st.topK {
					st.pq.Init()
					st.maxWeight = st.pq[0].priority
				}
				continue
			}

			if weight < st.maxWeight {
				ms, _ := st.pq[0].value.(MEdgeSeq)
				copy(ms, st.edges)
				st.pq.Update(st.pq[0], st.pq[0].value, weight)
				st.maxWeight = st.pq[0].priority
			}

			continue
		}

		st.nextMemoItem(node.src, node.target, node.level)
	}

	if st.pq.Len() < st.topK {
		st.pq.Init()
	}
}
