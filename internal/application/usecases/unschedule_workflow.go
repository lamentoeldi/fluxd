package usecases

import (
	"context"
	"github.com/google/uuid"
)

type UnscheduleWorkflowUseCase struct {
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
	// 1. remove workflows from storage
	panic("implement me")
}
