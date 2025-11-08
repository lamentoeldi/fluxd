package scheduler

import (
	"context"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
	"go.uber.org/zap"
	"sync"
	"time"
)

type WorkflowSchedulerConfig struct {
	SchedulerBackoff time.Duration `env:"SCHEDULER_BACKOFF" env-default:"1m"`
}

func NewWorkflowSchedulerConfig() (WorkflowSchedulerConfig, error) {
	var cfg WorkflowSchedulerConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return WorkflowSchedulerConfig{}, err
	}

	return cfg, nil
}

type WorkflowScheduler struct {
	cfg             WorkflowSchedulerConfig
	log             *zap.Logger
	getWorkflows    ports.GetWorkflowsUseCase
	publishWorkflow ports.PublishWorkflowUseCase
	updateWorkflow  ports.UpdateWorkflowUseCase
	wg              sync.WaitGroup
}

func NewWorkflowScheduler(
	log *zap.Logger,
	getWorkflows ports.GetWorkflowsUseCase,
	publishWorkflow ports.PublishWorkflowUseCase,
	updateWorkflow ports.UpdateWorkflowUseCase,
) (*WorkflowScheduler, error) {
	cfg, err := NewWorkflowSchedulerConfig()
	if err != nil {
		return nil, err
	}

	return &WorkflowScheduler{
		log:             log,
		getWorkflows:    getWorkflows,
		publishWorkflow: publishWorkflow,
		updateWorkflow:  updateWorkflow,
		cfg:             cfg,
	}, nil
}

func (w *WorkflowScheduler) Start(
	ctx context.Context,
) error {
	t := time.NewTicker(w.cfg.SchedulerBackoff)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			if err := w.scheduleWorkflows(ctx); err != nil {
				w.log.Error("failed to schedule workflows", zap.Error(err))
			}
		}
	}
}

func (w *WorkflowScheduler) Wait() {
	w.wg.Wait()
}

func (w *WorkflowScheduler) scheduleWorkflows(
	ctx context.Context,
) error {
	now := time.Now()

	workflows, err := w.getWorkflows.GetWorkflows(ctx)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			go func() {
				w.wg.Add(1)
				defer w.wg.Done()

				if err := w.enqueueWorkflow(ctx, now, workflow); err != nil {
					w.log.Error(
						"failed to plan workflow",
						zap.String("workflow-id", workflow.ID.String()),
						zap.Error(err),
					)
				} else {
					w.log.Info(
						"scheduled workflow",
						zap.String("workflow-id", workflow.ID.String()),
					)
				}
			}()
		}
	}

	w.wg.Wait()

	return nil
}

func (w *WorkflowScheduler) enqueueWorkflow(
	ctx context.Context,
	now time.Time,
	workflow models.Workflow,
) error {
	for _, schedule := range workflow.Schedules {
		if schedule.NextRun.After(now.Add(-1 * time.Second)) {
			if err := w.publishWorkflow.PublishWorkflow(ctx, workflow); err != nil {
				return err
			}

			if err := w.updateWorkflow.UpdateNext(ctx, workflow); err != nil {
				return err
			}
		}
	}

	return nil
}
