package data_structures

// Heap
// the implement in Golang is a Min Heap.
// DO NOT use the methods of Push and Pop we write below. Import 'heap' and Use heap.Init,heap.Push,heap.Pop instead.
type Heap []int

func (h *Heap) Len() int           { return len(*h) }
func (h *Heap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *Heap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *Heap) Pop() (v interface{}) {
	v, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]
	return v
}
