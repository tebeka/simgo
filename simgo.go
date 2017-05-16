// simgo is discreet event simulation framework
package simgo

import (
	"container/heap"
)

var (
	Version = "0.1.0"
)

type InChan chan int  // Process input channel
type OutChan chan int // Process output channel

// Process interface
type Process interface {
	In() InChan   // Input channel, receives time
	Out() OutChan // Output channel, send sleep time in ticks, -1 for done
}

type worker struct {
	proc Process
	time int
}

// container/heap interface
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

// Simulation environment
type Env struct {
	wh  *workerHeap
	now int
}

// NewEnv returns a new environment
func NewEnv() *Env {
	return &Env{&workerHeap{}, 0}
}

// Process adds a process to the simulation
func (env *Env) Process(w Process) {
	wk := &worker{proc: w, time: env.now}
	heap.Push(env.wh, wk)
}

// batch returns batch of processes to start now
func (env *Env) batch(time int) []*worker {
	var batch []*worker
	for env.nextTick() == time {
		w := heap.Pop(env.wh).(*worker)
		batch = append(batch, w)
	}
	return batch
}

// nextTick return next tick to jump to
func (env *Env) nextTick() int {
	if env.wh.Len() == 0 {
		return -1
	}
	return (*env.wh)[0].time
}

// Run runs the simulation
func (env *Env) Run(until int) {
	for env.wh.Len() > 0 {
		now := env.nextTick()
		if now > until {
			env.now = until
			return
		}
		env.now = now
		batch := env.batch(now)
		// Start all
		for _, w := range batch {
			w.proc.In() <- now
		}

		// Wait, if worker return positive number - it's sleep time
		for _, w := range batch {
			dt := <-w.proc.Out()
			if dt > 0 {
				w.time += dt
				heap.Push(env.wh, w)
			}
		}
	}
}

// Now returns the current tick
func (env *Env) Now() int {
	return env.now
}
