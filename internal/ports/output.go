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
