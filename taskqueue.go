package main

import (
	"sync"
	"sync/atomic"
)

type node struct {
	t    Task
	next *node
}

func (n *node) reset() *node {
	n.next = nil
	n.t = nil
	return n
}

/*
任务队列，任务给到线程池后放在任务队列中
*/
type taskQueue struct {
	head *node
	tail *node

	taskNum  int32
	nodePool sync.Pool
	lock     sync.Mutex
}

func NewTaskQueue() TaskQueue {
	t := &taskQueue{}
	t.nodePool.New = func() interface{} {
		return new(node)
	}
	return t
}

func (queue *taskQueue) Enqueue(t Task) {
	n := queue.nodePool.Get().(*node)
	n.t = t
	queue.lock.Lock()
	if queue.head == nil {
		queue.head = n
		queue.tail = n
	} else {
		queue.tail.next = n
		queue.tail = n
	}
	queue.taskNum++
	queue.lock.Unlock()
}

func (queue *taskQueue) Dequeue() (t Task) {
	queue.lock.Lock()
	head := queue.head
	if head != nil {
		t = queue.head.t
		if head.next != nil {
			queue.head = queue.head.next
		} else {
			queue.head, queue.tail = nil, nil
		}
		queue.nodePool.Put(head.reset())
	}
	queue.lock.Unlock()
	return
}

func (queue *taskQueue) Size() int32 {
	return atomic.LoadInt32(&queue.taskNum)
}
