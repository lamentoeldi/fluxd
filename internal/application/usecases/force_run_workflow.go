package usecases

import (
	"context"
	"github.com/google/uuid"
)

type ForceRunWorkflowUseCase struct{}

func (uc *ForceRunWorkflowUseCase) ForceRunWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	return uc.forceRunWorkflows(ctx, []uuid.UUID{workflowID})
}

func (uc *ForceRunWorkflowUseCase) ForceRunWorkflows(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) error {
	return uc.forceRunWorkflows(ctx, workflowIDs)
}

func (uc *ForceRunWorkflowUseCase) forceRunWorkflows(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) error {
	// 1. retrieve workflows from storage
	// 2. execute workflows
	panic("implement me")
}
