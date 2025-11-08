package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type UnscheduleWorkflowUseCase struct {
	workflowCommand ports.WorkflowCommand
}

func NewUnscheduleWorkflowUseCase(
	workflowCommand ports.WorkflowCommand,
) *UnscheduleWorkflowUseCase {
	return &UnscheduleWorkflowUseCase{
		workflowCommand: workflowCommand,
	}
}

func (uc *UnscheduleWorkflowUseCase) UnscheduleWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	return uc.unscheduleWorkflows(ctx, []uuid.UUID{workflowID})
}

func (uc *UnscheduleWorkflowUseCase) UnscheduleWorkflows(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) error {
	return uc.unscheduleWorkflows(ctx, workflowIDs)
}

func (uc *UnscheduleWorkflowUseCase) unscheduleWorkflows(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) error {
	return uc.workflowCommand.DeleteMany(ctx, workflowIDs)
}
