package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	StatusPending = "pending"
	StatusPlanned = "planned"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailure = "failure"
)

const (
	CondOnSuccess = "on-success"
	CondOnFailure = "on-failure"
)

type Workflow struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Schedules   []WorkflowSchedule `json:"schedules"`
	LogDriver   string             `json:"log_driver"`
	CreatedAt   int64              `json:"created_at"`
	Jobs        []Job              `json:"jobs"`
	Enabled     bool               `json:"enabled"`
}

type WorkflowSchedule struct {
	NextRun     time.Time     `json:"next_run"`
	RepeatEvery time.Duration `json:"repeat_every"`
	Times       int           `json:"times,omitempty"`
}

func (w *WorkflowSchedule) UpdateNext() {
	if w.Times > 0 || w.Times == -1 {
		w.NextRun = w.NextRun.Add(w.RepeatEvery)
	}

	if w.Times > 0 {
		w.Times--
	}
}

type WorkflowUpdate struct {
	ID        uuid.UUID
	Name      *string
	Enabled   *bool
	LogDriver *string
	Schedules []WorkflowSchedule
	Jobs      []Job
}

type Job struct {
	ID           int             `json:"id"`
	WorkflowID   uuid.UUID       `json:"workflow_id"`
	Name         string          `json:"name"`
	Executor     string          `json:"executor"`
	Commands     []string        `json:"commands"`
	Dependencies []JobDependency `json:"dependencies"`
	Retries      int             `json:"retries"`
	RetryBackoff time.Duration   `json:"retry_backoff"`
	Timeout      time.Duration   `json:"timeout"`
	StartedAt    time.Time       `json:"started_at"`
}

type JobDependency struct {
	DependencyID int    `json:"dependency_id"`
	When         string `json:"when"`
}

type JobLog struct {
	ID        uuid.UUID `json:"id"`
	LogStr    string    `json:"log"`
	Timestamp int64     `json:"timestamp"`
}

type JobResult struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	WorkflowID uuid.UUID
	Logs       []JobLog `json:"logs"`
	StartedAt  int64    `json:"started_at"`
	FinishedAt int64    `json:"finished_at"`
}

type DAG struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	M          map[int]*JobNode
}

type JobNode struct {
	Job          Job
	Dependencies []*Dependency
	Dependents   []*JobNode
	Status       string
}

type Dependency struct {
	*JobNode
	When string
}
