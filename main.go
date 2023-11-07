package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Fatalf("Error reading from stdin: %s", err)
	}

	var tasks [][]Task
	err = json.Unmarshal([]byte(inputData), &tasks)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var (
		totalDuration time.Duration
		minDuration   = time.Duration(math.MaxInt64)
		maxDuration   time.Duration
		taskCount     int
	)

	for _, taskGroup := range tasks {
		for _, task := range taskGroup {
			started, err := time.Parse(time.RFC3339Nano, task.StartedAt)
			if err != nil {
				log.Fatalf("Error parsing dates for StartedAt: %s", err)
			}
			created, err := time.Parse(time.RFC3339Nano, task.CreatedAt)
			if err != nil {
				log.Fatalf("Error parsing dates for CreatedAt: %s", err)
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
		log.Fatalln("No tasks to process.")
	}

	averageDuration := totalDuration / time.Duration(taskCount)
	fmt.Println("Average duration:", averageDuration)
	fmt.Println("Minimum duration:", minDuration)
	fmt.Println("Maximum duration:", maxDuration)
}
