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
