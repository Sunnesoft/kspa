package graph_shortest_paths

import "math"

var MIN_WEIGHT float64 = math.Inf(-1)
var MAX_WEIGHT float64 = math.Inf(1)

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

type TreeItem struct {
	minWeight float64
	maxWeight float64
	seq       []*TreeNode
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

func (st *DfsMemo) processOptEdges(src int, target int, level int) (bool, float64, float64) {
	if level >= st.deepLimit {
		return false, MIN_WEIGHT, MAX_WEIGHT
	}

	if item := st.inMemo(src, target, level); item != nil {
		return true, item.minWeight, item.maxWeight
	}

	res := make([]*TreeNode, 0)
	minWeight := MAX_WEIGHT
	maxWeight := MIN_WEIGHT

	for _, edge := range st.g.Succ(src) {
		v := edge.data.Id2
		if target == v {
			res = append(res, NewTreeNode(edge, -1, -1, -1))

			weight := edge.weight

			if weight < minWeight {
				minWeight = weight
			}

			if weight > maxWeight {
				maxWeight = weight
			}

			continue
		}

		ok, minw, maxw := st.processOptEdges(v, target, level+1)

		if ok {
			res = append(res, NewTreeNode(edge, v, target, level+1))

			minw += edge.weight
			maxw += edge.weight

			if minw < minWeight {
				minWeight = minw
			}

			if maxw > maxWeight {
				maxWeight = maxw
			}
		}
	}

	st.toMemo(&TreeItem{seq: res, minWeight: minWeight, maxWeight: maxWeight}, src, target, level)
	return true, minWeight, maxWeight
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
		st.nextMemoItem(src, target, level)
	}
}

func (st *DfsMemo) nextMemoItem(src int, target int, level int) {
	if src < 0 || target < 0 || level < 0 {
		return
	}

	nodes := st.memo[src][target][level]

	if st.psa[level]+nodes.minWeight >= 0 {
		return
	}

	if st.maxWeight != MIN_WEIGHT && st.psa[level]+nodes.minWeight > st.maxWeight {
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
