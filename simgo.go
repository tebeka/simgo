package simgo

import (
	"container/heap"
)

var (
	Version = "0.1.0"
)

type Worker interface {
	Chan() chan int
}

type worker struct {
	job  Worker
	time int
}

type workerHeap []*worker

func (h workerHeap) Len() int           { return len(h) }
func (h workerHeap) Less(i, j int) bool { return h[i].time < h[j].time }
func (h workerHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *workerHeap) Push(x interface{}) {
	*h = append(*h, x.(*worker))
}

func (h *workerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Env struct {
	wh  *workerHeap
	now int
}

func NewEnv() *Env {
	return &Env{&workerHeap{}, 0}
}

func (env *Env) Process(w Worker) {
	wk := &worker{job: w, time: env.now}
	heap.Push(env.wh, wk)
}

func (env *Env) Run(until int) {
	for env.wh.Len() > 0 {
		w := heap.Pop(env.wh).(*worker)
		if w.time > until {
			break
		}
		env.now = w.time
		w.job.Chan() <- env.now
		w.time = <-w.job.Chan()
		if w.time > 0 {
			heap.Push(env.wh, w)
		}
	}
}
