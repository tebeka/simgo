package simgo

import (
	"container/heap"

	"testing"
)

func TestHeap(t *testing.T) {
	h := &workerHeap{}
	heap.Init(h)
	heap.Push(h, &worker{nil, 2})
	heap.Push(h, &worker{nil, 1})
	heap.Push(h, &worker{nil, 3})

	if h.Len() != 3 {
		t.Fatalf("wrong size (%d != 3)", h.Len())
	}

	w := heap.Pop(h).(*worker)
	if h.Len() != 2 {
		t.Fatalf("wrong size (%d != 2)", h.Len())
	}

	if w.time != 1 {
		t.Fatalf("bad time (%d != 1)", w.time)
	}
}

type testWorker struct {
	env *Env
	ch  chan int
	id  int
}

func (tw *testWorker) Chan() chan int {
	return tw.ch
}

func newTestWorker(env *Env, id int, t *testing.T) *testWorker {
	tw := &testWorker{
		env: env,
		ch:  make(chan int),
		id:  id,
	}
	go func() {
		for {
			tick := <-tw.ch
			t.Logf("[%d] at %d\n", tw.id, tick)
			tw.ch <- tick + (tw.id+1)*10
		}
	}()
	return tw
}

func TestProcess(t *testing.T) {
	env := NewEnv()
	tw := newTestWorker(env, 0, t)

	env.Process(tw)
}

func TestRun(t *testing.T) {
	env := NewEnv()
	for i := 0; i < 3; i++ {
		env.Process(newTestWorker(env, i, t))
	}
	env.Run(100)
}
