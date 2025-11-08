package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type HandleJobResultUseCase struct {
	core             Core
	dagQuery         ports.DAGQuery
	jobs             ports.JobBus
	jobResultCommand ports.JobResultCommand
}

func NewHandleJobResultUseCase(
	core Core,
	dagQuery ports.DAGQuery,
	jobs ports.JobBus,
	jobResultCommand ports.JobResultCommand,
) *HandleJobResultUseCase {
	return &HandleJobResultUseCase{
		core:             core,
		dagQuery:         dagQuery,
		jobs:             jobs,
		jobResultCommand: jobResultCommand,
	}
}

func (uc *HandleJobResultUseCase) HandleJobResult(
	ctx context.Context,
	jobResult models.JobResult,
) error {
	dag, err := uc.dagQuery.GetByWorkflowID(ctx, jobResult.WorkflowID)
	if err != nil {
		return err
	}

	dag, err = uc.core.UpdateDAG(ctx, dag, jobResult)
	if err != nil {
		return err
	}

	jobs, err := uc.core.GetReadyJobs(ctx, dag)
	if err != nil {
		return err
	}

	if err := uc.jobs.SendMany(ctx, jobs); err != nil {
		return err
	}

	return uc.jobResultCommand.Add(ctx, jobResult)
}
