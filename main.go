package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Task struct {
	StartedAt string `json:"startedAt"`
	CreatedAt string `json:"createdAt"`
	TaskArn   string `json:"taskArn"`
}

func main() {
	// Read from standard input
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading from stdin:", err)
		return
	}

	var tasks [][]Task
	err = json.Unmarshal([]byte(inputData), &tasks)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var totalDuration time.Duration
	for _, taskGroup := range tasks {
		for _, task := range taskGroup {
			started, _ := time.Parse(time.RFC3339Nano, task.StartedAt)
			created, _ := time.Parse(time.RFC3339Nano, task.CreatedAt)
			duration := started.Sub(created)
			totalDuration += duration
		}
	}

	averageDuration := totalDuration / time.Duration(len(tasks)*len(tasks[0]))
	fmt.Println("Average duration:", averageDuration)
}
