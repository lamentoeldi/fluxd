package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type RunJobUseCase struct{}

func (uc *RunJobUseCase) RunJob(
	ctx context.Context,
	job models.Job,
) (models.JobResult, error) {
	// 1. choose specified job executor
	// 2. run in executor considering timeouts and retries
	// 3. return result
	panic("implement me")
}
