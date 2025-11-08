package core

import (
	"context"
	"fmt"
	"github.com/lamentoeldi/fluxd/internal/domain/errors"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

const (
	statusPending = "pending"
	statusPlanned = "planned"
	statusRunning = "running"
	statusSuccess = "success"
	statusFailure = "failure"
)

const (
	condOnSuccess = "on-success"
	condOnFailure = "on-failure"
)

type DAG map[int]*JobNode

type JobNode struct {
	Job          models.Job
	Dependencies []*Dependency
	Dependents   []*JobNode
	Status       string
}

type Dependency struct {
	*JobNode
	When string
}

func (n *JobNode) IsReady() bool {
	for _, dep := range n.Dependencies {
		switch dep.When {
		case condOnFailure:
			if dep.Status != statusFailure {
				return false
			}
		case condOnSuccess:
			fallthrough
		default:
			if dep.Status != statusSuccess {
				return false
			}
		}
	}

	return true
}

func BuildDAG(
	ctx context.Context,
	workflow models.Workflow,
) (DAG, error) {
	dag := make(DAG)

	if err := addNodes(ctx, dag, workflow.Jobs); err != nil {
		return nil, err
	}

	if err := linkNodes(ctx, dag); err != nil {
		return nil, err
	}

	if err := validateDAG(ctx, dag); err != nil {
		return nil, err
	}

	if err := setStatus(ctx, dag); err != nil {
		return nil, err
	}

	return dag, nil
}

func Update(
	ctx context.Context,
	dag DAG,
	jobResult models.JobResult,
) ([]*JobNode, error) {
	err := updateJobStatus(ctx, dag, jobResult)
	if err != nil {
		return nil, err
	}

	nodes, err := updateReadyJobs(ctx, dag, jobResult.ID)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func updateJobStatus(
	_ context.Context,
	dag DAG,
	jobResult models.JobResult,
) error {
	node, ok := dag[jobResult.ID]
	if !ok {
		return fmt.Errorf("node not found")
	}

	node.Status = jobResult.Status
	return nil
}

func updateReadyJobs(
	_ context.Context,
	dag DAG,
	jobID int,
) ([]*JobNode, error) {
	node, ok := dag[jobID]
	if !ok {
		return nil, fmt.Errorf("node not found")
	}

	ready := make([]*JobNode, 0)
	for _, dep := range node.Dependents {
		if dep.IsReady() && dep.Status == statusPending {
			dep.Status = statusPlanned
			ready = append(ready, dep)
		}
	}

	return ready, nil
}

func addNodes(
	ctx context.Context,
	dag DAG,
	jobs []models.Job,
) error {
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			dag[job.ID] = &JobNode{
				Job:          job,
				Dependencies: make([]*Dependency, 0, len(job.Dependencies)),
				Dependents:   make([]*JobNode, 0),
				Status:       statusPending,
			}
		}
	}

	return nil
}

func linkNodes(
	ctx context.Context,
	dag DAG,
) error {
	for _, node := range dag {
		for _, dep := range node.Job.Dependencies {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				depNode, ok := dag[dep.DependencyID]
				if !ok {
					return errors.ErrDepNotFound
				}

				node.Dependencies = append(node.Dependencies, &Dependency{
					JobNode: depNode,
					When:    dep.When,
				})
				depNode.Dependents = append(depNode.Dependents, node)
			}
		}
	}

	return nil
}

func validateDAG(
	_ context.Context,
	dag DAG,
) error {
	if hasCycle(dag) {
		return errors.ErrCycleFound
	}

	return nil
}

func setStatus(
	ctx context.Context,
	dag DAG,
) error {
	for _, node := range dag {
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

func hasCycle(d DAG) bool {
	const (
		grey  = "grey"
		black = "black"
	)

	colors := make(map[int]string)

	for id := range d {
		if colors[id] != "" {
			continue
		}

		stack := []int{id}

		for len(stack) > 0 {
			curr := stack[len(stack)-1]

			if colors[curr] == grey {
				colors[curr] = black
				stack = stack[:len(stack)-1]
				continue
			}

			if colors[curr] == black {
				stack = stack[:len(stack)-1]
				continue
			}

			colors[curr] = grey

			for _, dep := range d[curr].Dependents {
				childID := dep.Job.ID
				if colors[childID] == grey {
					return true
				}
				if colors[childID] == "" {
					stack = append(stack, childID)
				}
			}
		}
	}

	return false
}
