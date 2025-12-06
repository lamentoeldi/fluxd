package core

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/errors"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"time"
)

type Core struct {
}

func (c *Core) BuildDAG(
	ctx context.Context,
	workflow models.Workflow,
) (*models.DAG, error) {
	return buildDAG(ctx, workflow)
}

func (c *Core) UpdateDAG(
	ctx context.Context,
	dag *models.DAG,
	result models.JobResult,
) error {
	return updateDAG(ctx, dag, result)
}

func (c *Core) GetReadyJobs(
	ctx context.Context,
	dag *models.DAG,
) ([]models.Job, error) {
	ready := make([]models.Job, 0)

	walk := func(job *models.JobNode) error {
		if job.IsReady() {
			ready = append(ready, job.Job)
		}

		return nil
	}

	if err := dfsDAG(ctx, dag, walk); err != nil {
		return nil, err
	}

	return ready, nil
}

func (c *Core) UpdateWorkflowSchedules(
	ctx context.Context,
	workflow models.Workflow,
) (models.Workflow, error) {
	for _, schedule := range workflow.Schedules {
		select {
		case <-ctx.Done():
			return models.Workflow{}, ctx.Err()
		default:
			schedule.UpdateNext()
		}
	}

	return workflow, nil
}

func (c *Core) GetJobCtx(
	parent context.Context,
	job models.Job,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, job.Timeout)
}

func (c *Core) CheckDAGTimeouts(
	ctx context.Context,
	dag *models.DAG,
) error {
	walk := func(job *models.JobNode) error {
		exceedTime := job.Job.StartedAt.Add(job.Job.Timeout)
		if time.Now().After(exceedTime) && job.Status == models.StatusRunning {
			job.Status = models.StatusFailure
		}

		return nil
	}

	if err := bfsDAG(ctx, dag, walk); err != nil {
		return err
	}

	return nil
}

func (c *Core) CheckDAGRetries(
	ctx context.Context,
	dag *models.DAG,
) error {
	walk := func(job *models.JobNode) error {
		if job.Status != models.StatusFailure {
			return nil
		}

		if job.Job.Retries < 1 {
			return nil
		}

		job.Status = models.StatusPlanned
		job.Job.Retries--

		return nil
	}

	if err := bfsDAG(ctx, dag, walk); err != nil {
		return err
	}

	return nil
}

func (c *Core) CheckDAGIsCompleted(
	ctx context.Context,
	dag *models.DAG,
) (bool, error) {
	completed := true

	walk := func(job *models.JobNode) error {
		if job.Status != models.StatusSuccess {
			completed = false
			return errors.ErrStopTraversal
		}
		return nil
	}

	if err := bfsDAG(ctx, dag, walk); err != nil {
		return false, err
	}

	return completed, nil
}

func (c *Core) CheckDAGIsFailed(
	ctx context.Context,
	dag *models.DAG,
) (bool, error) {
	failed := false

	walk := func(job *models.JobNode) error {
		if job.Status == models.StatusFailure && job.Job.Retries < 1 && len(job.Dependents) > 0 {
			failed = true
		}
		return errors.ErrStopTraversal
	}

	if err := bfsDAG(ctx, dag, walk); err != nil {
		return false, err
	}

	return failed, nil
}
