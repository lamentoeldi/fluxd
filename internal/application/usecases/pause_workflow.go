package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type PauseWorkflowUseCase struct {
	workflowCommand ports.WorkflowCommand
}

func NewPauseWorkflowUseCase(
	workflowCommand ports.WorkflowCommand,
) *PauseWorkflowUseCase {
	return &PauseWorkflowUseCase{
		workflowCommand: workflowCommand,
	}
}

func (uc *PauseWorkflowUseCase) PauseWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	enabled := false

	upd := models.WorkflowUpdate{
		ID:      workflowID,
		Enabled: &enabled,
	}

	return uc.workflowCommand.Update(ctx, upd)
}
