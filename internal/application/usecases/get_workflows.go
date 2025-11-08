package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type GetWorkflowsUseCase struct {
	workflowQuery ports.WorkflowQuery
}

func NewGetWorkFlowsUseCase(
	workflowQuery ports.WorkflowQuery,
) *GetWorkflowsUseCase {
	return &GetWorkflowsUseCase{
		workflowQuery: workflowQuery,
	}
}

func (uc *GetWorkflowsUseCase) GetWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) (models.Workflow, error) {
	return uc.workflowQuery.Get(ctx, workflowID)
}

func (uc *GetWorkflowsUseCase) GetWorkflows(
	ctx context.Context,
) ([]models.Workflow, error) {
	return uc.workflowQuery.GetAll(ctx)
}
