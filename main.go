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

func parseTasks(input []byte) ([][]Task, error) {
	var tasks [][]Task
	err := json.Unmarshal(input, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}
	return tasks, nil
}

func parseTime(input string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, input)
}

func main() {
	// Read from standard input
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading from stdin: %s", err)
	}

	tasks, err := parseTasks(inputData)
	if err != nil {
		log.Fatalf("Error parsing tasks: %s", err)
	}

	var (
		totalDuration time.Duration
		minDuration   = time.Duration(math.MaxInt64)
		maxDuration   time.Duration
		taskCount     int
	)

	for _, taskGroup := range tasks {
		for _, task := range taskGroup {
			started, err := parseTime(task.StartedAt)
			if err != nil {
				log.Fatalf("Error parsing dates for StartedAt: %s", err)
			}
			created, err := parseTime(task.CreatedAt)
			if err != nil {
				log.Fatalf("Error parsing dates for CreatedAt: %s", err)
			}
			duration := started.Sub(created)

			minDuration = min(duration, minDuration)
			maxDuration = max(duration, maxDuration)
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
