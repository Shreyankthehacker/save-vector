package indexing




type FloatHeap []float32

func (f FloatHeap) Len() int {
	return len(f)
}
func (f FloatHeap) Less(i, j int) bool {
	return f[i] > f[j]
}

func (f FloatHeap) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f *FloatHeap) Push(x interface{}) {
	*f = append(*f, x.(float32))
}

func (f *FloatHeap) Pop() interface{} {
	old := *f
	n := len(old)
	x := old[n-1]
	*f = old[:n-1]
	return x
}