package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTasks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      []byte
		want    [][]task
		wantErr bool
	}{
		{name: "nil check", in: nil, want: nil, wantErr: true},
		{name: "empty check", in: []byte{}, want: nil, wantErr: true},
		{name: "invalid json", in: []byte("invalid"), want: nil, wantErr: true},
		{name: "valid json", in: []byte(`[[{"startedAt":"2021-08-01T00:00:00Z","createdAt":"2021-08-01T00:00:00Z","taskArn":"arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"},{"startedAt":"2021-08-01T00:00:00Z","createdAt":"2021-08-01T00:00:00Z","taskArn":"arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}]]`), want: [][]task{{{StartedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC), TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}, {StartedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC), TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTasks(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateDurations(t *testing.T) {
	t.Parallel()
	type out struct {
		totalDuration time.Duration
		minDuration   taskLaunchDetails
		maxDuration   taskLaunchDetails
		taskCount     int
	}
	tests := []struct {
		name    string
		in      [][]task
		want    out
		wantErr bool
	}{
		{name: "nil check", in: nil, want: out{0, taskLaunchDetails{}, taskLaunchDetails{}, 0}, wantErr: true},
		{name: "empty check", in: [][]task{}, want: out{0, taskLaunchDetails{}, taskLaunchDetails{}, 0}, wantErr: true},
		{name: "started at is empty", in: [][]task{{{StartedAt: (time.Time{})}}}, want: out{0, taskLaunchDetails{}, taskLaunchDetails{}, 0}, wantErr: true},
		{name: "created at is empty", in: [][]task{{{StartedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC), CreatedAt: (time.Time{})}}}, want: out{0, taskLaunchDetails{}, taskLaunchDetails{}, 0}, wantErr: true},
		{name: "started at is before created at", in: [][]task{{{StartedAt: time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC)}}}, want: out{0, taskLaunchDetails{}, taskLaunchDetails{}, 0}, wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			totalDuration, minDuration, maxDuration, taskCount, err := calculateDurations(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateDurations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if totalDuration != tt.want.totalDuration {
				t.Errorf("calculateDurations() totalDuration = %v, want %v", totalDuration, tt.want.totalDuration)
			}
			if minDuration != tt.want.minDuration {
				t.Errorf("calculateDurations() minDuration = %v, want %v", minDuration, tt.want.minDuration)
			}
			if maxDuration != tt.want.maxDuration {
				t.Errorf("calculateDurations() maxDuration = %v, want %v", maxDuration, tt.want.maxDuration)
			}
			if taskCount != tt.want.taskCount {
				t.Errorf("calculateDurations() taskCount = %v, want %v", taskCount, tt.want.taskCount)
			}
		})
	}
}
