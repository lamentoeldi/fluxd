package memdag

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
	"sync"
)

type DAGRepo struct {
	m  map[uuid.UUID]*models.DAG
	mu sync.RWMutex
}

func (dr *DAGRepo) Add(
	_ context.Context,
	dag *models.DAG,
) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	if _, ok := dr.m[dag.ID]; ok {
		return fmt.Errorf("dag with id %s already exists", dag.ID)
	}

	dr.m[dag.ID] = dag
	return nil
}

func (dr *DAGRepo) Delete(
	_ context.Context,
	dagID uuid.UUID,
) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	delete(dr.m, dagID)
	return nil
}

func (dr *DAGRepo) DeleteByWorkflowID(
	_ context.Context,
	workflowID uuid.UUID,
) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	for k, v := range dr.m {
		if bytes.Equal(v.WorkflowID[:], workflowID[:]) {
			delete(dr.m, k)
		}
	}

	return nil
}

func (dr *DAGRepo) Get(
	_ context.Context,
	dagID uuid.UUID,
) (*models.DAG, error) {
	if dag, ok := dr.m[dagID]; ok {
		return dag, nil
	}

	return nil, fmt.Errorf("dag not found")
}

func (dr *DAGRepo) GetByWorkflowID(
	_ context.Context,
	workflowID uuid.UUID,
) ([]*models.DAG, error) {
	dags := make([]*models.DAG, 0)

	dr.mu.RLock()
	defer dr.mu.RUnlock()

	for _, v := range dr.m {
		if bytes.Equal(v.WorkflowID[:], workflowID[:]) {
			dags = append(dags, v)
		}
	}

	if len(dags) < 1 {
		return nil, fmt.Errorf("no dags found")
	}

	return dags, nil
}

func (dr *DAGRepo) GetAll(
	_ context.Context,
) ([]*models.DAG, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	amount := len(dr.m)
	dags := make([]*models.DAG, 0, amount)

	for _, v := range dr.m {
		dags = append(dags, v)
	}

	return dags, nil
}
