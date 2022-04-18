package kspa

const (
	THR_NOT_USED = iota
	THR_CUSTOM
	THR_ZERO
)

const (
	LO_NOT_IN_CHAIN = iota
	LO_PARENT
	LO_CHILD
	LO_CURRENT
)

const (
	FN_ALL_PATHS = iota
	FN_LO_ONLY
)

type TreeNode struct {
	base   *MultiEdge
	src    int
	target int
	level  int
	status int
}

func NewTreeNode(base *MultiEdge, src int, target int, level int, status int) *TreeNode {
	return &TreeNode{
		base:   base,
		src:    src,
		target: target,
		level:  level,
		status: status,
	}
}

type TreeItemStat struct {
	minWeight float64
	status    int
}

type TreeItem struct {
	stat TreeItemStat
	seq  []*TreeNode
}

type DfsMemo struct {
	Dfs
	memo          map[int]map[int]map[int]*TreeItem
	threshold     float64
	thresholdMode int
	fnMode        int
	status        []int
}

func (st *DfsMemo) SetDeepLimit(v int) {
	st.Dfs.SetDeepLimit(v)
	st.status = make([]int, st.deepLimit+2)
	st.status[0] = LO_CHILD
}

func (st *DfsMemo) SetTreshold(t float64) {
	st.threshold = t
	st.thresholdMode = THR_CUSTOM
}

func (st *DfsMemo) SetTresholdMode(mode int) {
	st.thresholdMode = mode
}

func (st *DfsMemo) SetFnMode(mode int) {
	st.fnMode = mode
}

func (st *DfsMemo) Init() {
	st.initMemo()
}

func (st *DfsMemo) TopK(g *MultiGraph, srcId int, targetId int, topK int) (res PriorityQueue) {
	st.g = g
	st.SetTopKValue(topK)

	onlyLimitOrders := false
	if FN_LO_ONLY == st.fnMode {
		onlyLimitOrders = true
	}

	src := g.VertexIndex[srcId]
	target := g.VertexIndex[targetId]
	st.fillMemo(src, target, 0)
	st.traceMemo(src, target)
	res = ProcessOutsideEdges(st.pq, st.deepLimit, topK, false, onlyLimitOrders)
	return
}

func (st *DfsMemo) TopKOneToOne(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	st.g = g

	n := len(srcIds)
	res = make([]PriorityQueue, 0)

	onlyLimitOrders := false
	if FN_LO_ONLY == st.fnMode {
		onlyLimitOrders = true
	}

	for i := 0; i < n; i++ {
		src := g.VertexIndex[srcIds[i]]
		target := g.VertexIndex[targetIds[i]]

		st.SetTopKValue(topK)
		st.fillMemo(src, target, 0)
		st.traceMemo(src, target)
		res = append(res, ProcessOutsideEdges(st.pq, st.deepLimit, topK, false, onlyLimitOrders))
	}

	return
}

func (st *DfsMemo) TopKOneToMany(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue) {
	st.g = g

	n := len(srcIds)
	m := len(targetIds)

	res = make([]PriorityQueue, 0)

	onlyLimitOrders := false
	if FN_LO_ONLY == st.fnMode {
		onlyLimitOrders = true
	}

	for i := 0; i < n; i++ {
		st.SetTopKValue(topK)

		src := g.VertexIndex[srcIds[i]]

		for j := 0; j < m; j++ {
			target := g.VertexIndex[targetIds[j]]
			st.fillMemo(src, target, 0)
			st.traceMemo(src, target)
		}

		res = append(res, ProcessOutsideEdges(st.pq, st.deepLimit, topK, false, onlyLimitOrders))
	}

	return
}

func (st *DfsMemo) SetGraph(g *MultiGraph) {
	st.g = g
}

func (st *DfsMemo) Arbitrage(srcIds []int, topK int) (res []PriorityQueue) {
	n := len(srcIds)
	res = make([]PriorityQueue, 0)

	onlyLimitOrders := false
	if FN_LO_ONLY == st.fnMode {
		onlyLimitOrders = true
	}

	for i := 0; i < n; i++ {
		src := st.g.VertexIndex[srcIds[i]]

		st.SetTopKValue(topK)
		st.fillMemo(src, src, 0)
		st.traceMemo(src, src)
		res = append(res, ProcessOutsideEdges(st.pq, st.deepLimit, topK, false, onlyLimitOrders))
	}

	return
}

func (st *DfsMemo) AddLimitOrders(edges EdgeSeq) (MEdgeSeq, error) {
	for _, edge := range edges {
		edge.Status = LIMIT_ORDER
	}
	return st.g.Add(edges, LIMIT_ORDER)
}

func (st *DfsMemo) RemoveLimitOrders(medges MEdgeSeq) {
	st.g.Remove(medges, UNDEFINED)
}

func (st *DfsMemo) prepareThreshold(src int, target int) {
	if THR_CUSTOM != st.thresholdMode {
		if src == target {
			st.thresholdMode = THR_ZERO
		} else {
			st.thresholdMode = THR_NOT_USED
		}
	}

	switch st.thresholdMode {
	case THR_ZERO:
		st.threshold = 0
	}
}

func (st *DfsMemo) fillMemo(src int, target int, level int) (bool, TreeItemStat) {
	if level >= st.deepLimit {
		return false, TreeItemStat{
			minWeight: MIN_WEIGHT,
			status:    LO_NOT_IN_CHAIN,
		}
	}

	if item := st.inMemo(src, target, level); item != nil {
		return true, item.stat
	}

	res := make([]*TreeNode, 0)
	stat := TreeItemStat{
		minWeight: MAX_WEIGHT,
		status:    LO_NOT_IN_CHAIN,
	}

	for _, edge := range st.g.Succ(src) {
		v := edge.Data.Id2i

		localStatus := LO_NOT_IN_CHAIN
		if LIMIT_ORDER == edge.status {
			stat.status = LO_CURRENT
			localStatus = LO_CURRENT
		}

		if target == v {
			res = append(res, NewTreeNode(edge, -1, -1, -1, localStatus))

			weight := edge.Weight

			if weight < stat.minWeight {
				stat.minWeight = weight
			}

			continue
		}

		ok, statw := st.fillMemo(v, target, level+1)

		if ok {
			if LO_CURRENT != localStatus && (LO_CURRENT == statw.status || LO_CHILD == statw.status) {
				localStatus = LO_CHILD

				if LO_NOT_IN_CHAIN == stat.status {
					stat.status = LO_CHILD
				}
			}

			res = append(res, NewTreeNode(edge, v, target, level+1, localStatus))

			minw := statw.minWeight + edge.Weight

			if minw < stat.minWeight {
				stat.minWeight = minw
			}
		}
	}

	st.toMemo(&TreeItem{seq: res, stat: stat}, src, target, level)
	return true, stat
}

func (st *DfsMemo) initMemo() {
	st.memo = nil
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
	st.prepareThreshold(src, target)
	st.nextMemoItem(src, target, 0)
}

func (st *DfsMemo) nextMemoItem(src int, target int, level int) {
	if src < 0 || target < 0 || level < 0 {
		return
	}

	nodes := st.memo[src][target][level]

	if FN_LO_ONLY == st.fnMode &&
		LO_NOT_IN_CHAIN == st.status[level] &&
		LO_NOT_IN_CHAIN == nodes.stat.status {
		return
	}

	if THR_NOT_USED != st.thresholdMode &&
		st.psa[level]+nodes.stat.minWeight >= st.threshold {
		return
	}

	if MIN_WEIGHT != st.maxWeight &&
		st.psa[level]+nodes.stat.minWeight > st.maxWeight {
		return
	}

	for _, node := range nodes.seq {
		st.edges[level] = node.base
		weight := st.psa[level] + node.base.Weight
		st.psa[level+1] = weight

		if FN_LO_ONLY == st.fnMode {
			switch st.status[level] {
			case LO_CURRENT:
				st.status[level+1] = st.status[level]
			case LO_CHILD, LO_NOT_IN_CHAIN:
				st.status[level+1] = node.status
			}

			if LO_NOT_IN_CHAIN == st.status[level+1] {
				continue
			}
		}

		if node.src < 0 {
			if THR_NOT_USED != st.thresholdMode && weight >= st.threshold {
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
					st.maxWeight = st.pq[0].Priority
				}
				continue
			}

			if weight < st.maxWeight {
				ms, _ := st.pq[0].Value.(MEdgeSeq)
				copy(ms, st.edges)
				st.pq.Update(st.pq[0], st.pq[0].Value, weight)
				st.maxWeight = st.pq[0].Priority
			}

			continue
		}

		st.nextMemoItem(node.src, node.target, node.level)
	}

	if st.pq.Len() < st.topK {
		st.pq.Init()
	}
}
