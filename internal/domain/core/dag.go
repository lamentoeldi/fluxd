package core

import (
	"context"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/errors"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"sync"
)

const (
	statusPending = "pending"
	statusPlanned = "planned"
	statusRunning = "running"
	statusSuccess = "success"
	statusFailed  = "failed"
)

type jobGraph map[uuid.UUID]*JobNode

type JobNode struct {
	Job          models.Job
	Dependencies []*JobNode
	Dependents   []*JobNode
	Status       string
}

type DAG struct {
	jobs jobGraph
	mu   sync.RWMutex
}

func NewDAG() *DAG {
	return &DAG{
		jobs: make(jobGraph),
	}
}

func (d *DAG) Add(ctx context.Context, jobs []models.Job) error {
	m := make(jobGraph)

	if err := d.addNodes(ctx, m, jobs); err != nil {
		return err
	}

	if err := d.linkNodes(ctx, m); err != nil {
		return err
	}

	if err := d.validateDAG(ctx, m); err != nil {
		return err
	}

	if err := d.setStatus(ctx, m); err != nil {
		return err
	}

	if err := d.mergeDAG(ctx, m); err != nil {
		return err
	}

	return nil
}

func (d *DAG) addNodes(ctx context.Context, m jobGraph, jobs []models.Job) error {
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m[job.ID] = &JobNode{
				Job:          job,
				Dependencies: make([]*JobNode, 0, len(job.Dependencies)),
				Dependents:   make([]*JobNode, 0),
				Status:       statusPending,
			}
		}
	}

	return nil
}

func (d *DAG) linkNodes(ctx context.Context, m jobGraph) error {
	for _, node := range m {
		for _, depID := range node.Job.Dependencies {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				depNode, ok := m[depID]
				if !ok {
					return errors.ErrDepNotFound
				}

				node.Dependencies = append(node.Dependencies, depNode)
				depNode.Dependents = append(depNode.Dependents, node)
			}
		}
	}

	return nil
}

func (d *DAG) validateDAG(ctx context.Context, m jobGraph) error {
	return nil
}

func (d *DAG) setStatus(ctx context.Context, m jobGraph) error {
	for _, node := range m {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if len(node.Dependencies) < 1 {
				node.Status = statusPlanned
			}
		}
	}

	return nil
}

func (d *DAG) mergeDAG(ctx context.Context, m jobGraph) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.mu.Lock()
		defer d.mu.Unlock()

		for key, val := range m {
			d.jobs[key] = val
		}
	}

	return nil
}
