package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type PublishWorkflowUseCase struct {
	workflows ports.WorkflowBus
}

func NewPublishWorkflowUseCase(
	workflows ports.WorkflowBus,
) *PublishWorkflowUseCase {
	return &PublishWorkflowUseCase{
		workflows: workflows,
	}
}

func (uc *PublishWorkflowUseCase) PublishWorkflow(
	ctx context.Context,
	workflow models.Workflow,
) error {
	return uc.workflows.Send(ctx, workflow)
}
