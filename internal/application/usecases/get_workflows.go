package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type GetWorkflowsUseCase struct {
}

func (uc *GetWorkflowsUseCase) GetWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) (models.Workflow, error) {
	// 1. get specific workflow from storage
	panic("implement me")
}

func (uc *GetWorkflowsUseCase) GetWorkflows(
	ctx context.Context,
) ([]models.Workflow, error) {
	// 1. get non-disabled workflows from storage
	panic("implement me")
}
