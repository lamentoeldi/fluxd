package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type WorkflowCommand interface {
	Add(ctx context.Context, workflow models.Workflow) error
	AddMany(ctx context.Context, workflow []models.Workflow) error
	Update(ctx context.Context, upd models.WorkflowUpdate) error
	Replace(ctx context.Context, workflow models.Workflow) error
	ReplaceMany(ctx context.Context, workflow []models.Workflow) error
	Delete(ctx context.Context, workflowID uuid.UUID) error
	DeleteMany(ctx context.Context, workflowIDs []uuid.UUID) error
}

type WorkflowQuery interface {
	Get(ctx context.Context, workflow uuid.UUID) (models.Workflow, error)
	GetMany(ctx context.Context, workflowIDs []uuid.UUID) ([]models.Workflow, error)
	GetAll(ctx context.Context) ([]models.Workflow, error)
}

type MessageBus[T any] interface {
	Send(ctx context.Context, msg T) error
	SendMany(ctx context.Context, msgs []T) error
	Recv(ctx context.Context) (<-chan T, error)
}

type WorkflowBus interface {
	MessageBus[models.Workflow]
}

type JobBus interface {
	MessageBus[models.Job]
}

type JobResultBus interface {
	MessageBus[models.JobResult]
}

type JobResultCommand interface {
	Add(ctx context.Context, job models.JobResult) error
	AddMany(ctx context.Context, jobs []models.JobResult) error
}

type DAGCommand interface {
	Add(ctx context.Context, dag models.DAG) error
	Delete(ctx context.Context, dag models.DAG) error
}

type DAGQuery interface {
	GetByWorkflowID(ctx context.Context, workflowID uuid.UUID) (models.DAG, error)
	GetAll(ctx context.Context) ([]models.DAG, error)
}

type JobExecutor interface {
	ExecuteJob(ctx context.Context, job models.Job) (models.JobResult, error)
}
