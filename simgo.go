package simgo

import (
	"container/heap"
)

var (
	Version = "0.1.0"
)

type InChan chan int
type OutChan chan int

type Worker interface {
	In() InChan
	Out() OutChan
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
			env.now = until
			return
		}
		env.now = w.time
		var batch []*worker
		// Start all ones for this time slot at the same time
		// (e.g. Don't wait for one to finish before starting another)
		for {
			w.job.In() <- env.now
			batch = append(batch, w)
			if env.wh.Len() == 0 {
				break
			}
			if (*env.wh)[0].time != env.now {
				break
			}
			w = heap.Pop(env.wh).(*worker)
		}

		for _, w := range batch {
			dt := <-w.job.Out()
			if dt > 0 {
				w.time += dt
				heap.Push(env.wh, w)
			}
		}
	}
}

func (env *Env) Now() int {
	return env.now
}
