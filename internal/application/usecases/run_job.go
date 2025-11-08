package usecases

import (
	"context"
	"fmt"
	"github.com/lamentoeldi/fluxd/internal/domain/errors"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
	"sync"
)

type RunJobUseCase struct {
	core       Core
	executors  map[string]ports.JobExecutor
	executorMu sync.Mutex
	jobResults ports.JobResultBus
}

func NewRunJobUseCase(
	core Core,
	executors map[string]ports.JobExecutor,
	jobResults ports.JobResultBus,
) *RunJobUseCase {
	return &RunJobUseCase{
		core:       core,
		executors:  executors,
		jobResults: jobResults,
	}
}

func (uc *RunJobUseCase) RunJob(
	ctx context.Context,
	job models.Job,
) (models.JobResult, error) {
	ctx, cancel := uc.core.GetJobCtx(ctx, job)
	defer cancel()

	uc.executorMu.Lock()
	defer uc.executorMu.Unlock()

	executor, ok := uc.executors[job.Executor]
	if !ok {
		return models.JobResult{}, fmt.Errorf("%w: %s", errors.ErrUnknownExecutor, job.Executor)
	}

	res, err := executor.ExecuteJob(ctx, job)
	if err != nil {
		return models.JobResult{}, err
	}

	if err := uc.jobResults.Send(ctx, res); err != nil {
		return models.JobResult{}, err
	}

	return res, nil
}
