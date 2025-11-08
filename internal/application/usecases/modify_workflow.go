package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type ModifyWorkflowUseCase struct {
	workflowCommand ports.WorkflowCommand
}

func NewModifyWorkflowUseCase(
	workflowCommand ports.WorkflowCommand,
) *ModifyWorkflowUseCase {
	return &ModifyWorkflowUseCase{
		workflowCommand: workflowCommand,
	}
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
	return uc.workflowCommand.ReplaceMany(ctx, workflows)
}
