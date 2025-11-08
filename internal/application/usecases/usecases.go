package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type Core interface {
	BuildDAG(ctx context.Context, workflow models.Workflow) (models.DAG, error)
	UpdateDAG(ctx context.Context, dag models.DAG, result models.JobResult) (models.DAG, error)
	GetReadyJobs(ctx context.Context, dag models.DAG) ([]models.Job, error)
	UpdateWorkflowSchedules(ctx context.Context, workflow models.Workflow) error
	GetJobCtx(parent context.Context, job models.Job) (context.Context, context.CancelFunc)
}
