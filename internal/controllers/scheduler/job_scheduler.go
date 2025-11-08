package scheduler

import (
	"github.com/lamentoeldi/fluxd/internal/ports"
	"go.uber.org/zap"
)

type JobScheduler struct {
	log         *zap.Logger
	jobs        ports.JobBus
	jobsResults ports.JobResultBus
}

func (j *JobScheduler) Start() error {}
