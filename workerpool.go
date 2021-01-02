package main

import (
	"sync/atomic"
)

type WorkerPool struct {
	// 最大工作线程数
	maxWorkerCount int32

	// 当前工作线程数
	workersNum int32
	queue      TaskQueue
	close      int32
}

func NewWorkerPool(maxWorkerCount int32) *WorkerPool {
	if maxWorkerCount < 0 {
		panic("max worker count must bigger than 0")
	}
	wp := &WorkerPool{
		maxWorkerCount: maxWorkerCount,
	}
	return wp.init()
}

func (wp *WorkerPool) init() *WorkerPool {
	wp.queue = NewTaskQueue()
	return wp
}

func (wp *WorkerPool) Run(t Task) {
	if atomic.LoadInt32(&wp.close) > 0 {
		panic("workerpool has been closed")
	}
	wp.queue.Enqueue(t)
	if wp.getWorker() {
		wp.incrWorkersNum()
		go func() {
			wp.workerFunc()
			wp.decrWorkersNum()
		}()
	}
}

func (wp *WorkerPool) Shutdowm() {
	atomic.StoreInt32(&wp.close, 1)
}

func (wp *WorkerPool) getWorker() bool {
	// 可能出现实际启动线程数量大于maxWorkerCount的情况
	return atomic.LoadInt32(&wp.workersNum) < wp.maxWorkerCount
}

func (wp *WorkerPool) workerFunc() {
	for {
		t := wp.queue.Dequeue()
		if t == nil {
			return
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.PanicHandler()
				}
			}()
			t.ExecTask()
		}()
	}
}

func (wp *WorkerPool) incrWorkersNum() {
	atomic.AddInt32(&wp.workersNum, 1)
}

func (wp *WorkerPool) decrWorkersNum() {
	atomic.AddInt32(&wp.workersNum, -1)
}
