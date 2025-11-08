package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type WorkflowCommand interface {
	Add(ctx context.Context, workflow models.Workflow) error
	AddMany(ctx context.Context, workflow []models.Workflow) error
	Update(ctx context.Context, workflow models.Workflow) error
	Delete(ctx context.Context, workflowID uuid.UUID) error
}

type WorkflowQuery interface {
	Get(ctx context.Context, workflow uuid.UUID) (models.Workflow, error)
	GetAll(ctx context.Context) ([]models.Workflow, error)
}

type WorkflowBus interface {
	SendWorkflow(ctx context.Context, workflow models.Workflow) error
	SendWorkflows(ctx context.Context, workflow []models.Workflow) error
	RecvWorkflows(ctx context.Context) (<-chan models.Workflow, error)
}

type JobBus interface {
	SendJob(ctx context.Context, job models.Job) error
	SendJobs(ctx context.Context, job []models.Job) error
	RecvJobs(ctx context.Context) (<-chan models.Job, error)
}

type JobResultBus interface {
	SendJobResult(ctx context.Context, job models.JobResult) error
	SendJobResults(ctx context.Context, job []models.JobResult) error
	RecvJobResults(ctx context.Context) (<-chan models.JobResult, error)
}

type JobResultCommand interface {
	Add(ctx context.Context, job models.JobResult) error
	AddMany(ctx context.Context, jobs []models.JobResult) error
}
