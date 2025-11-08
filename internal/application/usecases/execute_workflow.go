package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/core"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type ExecuteWorkflowUseCase struct {
	core       Core
	dagCommand ports.DAGCommand
	jobsBus    ports.JobBus
}

func NewExecuteWorkflowUseCase(
	core Core,
	dagCommand ports.DAGCommand,
	jobsBus ports.JobBus,
) *ExecuteWorkflowUseCase {
	return &ExecuteWorkflowUseCase{
		core:       core,
		dagCommand: dagCommand,
		jobsBus:    jobsBus,
	}
}

func (uc *ExecuteWorkflowUseCase) ExecuteWorkflow(
	ctx context.Context,
	workflow models.Workflow,
) error {
	dag, err := core.BuildDAG(ctx, workflow)
	if err != nil {
		return err
	}

	if err := uc.dagCommand.Add(ctx, dag); err != nil {
		return err
	}

	jobs, err := uc.core.GetReadyJobs(ctx, dag)
	if err != nil {
		return err
	}

	return uc.jobsBus.SendMany(ctx, jobs)
}
