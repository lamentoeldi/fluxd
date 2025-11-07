package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type ScheduleJobUseCase interface {
	ScheduleJob(ctx context.Context, job models.Job) error
	ScheduleJobs(ctx context.Context, jobs []models.Job) error
}

type UnscheduleJobUseCase interface {
	UnscheduleJob(ctx context.Context, job models.Job) error
	UnscheduleJobs(ctx context.Context, jobs []models.Job) error
}

type ModifyJobUseCase interface {
	ModifyJob(ctx context.Context, job models.Job) error
	ModifyJobs(ctx context.Context, jobs []models.Job) error
}

type ForceRunJobUseCase interface {
	ForceRunJob(ctx context.Context, jobID uuid.UUID, includeDeps bool) error
	ForceRunJobs(ctx context.Context, jobs []models.Job, includeDeps bool) error
}

type GetJobUseCase interface {
	GetJob(ctx context.Context, jobID uuid.UUID) (models.Job, error)
	GetJobs(ctx context.Context) ([]models.Job, error)
}

type GetJobLogsUseCase interface {
	GetJobLogs(ctx context.Context, jobID uuid.UUID, getAll bool) ([]models.JobLog, error)
}
