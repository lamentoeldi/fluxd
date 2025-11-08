package scheduler

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
	"go.uber.org/zap"
)

type DAGScheduler struct {
	log               *zap.Logger
	jobs              ports.JobBus
	jobResults        ports.JobResultBus
	jobResultsHandler ports.HandleJobResultUseCase
	workflows         ports.WorkflowBus
	workflowExecutor  ports.ExecuteWorkFlowUseCase
}

func NewDAGScheduler(
	log *zap.Logger,
	jobs ports.JobBus,
	jobResults ports.JobResultBus,
	jobResultsHandler ports.HandleJobResultUseCase,
	workflows ports.WorkflowBus,
	workflowExecutor ports.ExecuteWorkFlowUseCase,
) *DAGScheduler {
	return &DAGScheduler{
		log:               log,
		jobs:              jobs,
		jobResults:        jobResults,
		jobResultsHandler: jobResultsHandler,
		workflows:         workflows,
		workflowExecutor:  workflowExecutor,
	}
}

func (d *DAGScheduler) Start(ctx context.Context) error {
	workflows, err := d.workflows.Recv(ctx)
	if err != nil {
		return err
	}

	jobResults, err := d.jobResults.Recv(ctx)
	if err != nil {
		return err
	}

	go d.startWorkflows(ctx, workflows)
	go d.startJobResults(ctx, jobResults)

	return nil
}

func (d *DAGScheduler) startWorkflows(
	ctx context.Context,
	ch <-chan models.Workflow,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case workflow, ok := <-ch:
			if !ok {
				d.log.Info("workflow channel closed")
				return
			}

			if err := d.workflowExecutor.ExecuteWorkflow(ctx, workflow); err != nil {
				d.log.Error(
					"failed to execute workflow",
					zap.String("workflow-id", workflow.ID.String()),
					zap.Error(err),
				)
			}
		}
	}
}

func (d *DAGScheduler) startJobResults(
	ctx context.Context,
	ch <-chan models.JobResult,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case jobResult, ok := <-ch:
			if !ok {
				d.log.Info("job result channel closed")
				return
			}

			if err := d.jobResultsHandler.HandleJobResult(ctx, jobResult); err != nil {
				d.log.Error(
					"failed to handle job result",
					zap.Int("job-id", jobResult.ID),
					zap.String("job-workflow-id", jobResult.WorkflowID.String()),
					zap.Error(err),
				)
			}
		}
	}
}
