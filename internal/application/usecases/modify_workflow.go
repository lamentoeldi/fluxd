package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type ModifyWorkflowUseCase struct {
}

func (uc *ModifyWorkflowUseCase) ModifyWorkflow(
	ctx context.Context,
	workflow models.Workflow,
) error {
	return uc.modifyWorkflows(ctx, []models.Workflow{workflow})
}

func (uc *ModifyWorkflowUseCase) ModifyWorkflows(
	ctx context.Context,
	workflows []models.Workflow,
) error {
	return uc.modifyWorkflows(ctx, workflows)
}

func (uc *ModifyWorkflowUseCase) modifyWorkflows(
	ctx context.Context,
	workflows []models.Workflow,
) error {
	// 1. update non-empty fields of workflow
	panic("implement me")
}
