package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type ResumeWorkflowUseCase struct {
	workflowCommand ports.WorkflowCommand
}

func NewResumeWorkflowUseCase(
	workflowCommand ports.WorkflowCommand,
) *ResumeWorkflowUseCase {
	return &ResumeWorkflowUseCase{
		workflowCommand: workflowCommand,
	}
}

func (uc *ResumeWorkflowUseCase) ResumeWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	enabled := true

	upd := models.WorkflowUpdate{
		ID:      workflowID,
		Enabled: &enabled,
	}

	return uc.workflowCommand.Update(ctx, upd)
}
