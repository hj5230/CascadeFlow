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

func (eq *ExecuteQueue) ExecuteNext() (map[string]interface{}, bool) {
	eq.mu.Lock()
	if len(eq.tasks) == 0 {
		eq.mu.Unlock()
		return nil, false
	}
	task := eq.tasks[0]
	eq.tasks = eq.tasks[1:]
	eq.mu.Unlock()

	resultChan := make(chan map[string]interface{})
	eq.wg.Add(1)
	go func() {
		defer eq.wg.Done()
		result := eq.executeFunc(task.Data)
		resultChan <- result
	}()

	result := <-resultChan
	return result, true
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
	for _, queue := range e.queues {
		queue.mu.Lock()
		if len(queue.tasks) < queue.maxSize {
			queue.tasks = append(queue.tasks, task)
			queue.mu.Unlock()
			return true
		}
		queue.mu.Unlock()
	}
	return false
}

func (e *Executor) Execute() {
	var wg sync.WaitGroup
	for {
		allEmpty := true
		for i, queue := range e.queues {
			if result, ok := queue.ExecuteNext(); ok {
				allEmpty = false
				wg.Add(1)
				go func(i int, result map[string]interface{}) {
					defer wg.Done()
					if i+1 < len(e.queues) {
						e.queues[i+1].mu.Lock()
						e.queues[i+1].tasks = append(e.queues[i+1].tasks, &Task{ID: -1, Data: result})
						e.queues[i+1].mu.Unlock()
					}
				}(i, result)
			}
		}
		wg.Wait()
		if allEmpty {
			break
		}
	}

	for _, queue := range e.queues {
		queue.wg.Wait()
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
