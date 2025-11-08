package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type UpdateWorkflowUseCase struct{}

func (uc *UpdateWorkflowUseCase) UpdateNext(
	ctx context.Context,
	workflow models.Workflow,
) error {
	// 1. update workflow exec date
	// 2. save workflow to storage
	panic("implement me")
}
