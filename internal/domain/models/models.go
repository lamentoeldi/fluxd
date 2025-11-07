package models

import (
	"github.com/google/uuid"
	"time"
)

type Job struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Schedules    []time.Time   `json:"schedules"`
	Executor     string        `json:"executor"`
	Commands     []string      `json:"commands"`
	Dependencies []uuid.UUID   `json:"dependencies"`
	Retries      int           `json:"retries"`
	RetryBackoff time.Duration `json:"retry_backoff"`
	Timeout      time.Duration `json:"timeout"`
	LogDriver    string        `json:"log_driver"`
	CreatedAt    int64         `json:"created_at"`
}

type JobLog struct {
	ID        uuid.UUID `json:"id"`
	LogStr    string    `json:"log"`
	Timestamp int64     `json:"timestamp"`
}

type JobResult struct {
	ID         uuid.UUID `json:"id"`
	Logs       []JobLog  `json:"logs"`
	StartedAt  int64     `json:"started_at"`
	FinishedAt int64     `json:"finished_at"`
}
