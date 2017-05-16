package simgo

import (
	"log"
	"os"

	"testing"
)

var (
	verbose = false
)

func init() {
	if len(os.Getenv("SIMGO_VERBOSE")) > 0 {
		verbose = true
	}
}

type testWorker struct {
	env   *Env
	ch    chan int
	id    int
	ticks []int
}

func (tw *testWorker) In() InChan {
	return tw.ch
}

func (tw *testWorker) Out() OutChan {
	return tw.ch
}

func newTestWorker(id int, t *testing.T) *testWorker {
	tw := &testWorker{
		ch: make(chan int),
		id: id,
	}
	go func() {
		for {
			tick := <-tw.ch
			tw.ticks = append(tw.ticks, tick)
			if verbose {
				t.Logf("[%d] at %d\n", tw.id, tick)
			}
			tw.ch <- (tw.id + 1) * 10
		}
	}()
	return tw
}

func sliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func TestProcess(t *testing.T) {
	env := NewEnv()
	tw := newTestWorker(0, t)

	env.Process(tw)
}

func TestRun(t *testing.T) {
	env := NewEnv()
	workers := make([]*testWorker, 3)
	for i := 0; i < 3; i++ {
		workers[i] = newTestWorker(i, t)
		env.Process(workers[i])
	}
	until := 100
	env.Run(until)

	if env.Now() != until {
		log.Fatalf("bad until %d != 100\n", env.Now())
	}

	expected := [][]int{
		{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		{0, 20, 40, 60, 80, 100},
		{0, 30, 60, 90},
	}

	for i, slice := range expected {
		if !sliceEq(workers[i].ticks, slice) {
			t.Fatalf("bad ticks for %d: %v\n", i, workers[i].ticks)
		}
	}
}
