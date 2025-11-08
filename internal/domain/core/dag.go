package core

import (
	"context"
	"github.com/lamentoeldi/fluxd/internal/domain/errors"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

const (
	statusPending = "pending"
	statusPlanned = "planned"
	statusRunning = "running"
	statusSuccess = "success"
	statusFailed  = "failed"
)

type DAG map[int]*JobNode

type JobNode struct {
	Job          models.Job
	Dependencies []*JobNode
	Dependents   []*JobNode
	Status       string
}

func (n *JobNode) IsReady() bool {
	if n.Status != statusPlanned {
		return false
	}

	for _, dep := range n.Dependencies {
		if dep.Status != statusSuccess {
			return false
		}
	}

	return true
}

type DAGManager struct {
}

func NewDAG() *DAGManager {
	return &DAGManager{}
}

func (d *DAGManager) BuildDAG(
	ctx context.Context,
	workflow models.Workflow,
) (DAG, error) {
	dag := make(DAG)

	if err := d.addNodes(ctx, dag, workflow.Jobs); err != nil {
		return nil, err
	}

	if err := d.linkNodes(ctx, dag); err != nil {
		return nil, err
	}

	if err := d.validateDAG(ctx, dag); err != nil {
		return nil, err
	}

	if err := d.setStatus(ctx, dag); err != nil {
		return nil, err
	}

	return dag, nil
}

func (d *DAGManager) addNodes(
	ctx context.Context,
	m DAG,
	jobs []models.Job,
) error {
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

func (d *DAGManager) linkNodes(
	ctx context.Context,
	m DAG,
) error {
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

func (d *DAGManager) validateDAG(
	ctx context.Context,
	m DAG,
) error {
	return nil
}

func (d *DAGManager) setStatus(
	ctx context.Context,
	m DAG,
) error {
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
