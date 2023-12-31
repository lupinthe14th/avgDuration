package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type task struct {
	StartedAt time.Time `json:"startedAt"`
	CreatedAt time.Time `json:"createdAt"`
	TaskArn   string    `json:"taskArn"`
}

type taskLaunchDetails struct {
	Duration time.Duration
	TaskArn  string
}

// parseTasks parses a JSON array of tasks.
func parseTasks(input []byte) ([][]task, error) {
	var tasks [][]task
	err := json.Unmarshal(input, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}
	return tasks, nil
}

// calculateDurations calculates the total, minimum, and maximum durations of the tasks.
func calculateDurations(tasks [][]task) (time.Duration, taskLaunchDetails, taskLaunchDetails, int, error) {
	var (
		totalDuration time.Duration
		minDuration   taskLaunchDetails
		maxDuration   taskLaunchDetails
		taskCount     int
		first         = true
	)

	if len(tasks) == 0 {
		return 0, taskLaunchDetails{}, taskLaunchDetails{}, 0, fmt.Errorf("no tasks to process")
	}
	for _, taskGroup := range tasks {
		for _, task := range taskGroup {
			if task.StartedAt == (time.Time{}) {
				return 0, taskLaunchDetails{}, taskLaunchDetails{}, 0, fmt.Errorf("task has no StartedAt time")
			}
			if task.CreatedAt == (time.Time{}) {
				return 0, taskLaunchDetails{}, taskLaunchDetails{}, 0, fmt.Errorf("task has no CreatedAt time")
			}
			started := task.StartedAt
			created := task.CreatedAt
			if started.Before(created) {
				return 0, taskLaunchDetails{}, taskLaunchDetails{}, 0, fmt.Errorf("started time is before created time")
			}
			duration := started.Sub(created)
			if first || duration < minDuration.Duration {
				minDuration = taskLaunchDetails{
					Duration: duration,
					TaskArn:  task.TaskArn,
				}
			}
			if duration > maxDuration.Duration {
				maxDuration = taskLaunchDetails{
					Duration: duration,
					TaskArn:  task.TaskArn,
				}
			}
			totalDuration += duration
			taskCount++
			first = false
		}
	}

	return totalDuration, minDuration, maxDuration, taskCount, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		logger.Error("Error reading from stdin: %s", err)
		os.Exit(1)
	}

	tasks, err := parseTasks(inputData)
	if err != nil {
		logger.Error("Error parsing tasks: %s", err)
		os.Exit(1)
	}

	totalDuration, minDuration, maxDuration, taskCount, err := calculateDurations(tasks)
	if err != nil {
		logger.Error("Error calculating durations: %s", err)
		os.Exit(1)
	}

	if taskCount == 0 {
		logger.Error("No tasks to process.")
		os.Exit(1)
	}

	averageDuration := totalDuration / time.Duration(taskCount)
	fmt.Println("Average duration:", averageDuration)
	fmt.Println("Minimum duration: ", minDuration.Duration)
	fmt.Println("Maximum duration: ", maxDuration.Duration)
	fmt.Println("Minimum task ARN: ", minDuration.TaskArn)
	fmt.Println("Maximum task ARN: ", maxDuration.TaskArn)
}
