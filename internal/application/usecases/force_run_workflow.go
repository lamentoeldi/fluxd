package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type ForceRunWorkflowUseCase struct {
	workflowQuery    ports.WorkflowQuery
	workflowExecutor ports.ExecuteWorkFlowUseCase
}

func NewForceRunWorkflowUseCase(
	workflowQuery ports.WorkflowQuery,
	workflowExecutor ports.ExecuteWorkFlowUseCase,
) *ForceRunWorkflowUseCase {
	return &ForceRunWorkflowUseCase{
		workflowQuery:    workflowQuery,
		workflowExecutor: workflowExecutor,
	}
}

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
	workflows, err := uc.workflowQuery.GetMany(ctx, workflowIDs)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		if err := uc.workflowExecutor.ExecuteWorkflow(ctx, workflow); err != nil {
			return err
		}
	}

	return nil
}
