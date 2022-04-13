package kspa

import "math/big"

type BiBfs struct{}

func (b *BiBfs) ForwardBackward(dg *DexGraph, source Dexer, middle Dexer, target Dexer, depthLimit int) (forw []*DexTreeNode, back []*DexTreeNode) {
	bfs := new(Bfs)

	forw = bfs.Forward(dg, source, middle, depthLimit)
	back = bfs.Backward(dg, target, middle, depthLimit)
	return
}

func (b *BiBfs) ForwardBackwardArbitrage(dg *DexGraph, source Dexer, middle Dexer, depthLimit int) (forw []*DexTreeNode, back []*DexTreeNode) {
	queue := make([]*DexTreeNode, 1)
	forw = make([]*DexTreeNode, 0)
	back = make([]*DexTreeNode, 0)

	queue = append(queue, NewDexTreeNode(source, 0, nil))
	tokenOutIdForward := middle.TokenInId()
	tokenOutIdBackward := middle.TokenOutId()

	for len(queue) != 0 {
		var x *DexTreeNode
		x, queue = queue[0], queue[1:]

		if x.Level > depthLimit {
			break
		}

		for _, s := range dg.Succ(x.Val.TokenOutId()) {
			node := NewDexTreeNode(s, x.Level+1, x)

			if s.TokenOutId() == tokenOutIdForward {
				forw = append(forw, node)
				continue
			}

			if s.TokenOutId() == tokenOutIdBackward {
				back = append(back, node)
				continue
			}

			queue = append(queue, node)
		}
	}

	return
}

type BruteForceConfig struct {
	StartTokenInAmount *big.Int
	EndTokenInAmount *big.Int
	Step *big.Int
}

type LimitOrder struct {
	tokenInAmount *big.Int
	tokenOutAmount *big.Int
}

func AllPaths(cfg *BruteForceConfig, lo *LimitOrder, dex Dexer, forw []*DexTreeNode, back []*DexTreeNode) {
	startTokenInAmount := 
	for _, f := range forw {

		node := f
		out :=
		for node != nil {
			in := node.Val.GetAmountOutByAmountIn()
			node = node.Parent
		}

		for _, b := range back {
			node = b
			for node != nil {
				node = node.Parent
			}
		}
	}
}
