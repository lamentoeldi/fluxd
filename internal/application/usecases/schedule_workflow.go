package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"github.com/lamentoeldi/fluxd/internal/ports"
)

type ScheduleWorkflowUseCase struct {
	core            Core
	workflowCommand ports.WorkflowCommand
}

func NewScheduleWorkflowUseCase(
	core Core,
	workflowCommand ports.WorkflowCommand,
) *ScheduleWorkflowUseCase {
	return &ScheduleWorkflowUseCase{
		core:            core,
		workflowCommand: workflowCommand,
	}
}

func (uc *ScheduleWorkflowUseCase) ScheduleWorkflow(
	ctx context.Context,
	workflow models.Workflow,
) error {
	return uc.scheduleWorkflows(ctx, []models.Workflow{workflow})
}

func (uc *ScheduleWorkflowUseCase) ScheduleWorkflows(
	ctx context.Context,
	workflows []models.Workflow,
) error {
	return uc.scheduleWorkflows(ctx, workflows)
}

func (uc *ScheduleWorkflowUseCase) scheduleWorkflows(
	ctx context.Context,
	workflows []models.Workflow,
) error {
	for _, workflow := range workflows {
		_, err := uc.core.BuildDAG(ctx, workflow)
		if err != nil {
			return err
		}
	}

	return uc.workflowCommand.AddMany(ctx, workflows)
}
