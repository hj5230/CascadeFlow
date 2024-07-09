# CascadeFlow

A simple and efficient multi-stage task executor written in Go.

## Features:

Defines a clear task structure and execution flow.

Supports multiple execution queues, each handling a specific type of task.

Ensures sequential execution of tasks within a queue and chained execution across queues.

Utilizes goroutines and channels for asynchronous task execution.

Provides thread-safety through mutexes and wait groups.

## How it works:

The core components of this library are:

Task: Represents a unit of work with an ID and data.

ExecuteQueue: Manages a queue of tasks and executes them using a defined function.

Executor: Orchestrates the execution flow across multiple queues.

Tasks are added to specific queues within the Executor. Each queue processes tasks sequentially, and the output of one queue can be fed as input to the next, creating a pipeline-like execution flow.

## Example Usage:

```go
package main

import (
	"fmt"
	"time"

	"github.com/hj5230/CascadeFlow/cascadeflow"
)

func main() {
	// Define task execution functions
	addFunc := func(data map[string]interface{}) map[string]interface{} {
		// ... task logic ...
		return map[string]interface{}{"result": a + b}
	}

	numToStringFunc := func(data map[string]interface{}) map[string]interface{} {
		// ... task logic ...
		return map[string]interface{}{"result": fmt.Sprintf("%d", result)}
	}

	// Create an Executor with multiple queues
	executor := staged-executor.NewExecutor(
		"sampleExecutor",		[]interface{}{
			staged-executor.NewExecuteQueue("Queue1", 3, addFunc),
			staged-executor.NewExecuteQueue("Queue2", 3, numToStringFunc),
		},
	)

	// Add tasks to the executor
	executor.AddTask(&staged-executor.Task{ID: 1, Data: map[string]interface{}{"a": 1, "b": 2}})
	executor.AddTask(&staged-executor.Task{ID: 2, Data: map[string]interface{}{"a": 3, "b": 4}})

	// Execute tasks
	executor.Execute()
}
```
