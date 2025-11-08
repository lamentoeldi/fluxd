package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type ExecuteWorkFlowUseCase struct{}

func (uc *ExecuteWorkFlowUseCase) ExecuteWorkflow(
	ctx context.Context,
	workflow models.Workflow,
) error {
	// 1. build DAG
	// 2. get ready jobs
	// 3. publish ready jobs to bus
	panic("implement me")
}
