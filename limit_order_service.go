package kspa

import "math"

// TODO
type LimitOrderService struct {
	g           *MultiGraph
	unprocessed EdgeSeq
}

func (s *LimitOrderService) Add(id string, tokenIn int, tokenOut int, tokenInAmount float64, tokenOutAmount float64) {
	relation := tokenOutAmount / tokenInAmount
	s.AddWrapped(&SingleEdge{
		Data: &Entity{
			EntityId: id,
			Id1:      tokenIn,
			Id2:      tokenOut,
			Relation: relation,
		},
		Weight: -math.Log(relation),
		Status: LIMIT_ORDER,
	})
}

func (s *LimitOrderService) AddWrapped(lo *SingleEdge) {
	s.unprocessed = append(s.unprocessed, lo)
}

func (s *LimitOrderService) AddWrappedList(los EdgeSeq) {
	s.unprocessed = append(s.unprocessed, los...)
}

func (s *LimitOrderService) Remove(id string) {
	// TODO
}
