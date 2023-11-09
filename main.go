package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Task struct {
	StartedAt string `json:"startedAt"`
	CreatedAt string `json:"createdAt"`
	TaskArn   string `json:"taskArn"`
}

// parseTasks parses a JSON array of tasks.
func parseTasks(input []byte) ([][]Task, error) {
	var tasks [][]Task
	err := json.Unmarshal(input, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}
	return tasks, nil
}

// parseTime parses a time string in RFC3339Nano format.
func parseTime(input string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, input)
}

// calculateDurations calculates the total, minimum, and maximum durations of the tasks.
func calculateDurations(tasks [][]Task) (time.Duration, time.Duration, time.Duration, int, error) {
	var (
		totalDuration time.Duration
		minDuration   time.Duration
		maxDuration   time.Duration
		taskCount     int
		first         = true
	)

	if tasks == nil || len(tasks) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("no tasks to process")
	}
	for _, taskGroup := range tasks {
		for _, task := range taskGroup {
			if task.StartedAt == "" {
				return 0, 0, 0, 0, fmt.Errorf("task has no StartedAt time")
			}
			if task.CreatedAt == "" {
				return 0, 0, 0, 0, fmt.Errorf("task has no CreatedAt time")
			}
			started, err := parseTime(task.StartedAt)
			if err != nil {
				return 0, 0, 0, 0, fmt.Errorf("error parsing StartedAt: %w", err)
			}
			created, err := parseTime(task.CreatedAt)
			if err != nil {
				return 0, 0, 0, 0, fmt.Errorf("error parsing CreatedAt: %w", err)
			}
			if started.Before(created) {
				return 0, 0, 0, 0, fmt.Errorf("started time is before created time")
			}
			duration := started.Sub(created)
			if first || duration < minDuration {
				minDuration = duration
			}
			if duration > maxDuration {
				maxDuration = duration
			}
			totalDuration += duration
			taskCount++
			first = false
		}
	}

	return totalDuration, minDuration, maxDuration, taskCount, nil
}

func main() {
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading from stdin: %s", err)
	}

	tasks, err := parseTasks(inputData)
	if err != nil {
		log.Fatalf("Error parsing tasks: %s", err)
	}

	totalDuration, minDuration, maxDuration, taskCount, err := calculateDurations(tasks)
	if err != nil {
		log.Fatalf("Error calculating durations: %s", err)
	}

	if taskCount == 0 {
		log.Fatalln("No tasks to process.")
	}

	averageDuration := totalDuration / time.Duration(taskCount)
	fmt.Println("Average duration:", averageDuration)
	fmt.Println("Minimum duration:", minDuration)
	fmt.Println("Maximum duration:", maxDuration)
}
