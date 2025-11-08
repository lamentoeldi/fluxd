package usecases

import (
	"context"
	"github.com/google/uuid"
)

type PauseWorkflowUseCase struct {
}

func (uc *PauseWorkflowUseCase) PauseWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	// 1. update workflow: Enabled = false
	panic("implement me")
}
