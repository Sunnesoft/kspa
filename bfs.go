package kspa

type DexTreeNode struct {
	Val    Dexer
	Level  int
	Parent *DexTreeNode
}

func NewDexTreeNode(val Dexer, level int, parent *DexTreeNode) *DexTreeNode {
	return &DexTreeNode{
		Val:    val,
		Level:  level,
		Parent: parent,
	}
}

type Bfs struct{}

func (b *Bfs) Forward(dg *DexGraph, source Dexer, target Dexer, depthLimit int) []*DexTreeNode {
	queue := make([]*DexTreeNode, 1)
	res := make([]*DexTreeNode, 0)

	queue = append(queue, NewDexTreeNode(source, 0, nil))
	tokenOutId := target.TokenInId()

	for len(queue) != 0 {
		var x *DexTreeNode
		x, queue = queue[0], queue[1:]

		if x.Level > depthLimit {
			break
		}

		for _, s := range dg.Succ(x.Val.TokenOutId()) {
			node := NewDexTreeNode(s, x.Level+1, x)

			if s.TokenOutId() == tokenOutId {
				res = append(res, node)
				continue
			}

			queue = append(queue, node)
		}
	}

	return res
}

func (b *Bfs) Backward(dg *DexGraph, source Dexer, target Dexer, depthLimit int) []*DexTreeNode {
	queue := make([]*DexTreeNode, 1)
	res := make([]*DexTreeNode, 0)

	queue = append(queue, NewDexTreeNode(source, 0, nil))
	tokenOutId := target.TokenOutId()

	for len(queue) != 0 {
		var x *DexTreeNode
		x, queue = queue[0], queue[1:]

		if x.Level > depthLimit {
			break
		}

		for _, s := range dg.Succ(x.Val.TokenInId()) {
			node := NewDexTreeNode(s, x.Level+1, x)

			if s.TokenInId() == tokenOutId {
				res = append(res, node)
				continue
			}

			queue = append(queue, node)
		}
	}

	return res
}
