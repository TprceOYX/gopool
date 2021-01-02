package main

type Task interface {
	ExecTask()
	PanicHandler()
}

type TaskQueue interface {
	Enqueue(t Task)
	Dequeue() (t Task)
	Size() int32
}
