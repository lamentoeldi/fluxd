package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type ScheduleWorkflowUseCase interface {
	ScheduleWorkflow(ctx context.Context, workflow models.Workflow) error
	ScheduleWorkflows(ctx context.Context, workflows []models.Workflow) error
}

type UnscheduleWorkflowUseCase interface {
	UnscheduleWorkflow(ctx context.Context, workflowID uuid.UUID) error
	UnscheduleWorkflows(ctx context.Context, workflowIDs []uuid.UUID) error
}

type PauseWorkflowUseCase interface {
	PauseWorkflow(ctx context.Context, workflowID uuid.UUID) error
}

type ResumeWorkflowUseCase interface {
	ResumeWorkflow(ctx context.Context, workflowID uuid.UUID) error
}

type ModifyWorkflowUseCase interface {
	ModifyWorkflow(ctx context.Context, workflow models.Workflow) error
	ModifyWorkflows(ctx context.Context, workflows []models.Workflow) error
}

type ForceRunWorkflowUseCase interface {
	ForceRunWorkflow(ctx context.Context, workflowID uuid.UUID) error
	ForceRunWorkflows(ctx context.Context, workflowIDs []uuid.UUID) error
}

type GetWorkflowsUseCase interface {
	GetWorkflow(ctx context.Context, workflowID uuid.UUID) (models.Workflow, error)
	GetWorkflows(ctx context.Context) ([]models.Workflow, error)
}

type PublishWorkflowUseCase interface {
	PublishWorkflow(ctx context.Context, workflow models.Workflow) error
}

type UpdateWorkflowUseCase interface {
	UpdateNext(ctx context.Context, workflow models.Workflow) error
}

type ExecuteWorkFlowUseCase interface {
	ExecuteWorkflow(ctx context.Context, workflow models.Workflow) error
}

type HandleJobResultUseCase interface {
	HandleJobResult(ctx context.Context, jobResult models.JobResult) error
}

type RunJobUseCase interface {
	RunJob(ctx context.Context, job models.Job) (models.JobResult, error)
}

type CheckDAGsUseCase interface {
	CheckDAGs(ctx context.Context, maxWorkers int) error
}

type RetryJobUseCase interface {
	RetryJob(ctx context.Context, job models.Job) error
}
