package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
	"github.com/lamentoeldi/fluxd/pkg/log"
	"go.uber.org/zap"
	"sync"
)

type CheckDAGsUseCase struct {
	core       Core
	dagCommand ports.DAGCommand
	dagQuery   ports.DAGQuery
	jobs       ports.JobBus
}

func (uc *CheckDAGsUseCase) CheckDAGs(
	ctx context.Context,
	maxWorkers int,
) error {
	// 1. retrieve DAGs
	// 2. check DAGs for job timeouts
	// 3. check DAGs for possible retries
	// 4. get ready tasks and plan them
	// 5. get completed DAGs and finish them
	// 6. get blocked DAGs and mark them as failed

	// DAG DoD:
	// - all jobs are in terminal condition (success or failed with no retries left)
	// - no jobs are ready to go

	dags, err := uc.dagQuery.GetAll(ctx)
	if err != nil {
		return err
	}

	ch := make(chan models.DAG)

	go func() {
		defer close(ch)

		for _, dag := range dags {
			ch <- dag
		}
	}()

	wg := &sync.WaitGroup{}

	wg.Add(maxWorkers)
	for range maxWorkers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					log.
						FromContext(ctx).
						Info("dag checking cancelled")
					return
				case dag, ok := <-ch:
					if !ok {
						log.
							FromContext(ctx).
							Info("dag channel closed")
						return
					}

					if err := uc.checkDAG(ctx, dag); err != nil {
						log.
							FromContext(ctx).
							Error(
								"failed to check DAG",
								zap.Error(err),
							)
					}
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func (uc *CheckDAGsUseCase) checkDAG(
	ctx context.Context,
	dag models.DAG,
) error {
	if err := uc.core.CheckDAGTimeouts(ctx, dag); err != nil {
		return err
	}

	if err := uc.core.CheckDAGRetries(ctx, dag); err != nil {
		return err
	}

	if completed, err := uc.core.CheckDAGIsCompleted(ctx, dag); completed && err == nil {
		return uc.dagCommand.Delete(ctx, dag)
	}

	if failed, err := uc.core.CheckDAGIsFailed(ctx, dag); failed && err == nil {
		return uc.dagCommand.Delete(ctx, dag)
	}

	jobs, err := uc.core.GetReadyJobs(ctx, dag)
	if err != nil {
		return err
	}

	return uc.jobs.SendMany(ctx, jobs)
}
