package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func()

type Queue struct {
	tasks []Task
	mu    sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		tasks: make([]Task, 0),
	}
}

func (q *Queue) AddTask(t Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, t)
}

func (q *Queue) Run() {
	for len(q.tasks) > 0 {
		q.mu.Lock()
		t := q.tasks[0]
		q.tasks = q.tasks[1:]
		q.mu.Unlock()
		t()
	}
}

type Scheduler struct {
	queues []*Queue
	wg     sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		queues: make([]*Queue, 0),
	}
}

func (s *Scheduler) AddQueue(q *Queue) {
	s.queues = append(s.queues, q)
}

func (s *Scheduler) Run() {
	for _, q := range s.queues {
		s.wg.Add(1)
		go func(q *Queue) {
			defer s.wg.Done()
			q.Run()
		}(q)
	}
	s.wg.Wait()
}

func main() {
	s := NewScheduler()

	q1 := NewQueue()
	q1.AddTask(func() {
		time.Sleep(time.Second * 3)
		fmt.Println("Task 1 from Queue 1")
	})
	q1.AddTask(func() {
		time.Sleep(time.Second * 2)
		fmt.Println("Task 2 from Queue 1")
	})

	q2 := NewQueue()
	q2.AddTask(func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Task 1 from Queue 2")
	})
	q2.AddTask(func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Task 2 from Queue 2")
	})

	s.AddQueue(q1)
	s.AddQueue(q2)

	s.Run()
}
