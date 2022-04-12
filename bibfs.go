package kspa

type BiBfsProcessor interface {
	Run(*DexTreeNode, *DexTreeNode)
}

type BiBfs struct{}

func (b *BiBfs) Search(dg *DexGraph, source Dexer, middle Dexer, target Dexer, depthLimit int) (forw []*DexTreeNode, back []*DexTreeNode) {
	bfs := new(Bfs)

	forw = bfs.Forward(dg, source, middle, depthLimit)
	back = bfs.Backward(dg, target, middle, depthLimit)
	return
}

func (b *BiBfs) AllPaths(processor BiBfsProcessor, forw []*DexTreeNode, back []*DexTreeNode) {
	for _, f := range forw {
		for _, b := range back {
			processor.Run(f, b)
		}
	}
}
