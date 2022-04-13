package kspa

type Memo [][]*MemoNode

func NewMemo(depthLimit int, vertexCount int, targetCount int) {
	m := make(Memo, depthLimit)
	for i := 0; i < depthLimit; i++ {
		m[i] = make([]*MemoNode, edgesCount)
	}
}

type MemoNode struct {
	items []*MemoItem
}

func NewMemoNode(dexId int, level int, capacity int) *MemoNode {
	return &MemoNode{items: make([]*MemoItem, 0, capacity)}
}

type MemoItem struct {
	src    int
	target int
	level  int
}

func NewMemoItem(dexId int, level int) *MemoItem {
	return &MemoItem{
		dexId: dexId,
		level: level,
	}
}

type Dfsm struct {
	depthLimit int
}

func (st *Dfsm) fillMemo(src, target, level int) bool {
	if level >= st.depthLimit {
		return false
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
		v := edge.data.Id2i
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
