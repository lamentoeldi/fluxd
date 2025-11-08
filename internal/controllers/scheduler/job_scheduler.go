package scheduler

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
	"go.uber.org/zap"
	"sync"
)

type JobSchedulerConfig struct {
	Workers int `env:"JOB_WORKERS" env-default:"1"`
}

type JobScheduler struct {
	cfg           JobSchedulerConfig
	log           *zap.Logger
	jobs          ports.JobBus
	jobsResults   ports.JobResultBus
	runJobUseCase ports.RunJobUseCase
	wg            sync.WaitGroup
}

func (j *JobScheduler) Start(
	ctx context.Context,
) error {
	ch, err := j.jobs.Recv(ctx)
	if err != nil {
		return err
	}

	j.wg.Add(j.cfg.Workers)
	for range j.cfg.Workers {
		go func() {
			defer j.wg.Done()

			j.executeJobs(ctx, ch)
		}()
	}

	return nil
}

func (j *JobScheduler) Wait() {
	j.wg.Wait()
}

func (j *JobScheduler) executeJobs(
	ctx context.Context,
	ch <-chan models.Job,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-ch:
			if !ok {
				j.log.Info("job channel closed")
				return
			}

			_, err := j.runJobUseCase.RunJob(ctx, job)
			if err != nil {
				j.log.Error(
					"job failed",
					zap.Int("job_id", job.ID),
					zap.String("job-workflow-id", job.WorkflowID.String()),
					zap.Error(err),
				)
			}
		}
	}
}
