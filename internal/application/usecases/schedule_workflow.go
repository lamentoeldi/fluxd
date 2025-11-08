package usecases

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type ScheduleWorkflowUseCase struct {
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
	// 1. build DAG and validate workflow
	// 2. save workflow do storage
	panic("unimplemented")
}
