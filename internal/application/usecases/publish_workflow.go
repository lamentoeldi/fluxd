package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type PublishWorkflowUseCase struct{}

func (uc *PublishWorkflowUseCase) PublishWorkflow(
	ctx context.Context,
	workflow models.Workflow,
) error {
	// 1. publish workflow to message bus
	panic("implement me")
}
