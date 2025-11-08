package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type HandleJobResultUseCase struct{}

func (uc *HandleJobResultUseCase) HandleJobResult(
	ctx context.Context,
	jobResult models.JobResult,
) error {
	// 1. get associated DAG
	// 2. update DAG
	// 3. get ready jobs from DAG
	// 4. publish ready jobs
	// 5. save result to storage
	return nil
}
