package cascadeflow

import (
	"fmt"
	"sync"
)

type Task struct {
	ID   int
	Data map[string]interface{}
}

type ExecuteFunc func(map[string]interface{}) map[string]interface{}

type ExecuteQueue struct {
	name        string
	maxSize     int
	tasks       []*Task
	executeFunc ExecuteFunc
	mu          sync.Mutex
	wg          sync.WaitGroup
}

func NewExecuteQueue(name string, maxSize int, executeFunc ExecuteFunc) *ExecuteQueue {
	return &ExecuteQueue{
		name:        name,
		maxSize:     maxSize,
		tasks:       make([]*Task, 0),
		executeFunc: executeFunc,
	}
}

func (eq *ExecuteQueue) AddTask(task *Task) bool {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if len(eq.tasks) < eq.maxSize {
		eq.tasks = append(eq.tasks, task)
		return true
	}
	return false
}

func (eq *ExecuteQueue) ExecuteTasks(nextQueue *ExecuteQueue) {
	eq.mu.Lock()
	tasks := eq.tasks
	eq.tasks = nil
	eq.mu.Unlock()

	for _, task := range tasks {
		eq.wg.Add(1)
		go func(t *Task) {
			defer eq.wg.Done()
			result := eq.executeFunc(t.Data)
			if nextQueue != nil {
				nextQueue.AddTask(&Task{ID: t.ID, Data: result})
			}
		}(task)
	}

	eq.wg.Wait()
}

type Executor struct {
	name   string
	queues []*ExecuteQueue
}

func NewExecutor(name string, queues interface{}) *Executor {
	queueSlice, ok := queues.([]interface{})
	if !ok {
		panic("Invalid queues parameter")
	}

	var executorQueues []*ExecuteQueue
	for _, q := range queueSlice {
		if eq, ok := q.(*ExecuteQueue); ok {
			executorQueues = append(executorQueues, eq)
		}
	}

	return &Executor{
		name:   name,
		queues: executorQueues,
	}
}

func (e *Executor) AddTask(task *Task) bool {
	return e.queues[0].AddTask(task)
}

func (e *Executor) Execute() {
	for i := 0; i < len(e.queues); i++ {
		nextQueue := (*ExecuteQueue)(nil)
		if i+1 < len(e.queues) {
			nextQueue = e.queues[i+1]
		}
		e.queues[i].ExecuteTasks(nextQueue)
	}
}

func (e *Executor) GetAllQueueStatus() string {
	status := fmt.Sprintf("%s@%p{", e.name, e)
	for i, queue := range e.queues {
		queue.mu.Lock()
		status += fmt.Sprintf("%s@%p{", queue.name, queue)
		for j, task := range queue.tasks {
			if j > 0 {
				status += ", "
			}
			status += fmt.Sprintf("task%d@%p", task.ID, task)
		}
		status += "}"
		queue.mu.Unlock()
		if i < len(e.queues)-1 {
			status += ", "
		}
	}
	status += "}"
	return status
}
