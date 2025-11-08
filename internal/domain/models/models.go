package models

import (
	"github.com/google/uuid"
	"time"
)

type Workflow struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Schedules []WorkflowSchedule `json:"schedules"`
	LogDriver string             `json:"log_driver"`
	CreatedAt int64              `json:"created_at"`
	Jobs      []Job              `json:"jobs"`
}

type WorkflowSchedule struct {
	NextRun     time.Time     `json:"next_run"`
	RepeatEvery time.Duration `json:"repeat_every"`
	Times       []int         `json:"times,omitempty"`
}

func (w *WorkflowSchedule) UpdateNext() {
	w.NextRun = w.NextRun.Add(w.RepeatEvery)
}

type Job struct {
	ID           int `json:"id"`
	WorkflowID   uuid.UUID
	Name         string        `json:"name"`
	Executor     string        `json:"executor"`
	Commands     []string      `json:"commands"`
	Dependencies []int         `json:"dependencies"`
	Retries      int           `json:"retries"`
	RetryBackoff time.Duration `json:"retry_backoff"`
	Timeout      time.Duration `json:"timeout"`
}

type JobLog struct {
	ID        uuid.UUID `json:"id"`
	LogStr    string    `json:"log"`
	Timestamp int64     `json:"timestamp"`
}

type JobResult struct {
	ID         int `json:"id"`
	WorkflowID uuid.UUID
	Logs       []JobLog `json:"logs"`
	StartedAt  int64    `json:"started_at"`
	FinishedAt int64    `json:"finished_at"`
}
