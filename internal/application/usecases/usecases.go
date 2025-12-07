package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type Core interface {
	BuildDAG(ctx context.Context, workflow models.Workflow) (models.DAG, error)
	UpdateDAG(ctx context.Context, dag *models.DAG, result models.JobResult) error
	GetReadyJobs(ctx context.Context, dag *models.DAG) ([]models.Job, error)
	UpdateWorkflowSchedules(ctx context.Context, workflow models.Workflow) (models.Workflow, error)
	GetJobCtx(parent context.Context, job models.Job) (context.Context, context.CancelFunc)

	CheckDAGTimeouts(ctx context.Context, dag *models.DAG) error
	CheckDAGRetries(ctx context.Context, dag *models.DAG) error
	CheckDAGIsCompleted(ctx context.Context, dag *models.DAG) (bool, error)
	CheckDAGIsFailed(ctx context.Context, dag *models.DAG) (bool, error)
}
