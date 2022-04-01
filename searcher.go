package kspa

type Searcher interface {
	TopK(g *MultiGraph, srcId int, targetId int, topK int) (res PriorityQueue)
	TopKOneToOne(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue)
	TopKOneToMany(g *MultiGraph, srcIds []int, targetIds []int, topK int) (res []PriorityQueue)
}
