package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

func makeJob(
	id int,
	deps []models.JobDependency,
) models.Job {
	return models.Job{
		ID:           id,
		WorkflowID:   uuid.New(),
		Name:         "job",
		Executor:     "echo",
		Commands:     []string{"echo hello"},
		Dependencies: deps,
		Retries:      0,
		RetryBackoff: 0,
		Timeout:      0,
	}
}

func (d DAG) debugPrint() {
	for _, node := range d {
		fmt.Println("job: ", node.Job)
		for _, dep := range node.Dependencies {
			fmt.Println("dependency: ", dep.Job.ID, dep.When)
		}
		fmt.Println(node.Status)
		fmt.Println("is ready: ", node.IsReady())
		for _, dep := range node.Dependents {
			fmt.Println("dependent: ", dep.Job.ID)
		}
	}
}

func TestBuildDAGBasic(t *testing.T) {
	// 1 -> 2
	// 1 -> 3
	// 2, 3 -> 4
	jobs := []models.Job{
		makeJob(1, nil),
		makeJob(2, []models.JobDependency{
			{DependencyID: 1, When: condOnSuccess},
		}),
		makeJob(3, []models.JobDependency{
			{DependencyID: 1, When: condOnSuccess},
		}),
		makeJob(4, []models.JobDependency{
			{DependencyID: 2, When: condOnSuccess},
			{DependencyID: 3, When: condOnSuccess},
		}),
	}

	d := NewDAG()
	dag, err := d.BuildDAG(context.Background(), models.Workflow{
		ID:   uuid.New(),
		Jobs: jobs,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 4; i++ {
		if _, ok := dag[i]; !ok {
			t.Fatalf("node %d not found", i)
		}
	}

	if len(dag[4].Dependencies) != 2 {
		t.Errorf("expected 2 dependencies for node 4, got %d", len(dag[4].Dependencies))
	}

	if dag[1].Status != statusPlanned {
		t.Errorf("expected node 1 to be planned, got %s", dag[1].Status)
	}

	dag.debugPrint()
}

func TestBuildDAGWithConditions(t *testing.T) {
	// 1 -> 2
	// 1 [on failure] -> 3
	jobs := []models.Job{
		makeJob(1, nil),
		makeJob(2, []models.JobDependency{
			{DependencyID: 1, When: condOnSuccess},
		}),
		makeJob(3, []models.JobDependency{
			{DependencyID: 1, When: condOnFailure},
		}),
	}

	d := NewDAG()
	dag, err := d.BuildDAG(context.Background(), models.Workflow{
		ID:   uuid.New(),
		Jobs: jobs,
	})
	if err != nil {
		t.Fatal(err)
	}

	node1 := dag[1]
	node2 := dag[2]
	node3 := dag[3]

	found := false
	for _, dep := range node2.Dependencies {
		if dep.JobNode.Job.ID == node1.Job.ID && dep.When == condOnSuccess {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected node2 to depend on node1 with condOnSuccess")
	}

	found = false
	for _, dep := range node3.Dependencies {
		if dep.JobNode.Job.ID == node1.Job.ID && dep.When == condOnFailure {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected node3 to depend on node1 with condOnFailure")
	}

	node1.Status = statusSuccess
	if node2.IsReady() != true {
		t.Errorf("node2 should be ready when dependency succeeded and condition = condOnSuccess")
	}
	if node3.IsReady() {
		t.Errorf("node3 should not be ready when dependency succeeded but condition = condOnFailure")
	}

	node1.Status = statusFailure
	if !node3.IsReady() {
		t.Errorf("node3 should be ready when dependency failed and condition = condOnFailure")
	}

	dag.debugPrint()
}

func TestHasCycle(t *testing.T) {
	cases := []struct {
		name     string
		jobs     []models.Job
		expected bool
	}{
		{
			name: "no cycle",
			// 1 -> 2
			// 1 -> 3
			// 2, 3 -> 4
			jobs: []models.Job{
				makeJob(1, nil),
				makeJob(2, []models.JobDependency{
					{DependencyID: 1, When: condOnSuccess},
				}),
				makeJob(3, []models.JobDependency{
					{DependencyID: 1, When: condOnSuccess},
				}),
				makeJob(4, []models.JobDependency{
					{DependencyID: 2, When: condOnSuccess},
					{DependencyID: 3, When: condOnSuccess},
				}),
			},
			expected: false,
		},
		{
			name: "simple cycle",
			// 1 -> 2
			// 2 -> 1
			jobs: []models.Job{
				makeJob(1, []models.JobDependency{
					{DependencyID: 2, When: condOnSuccess},
				}),
				makeJob(2, []models.JobDependency{
					{DependencyID: 1, When: condOnSuccess},
				}),
			},
			expected: true,
		},
		{
			name: "medium cycle",
			// 1 -> 2
			// 1 -> 3
			// 2 -> 3
			// 2 -> 4
			// 2 -> 5
			// 3 -> 5
			// 5 -> 3
			jobs: []models.Job{
				makeJob(1, nil),
				makeJob(2, []models.JobDependency{
					{DependencyID: 1, When: condOnSuccess},
				}),
				makeJob(3, []models.JobDependency{
					{DependencyID: 1, When: condOnSuccess},
					{DependencyID: 2, When: condOnSuccess},
					{DependencyID: 5, When: condOnSuccess},
				}),
				makeJob(4, []models.JobDependency{
					{DependencyID: 2, When: condOnSuccess},
				}),
				makeJob(5, []models.JobDependency{
					{DependencyID: 2, When: condOnSuccess},
					{DependencyID: 3, When: condOnSuccess},
				}),
			},
			expected: true,
		},
	}

	d := NewDAG()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			dag := make(DAG)
			err := d.addNodes(ctx, dag, tc.jobs)
			err = d.linkNodes(ctx, dag)
			if err != nil {
				t.Fatal(err)
			}

			cyclic := hasCycle(dag)
			if tc.expected != cyclic {
				t.Errorf("expected %v, got %v", tc.expected, cyclic)
			}
		})
	}
}
