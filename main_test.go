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
		want    [][]Task
		wantErr bool
	}{
		{name: "nil check", in: nil, want: nil, wantErr: true},
		{name: "empty check", in: []byte{}, want: nil, wantErr: true},
		{name: "invalid json", in: []byte("invalid"), want: nil, wantErr: true},
		{name: "valid json", in: []byte(`[[{"startedAt":"2021-08-01T00:00:00Z","createdAt":"2021-08-01T00:00:00Z","taskArn":"arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"},{"startedAt":"2021-08-01T00:00:00Z","createdAt":"2021-08-01T00:00:00Z","taskArn":"arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}]]`), want: [][]Task{{{StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}, {StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseTasks(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTasks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		want    time.Time
		wantErr bool
	}{
		{name: "empty check", in: "", want: time.Time{}, wantErr: true},
		{name: "invalid check", in: "invalid", want: time.Time{}, wantErr: true},
		{name: "valid check", in: "2021-08-01T00:00:00Z", want: time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC), wantErr: false},
		{name: "valid check with nanoseconds", in: "2021-08-01T00:00:00.123456789Z", want: time.Date(2021, 8, 1, 0, 0, 0, 123456789, time.UTC), wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("parseTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateDurations(t *testing.T) {
	t.Parallel()
	type out struct {
		totalDuration time.Duration
		minDuration   time.Duration
		maxDuration   time.Duration
		taskCount     int
	}
	tests := []struct {
		name    string
		in      [][]Task
		want    out
		wantErr bool
	}{
		{name: "nil check", in: nil, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "empty check", in: [][]Task{}, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "started at is empty", in: [][]Task{{{StartedAt: "", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "created at is empty", in: [][]Task{{{StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "started at parse error", in: [][]Task{{{StartedAt: "invalid", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}, {StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "created at parse error", in: [][]Task{{{StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "invalid", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}, {StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "started time is after created time", in: [][]Task{{{StartedAt: "2021-08-01T00:00:00Z", CreatedAt: "2021-07-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}, {StartedAt: "2021-07-01T00:00:00Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, want: out{0, 0, 0, 0}, wantErr: true},
		{name: "valid check", in: [][]Task{{{StartedAt: "2021-08-01T00:00:10Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}, {StartedAt: "2021-08-01T00:00:11Z", CreatedAt: "2021-08-01T00:00:00Z", TaskArn: "arn:aws:ecs:us-east-1:123456789012:task/12345678901234567890123456789012"}}}, want: out{21 * time.Second, 10 * time.Second, 11 * time.Second, 2}, wantErr: false},
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
