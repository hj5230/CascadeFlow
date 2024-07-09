package main

import (
	"fmt"
	"time"

	"github.com/hj5230/CascadeFlow/cascadeflow"
)

func main() {
	addFunc := func(data map[string]interface{}) map[string]interface{} {
		fmt.Printf("Adding %d and %d\n", data["a"], data["b"])
		time.Sleep(3 * time.Second)
		a := data["a"].(int)
		b := data["b"].(int)
		return map[string]interface{}{"result": a + b}
	}

	numToStringFunc := func(data map[string]interface{}) map[string]interface{} {
		fmt.Printf("Converting %d to string\n", data["result"])
		time.Sleep(5 * time.Second)
		result := data["result"].(int)
		return map[string]interface{}{"result": fmt.Sprintf("%d", result)}
	}

	printOutputFunc := func(data map[string]interface{}) map[string]interface{} {
		fmt.Println("Final output:", data["result"])
		return data
	}

	executor := cascadeflow.NewExecutor(
		"sampleExecutor",
		[]interface{}{
			cascadeflow.NewExecuteQueue("Queue1", 3, addFunc),
			cascadeflow.NewExecuteQueue("Queue2", 3, numToStringFunc),
			cascadeflow.NewExecuteQueue("Queue3", 3, printOutputFunc),
		},
	)

	executor.AddTask(&cascadeflow.Task{ID: 1, Data: map[string]interface{}{"a": 1, "b": 2}})
	executor.AddTask(&cascadeflow.Task{ID: 2, Data: map[string]interface{}{"a": 3, "b": 4}})
	executor.AddTask(&cascadeflow.Task{ID: 3, Data: map[string]interface{}{"a": 5, "b": 6}})

	fmt.Println(executor.GetAllQueueStatus())

	executor.Execute()
}
