package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type JobCommand interface {
	Add(ctx context.Context, job models.Job) error
	AddMany(ctx context.Context, jobs []models.Job) error
	Update(ctx context.Context, job models.Job) error
	Delete(ctx context.Context, jobID uuid.UUID) error
}

type JobQuery interface {
	Get(ctx context.Context, jobID uuid.UUID) (models.Job, error)
	GetDeps(ctx context.Context, jobID uuid.UUID) ([]models.Job, error)
	GetMany(ctx context.Context, cursor uuid.UUID, limit int) ([]models.Job, error)
	GetWithDeps(ctx context.Context, jobID, cursor uuid.UUID, limit int) ([]models.Job, error)
	GetReady(ctx context.Context, cursor uuid.UUID, limit int) ([]models.Job, error)
}

type JobLogsCommand interface {
	Add(ctx context.Context, log models.JobLog) error
	AddMany(ctx context.Context, logs []models.JobLog) error
}

type JobLogsQuery interface {
	Get(ctx context.Context, jobID uuid.UUID, limit int, before, after int64) ([]models.JobLog, error)
}

type JobExecutor interface {
	ExecuteJob(ctx context.Context, job models.Job) (models.JobResult, error)
}
