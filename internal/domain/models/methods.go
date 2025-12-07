package models

import "github.com/google/uuid"

func (w *WorkflowSchedule) UpdateNext() {
	if w.Times > 0 || w.Times == -1 {
		w.NextRun = w.NextRun.Add(w.RepeatEvery)
	}

	if w.Times > 0 {
		w.Times--
	}
}

func NewDAG(workflowID uuid.UUID) (*DAG, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	dag := &DAG{
		ID:         id,
		WorkflowID: workflowID,
		M:          make(map[int]*JobNode),
	}

	return dag, nil
}

func (n *JobNode) IsReady() bool {
	for _, dep := range n.Dependencies {
		switch dep.When {
		case CondOnFailure:
			if dep.Status != StatusFailure {
				return false
			}
		case CondOnSuccess:
			fallthrough
		default:
			if dep.Status != StatusSuccess {
				return false
			}
		}
	}

	if n.Status != StatusPlanned {
		return false
	}

	return true
}
