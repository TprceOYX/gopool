package main

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// basic test
var a int32

type normalTask struct {
	wg *sync.WaitGroup
}

func (t *normalTask) ExecTask() {
	atomic.AddInt32(&a, 1)
	t.wg.Done()
}

func (t *normalTask) PanicHandler() {}

type panicTask struct {
	wg *sync.WaitGroup
}

func (t *panicTask) ExecTask() {
	panic("task panic")
}

func (t *panicTask) PanicHandler() {
	atomic.AddInt32(&a, 1)
	// fmt.Println("a")
	t.wg.Done()
}

func TestWorkerPool(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1000)
	wp := NewWorkerPool(int32(runtime.GOMAXPROCS(0)))
	a = 0
	for i := 0; i < 1000; i++ {
		var t Task
		if i%2 == 0 {
			t = &normalTask{
				wg: wg,
			}
		} else {
			t = &panicTask{
				wg: wg,
			}
		}
		wp.Run(t)
	}
	wg.Wait()
	if a != 1000 {
		t.Errorf("error value, a:%v", a)
	}
}

// benchmark test
const benchmarkCount = 10000

type benchmarkTask struct {
	wg *sync.WaitGroup
}

func (t *benchmarkTask) ExecTask() {
	rand.Intn(benchmarkCount)
	// fmt.Println("a")
	t.wg.Done()
}

func (t *benchmarkTask) PanicHandler() {}

func BenchmarkWorkerPool(b *testing.B) {
	wp := NewWorkerPool(int32(runtime.GOMAXPROCS(0)))
	wg := &sync.WaitGroup{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkCount)
		for n := 0; n < benchmarkCount; n++ {
			t := &benchmarkTask{
				wg: wg,
			}
			wp.Run(t)
		}
		wg.Wait()
	}
}
