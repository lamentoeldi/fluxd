package core

import (
	"context"
	stderr "errors"
	"fmt"
	"github.com/lamentoeldi/fluxd/internal/domain/errors"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"sync"
)

func buildDAG(
	ctx context.Context,
	workflow models.Workflow,
) (models.DAG, error) {
	dag := make(models.DAG)

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

func updateDAG(
	ctx context.Context,
	dag models.DAG,
	jobResult models.JobResult,
) error {
	err := updateJobStatus(ctx, dag, jobResult)
	if err != nil {
		return err
	}

	_, err = updateReadyJobs(ctx, dag, jobResult.ID)
	if err != nil {
		return err
	}

	return nil
}

func updateJobStatus(
	_ context.Context,
	dag models.DAG,
	jobResult models.JobResult,
) error {
	node, ok := dag[jobResult.ID]
	if !ok {
		return errors.ErrNodeNotFound
	}

	if node.Status != models.StatusRunning {
		return errors.ErrInvalidStatus
	}

	node.Status = jobResult.Status
	return nil
}

func updateReadyJobs(
	_ context.Context,
	dag models.DAG,
	jobID int,
) ([]*models.JobNode, error) {
	node, ok := dag[jobID]
	if !ok {
		return nil, fmt.Errorf("node not found")
	}

	ready := make([]*models.JobNode, 0)
	for _, dep := range node.Dependents {
		if dep.IsReady() && dep.Status == models.StatusPending {
			dep.Status = models.StatusPlanned
			ready = append(ready, dep)
		}
	}

	return ready, nil
}

func addNodes(
	ctx context.Context,
	dag models.DAG,
	jobs []models.Job,
) error {
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			dag[job.ID] = &models.JobNode{
				Job:          job,
				Dependencies: make([]*models.Dependency, 0, len(job.Dependencies)),
				Dependents:   make([]*models.JobNode, 0),
				Status:       models.StatusPending,
			}
		}
	}

	return nil
}

func linkNodes(
	ctx context.Context,
	dag models.DAG,
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

				node.Dependencies = append(node.Dependencies, &models.Dependency{
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
	dag models.DAG,
) error {
	if hasCycle(dag) {
		return errors.ErrCycleFound
	}

	return nil
}

func setStatus(
	ctx context.Context,
	dag models.DAG,
) error {
	for _, node := range dag {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if len(node.Dependencies) < 1 {
				node.Status = models.StatusPlanned
			}
		}
	}

	return nil
}

func hasCycle(
	d models.DAG,
) bool {
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

type walkFunc func(node *models.JobNode) error

type stack[T any] struct {
	stack []T
	mu    sync.RWMutex
}

func (s *stack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.stack)
}

func (s *stack[T]) Push(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stack = append(s.stack, v)
}

func (s *stack[T]) Pop() T {
	if s.Len() < 1 {
		var zero T
		return zero
	}

	lastIdx := s.Len() - 1

	s.mu.Lock()
	defer s.mu.Unlock()

	last := s.stack[lastIdx]

	s.stack = s.stack[:lastIdx]
	return last
}

type queue[T any] struct {
	queue []T
	mu    sync.RWMutex
}

func (q *queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.queue)
}

func (q *queue[T]) Push(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, v)
}

func (q *queue[T]) Pop() T {
	if q.Len() < 1 {
		var zero T
		return zero
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	first := q.queue[0]

	q.queue = q.queue[1:]
	return first
}

func dfsDAG(
	ctx context.Context,
	dag models.DAG,
	walk walkFunc,
) error {
	s := stack[int]{}
	visited := make(map[int]struct{})

	for id := range dag {
		if _, ok := visited[id]; ok {
			continue
		}

		s.Push(id)

		for s.Len() > 0 {
			curr := s.Pop()

			if _, ok := visited[curr]; ok {
				continue
			}

			visited[curr] = struct{}{}

			node, ok := dag[curr]
			if !ok {
				continue
			}

			err := walk(node)
			if stderr.Is(err, errors.ErrStopTraversal) {
				return nil
			}

			if err != nil {
				return err
			}

			for _, dep := range node.Dependents {
				if _, ok := visited[dep.Job.ID]; !ok {
					s.Push(dep.Job.ID)
				}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

func bfsDAG(
	ctx context.Context,
	dag models.DAG,
	walk walkFunc,
) error {
	q := queue[int]{}
	visited := make(map[int]struct{})

	for id := range dag {
		if _, ok := visited[id]; ok {
			continue
		}

		q.Push(id)

		for q.Len() > 0 {
			curr := q.Pop()

			if _, ok := visited[curr]; ok {
				continue
			}

			visited[curr] = struct{}{}

			node, ok := dag[curr]
			if !ok {
				continue
			}

			err := walk(node)
			if stderr.Is(err, errors.ErrStopTraversal) {
				return nil
			}

			if err != nil {
				return err
			}

			for _, dep := range node.Dependents {
				if _, ok := visited[dep.Job.ID]; !ok {
					q.Push(dep.Job.ID)
				}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}
