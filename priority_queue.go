package graph_shortest_paths

import (
	"container/heap"
	"encoding/json"
)

type Item struct {
	value    interface{}
	priority float64
	index    int
}

func (pq Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pq.value)
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item, _ := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Append(value interface{}, priority float64) {
	item := &Item{value: value, priority: priority}
	pq.Push(item)
}

func (pq *PriorityQueue) Update(item *Item, value interface{}, priority float64) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) Init() {
	heap.Init(pq)
}

func (pq *PriorityQueue) PushItem(item *Item) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue) PopItem() (item *Item) {
	item = heap.Pop(pq).(*Item)
	return
}

func NewPriorityQueue(size int, capacity int) PriorityQueue {
	pq := make(PriorityQueue, size, capacity)
	return pq
}

func PriorityQueue2SortedArray(pq PriorityQueue, asc bool) (sa PriorityQueue) {
	n := pq.Len()
	sa = NewPriorityQueue(pq.Len(), pq.Len())
	index := 0
	inc := +1

	if !asc {
		index = n - 1
		inc = -1
	}

	for pq.Len() > 0 {
		item := pq.PopItem()
		item.index = index
		sa[index] = item
		index += inc
	}

	return
}
