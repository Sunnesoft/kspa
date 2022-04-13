package kspa

import "math/big"

type Dexer interface {
	Id() int
	GroupId() int
	SetId(int)
	SetGroupId(int)
	Alias() string
	TokenInId() int
	TokenOutId() int
	GetAmountOutByAmountIn(in *big.Int) *big.Int
	GetAmountInByAmountOut(out *big.Int) *big.Int
	GetRelation(in *big.Int) *big.Int
}

type DexChain []Dexer

func (dc DexChain) Id() int {
	return 0
}

func (dc DexChain) Alias() string {
	return ""
}

func (dc DexChain) TokenInId() int {
	return dc[0].TokenInId()
}

func (dc DexChain) TokenOutId() int {
	return dc[len(dc)-1].TokenOutId()
}

func (dc DexChain) GetAmountOutByAmountIn(in *big.Int) *big.Int {
	out := in
	for _, dex := range dc {
		out = dex.GetAmountOutByAmountIn(out)
	}
	return out
}

func (dc DexChain) GetAmountInByAmountOut(out *big.Int) *big.Int {
	in := out
	for i := len(dc) - 1; i >= 0; i-- {
		in = dc[i].GetAmountInByAmountOut(in)
	}
	return in
}

func (dc DexChain) GetRelation(in *big.Int) *big.Int {
	out := in
	res := big.NewInt(1)
	for _, dex := range dc {
		lout := dex.GetAmountOutByAmountIn(out)
		out = lout.Div(lout, out)
		res.Mul(res, out)
	}
	return res
}

func NewDexChain(size int, capacity int) DexChain {
	dc := make(DexChain, size, capacity)
	return dc
}

type DexBase struct {
	id         int
	alias      string
	tokenInId  int
	tokenOutId int
}

func (dc DexBase) Id() int {
	return dc.id
}

func (dc DexBase) Alias() string {
	return dc.alias
}

func (dc DexBase) TokenInId() int {
	return dc.tokenInId
}

func (dc DexBase) TokenOutId() int {
	return dc.tokenOutId
}

func (dc DexBase) GetAmountOutByAmountIn(in *big.Int) *big.Int {
	return in
}

func (dc DexBase) GetAmountInByAmountOut(out *big.Int) *big.Int {
	return out
}

func (dc DexBase) GetRelation(in *big.Int) *big.Int {
	lout := dc.GetAmountOutByAmountIn(in)
	return lout.Div(lout, in)
}

func NewDexBase(id int, alias string, tokenInId int, tokenOutId int) *DexBase {
	return &DexBase{
		id:         id,
		alias:      alias,
		tokenInId:  tokenInId,
		tokenOutId: tokenOutId,
	}
}

type DexGraph struct {
	edges      []Dexer
	successors [][]Dexer
}

func (dg *DexGraph) Build(edges []Dexer) {
	dg.edges = edges
	dg.initSuccessors()
}

func (dg *DexGraph) Succ(source int) []Dexer {
	return dg.successors[source]
}

func (dg *DexGraph) calcVertexCount() int {
	verteces := make(map[int]bool)

	for _, e := range dg.edges {
		verteces[e.TokenInId()] = true
		verteces[e.TokenOutId()] = true
	}

	return len(verteces)
}

func (dg *DexGraph) initSuccessors() {
	n := dg.calcVertexCount()
	dg.successors = make([][]Dexer, n)

	for _, v := range dg.edges {
		if dg.successors[v.TokenInId()] == nil {
			dg.successors[v.TokenInId()] = make([]Dexer, 0, 1)
		}

		dg.successors[v.TokenInId()] = append(dg.successors[v.TokenInId()], v)
	}
}
