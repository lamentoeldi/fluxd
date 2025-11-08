package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type UpdateWorkflowUseCase struct {
	core            Core
	workflowCommand ports.WorkflowCommand
}

func NewUpdateWorkflowUseCase(
	core Core,
	workflowCommand ports.WorkflowCommand,
) *UpdateWorkflowUseCase {
	return &UpdateWorkflowUseCase{
		core:            core,
		workflowCommand: workflowCommand,
	}
}

func (uc *UpdateWorkflowUseCase) UpdateNext(
	ctx context.Context,
	workflow models.Workflow,
) error {
	if err := uc.core.UpdateWorkflowSchedules(ctx, workflow); err != nil {
		return err
	}

	return uc.workflowCommand.Replace(ctx, workflow)
}
