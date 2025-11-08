package usecases

import (
	"context"
	"github.com/google/uuid"
)

type ResumeWorkflowUseCase struct{}

func (uc *ResumeWorkflowUseCase) ResumeWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	// 1. update workflow: Enabled = true
	panic("implement me")
}
