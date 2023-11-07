package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
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

	var (
		totalDuration time.Duration
		minDuration   = time.Duration(math.MaxInt64)
		maxDuration   time.Duration
		taskCount     int
	)

	for _, taskGroup := range tasks {
		for _, task := range taskGroup {
			started, errStart := time.Parse(time.RFC3339Nano, task.StartedAt)
			created, errCreated := time.Parse(time.RFC3339Nano, task.CreatedAt)
			if errStart != nil || errCreated != nil {
				fmt.Println("Error parsing dates:", errStart, errCreated)
				continue // Skip this task if there's an error parsing dates
			}
			duration := started.Sub(created)

			if duration < minDuration {
				minDuration = duration
			}
			if duration > maxDuration {
				maxDuration = duration
			}
			totalDuration += duration
			taskCount++
		}
	}

	if taskCount == 0 {
		fmt.Println("No tasks to process.")
		return
	}

	averageDuration := totalDuration / time.Duration(taskCount)
	fmt.Println("Average duration:", averageDuration)
	fmt.Println("Minimum duration:", minDuration)
	fmt.Println("Maximum duration:", maxDuration)
}
