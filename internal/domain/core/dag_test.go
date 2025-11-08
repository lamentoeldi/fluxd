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

func debugPrint(dag models.DAG) {
	for _, node := range dag {
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
			{DependencyID: 1, When: models.CondOnSuccess},
		}),
		makeJob(3, []models.JobDependency{
			{DependencyID: 1, When: models.CondOnSuccess},
		}),
		makeJob(4, []models.JobDependency{
			{DependencyID: 2, When: models.CondOnSuccess},
			{DependencyID: 3, When: models.CondOnSuccess},
		}),
	}

	dag, err := BuildDAG(context.Background(), models.Workflow{
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

	if dag[1].Status != models.StatusPlanned {
		t.Errorf("expected node 1 to be planned, got %s", dag[1].Status)
	}

	debugPrint(dag)
}

func TestBuildDAGWithConditions(t *testing.T) {
	// 1 -> 2
	// 1 [on failure] -> 3
	jobs := []models.Job{
		makeJob(1, nil),
		makeJob(2, []models.JobDependency{
			{DependencyID: 1, When: models.CondOnSuccess},
		}),
		makeJob(3, []models.JobDependency{
			{DependencyID: 1, When: models.CondOnFailure},
		}),
	}

	dag, err := BuildDAG(context.Background(), models.Workflow{
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
		if dep.JobNode.Job.ID == node1.Job.ID && dep.When == models.CondOnSuccess {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected node2 to depend on node1 with condOnSuccess")
	}

	found = false
	for _, dep := range node3.Dependencies {
		if dep.JobNode.Job.ID == node1.Job.ID && dep.When == models.CondOnFailure {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected node3 to depend on node1 with condOnFailure")
	}

	node1.Status = models.StatusSuccess
	if node2.IsReady() != true {
		t.Errorf("node2 should be ready when dependency succeeded and condition = condOnSuccess")
	}
	if node3.IsReady() {
		t.Errorf("node3 should not be ready when dependency succeeded but condition = condOnFailure")
	}

	node1.Status = models.StatusFailure
	if !node3.IsReady() {
		t.Errorf("node3 should be ready when dependency failed and condition = condOnFailure")
	}

	debugPrint(dag)
}

func TestUpdateJobStatusAndReady(t *testing.T) {
	// 1 -> 2 -> 3
	jobs := []models.Job{
		makeJob(1, nil),
		makeJob(2, []models.JobDependency{
			{DependencyID: 1, When: models.CondOnSuccess},
		}),
		makeJob(3, []models.JobDependency{
			{DependencyID: 2, When: models.CondOnSuccess},
		}),
	}

	dag, err := BuildDAG(t.Context(), models.Workflow{
		ID:   uuid.UUID{},
		Jobs: jobs,
	})
	if err != nil {
		t.Fatal(err)
	}

	jobResult := models.JobResult{
		ID:     1,
		Status: models.StatusSuccess,
	}

	ready, err := Update(context.Background(), dag, jobResult)
	if err != nil {
		t.Fatal(err)
	}

	if dag[2].Status != models.StatusPlanned {
		t.Errorf("expected job 2 status = planned, got %s", dag[2].Status)
	}

	if dag[3].Status != models.StatusPending {
		t.Errorf("expected job 3 status = pending, got %s", dag[3].Status)
	}

	if len(ready) != 1 || ready[0].Job.ID != 2 {
		t.Errorf("expected ready nodes = [2], got %+v", ready)
	}
}

func TestUpdateJobStatusWithConditionalDeps_IndependentCases(t *testing.T) {
	cases := []struct {
		name       string
		jobs       []models.Job
		jobResult  models.JobResult
		expPlanned []int
		expPending []int
	}{
		{
			// 1 -> 2
			name: "job1 success triggers job2 [on-success]",
			jobs: []models.Job{
				makeJob(1, nil),
				makeJob(2, []models.JobDependency{{DependencyID: 1, When: models.CondOnSuccess}}),
			},
			jobResult:  models.JobResult{ID: 1, Status: models.StatusSuccess},
			expPlanned: []int{2},
			expPending: []int{},
		},
		{
			// 1 [on failure] -> 2
			name: "job1 failure triggers job2 [on-failure]",
			jobs: []models.Job{
				makeJob(1, nil),
				makeJob(2, []models.JobDependency{{DependencyID: 1, When: models.CondOnFailure}}),
			},
			jobResult:  models.JobResult{ID: 1, Status: models.StatusFailure},
			expPlanned: []int{2},
			expPending: []int{},
		},
		{
			// 1 [on failure] -> 2
			name: "job1 success does not trigger job2 [on-failure]",
			jobs: []models.Job{
				makeJob(1, nil),
				makeJob(2, []models.JobDependency{{DependencyID: 1, When: models.CondOnFailure}}),
			},
			jobResult:  models.JobResult{ID: 1, Status: models.StatusSuccess},
			expPlanned: []int{},
			expPending: []int{2},
		},
		{
			// 1 -> 2
			name: "job1 failure does not trigger job2 [on-success]",
			jobs: []models.Job{
				makeJob(1, nil),
				makeJob(2, []models.JobDependency{{DependencyID: 1, When: models.CondOnSuccess}}),
			},
			jobResult:  models.JobResult{ID: 1, Status: models.StatusFailure},
			expPlanned: []int{},
			expPending: []int{2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			workflow := models.Workflow{
				ID:   uuid.New(),
				Jobs: c.jobs,
			}

			dag, err := BuildDAG(context.Background(), workflow)
			if err != nil {
				t.Fatal(err)
			}

			ready, err := Update(context.Background(), dag, c.jobResult)
			if err != nil {
				t.Fatal(err)
			}

			for _, id := range c.expPlanned {
				if dag[id].Status != models.StatusPlanned {
					t.Errorf("expected job %d status = planned, got %s", id, dag[id].Status)
				}
			}
			for _, id := range c.expPending {
				if dag[id].Status != models.StatusPending {
					t.Errorf("expected job %d status = pending, got %s", id, dag[id].Status)
				}
			}

			expReadyIDs := make(map[int]struct{})
			for _, id := range c.expPlanned {
				expReadyIDs[id] = struct{}{}
			}
			if len(ready) != len(expReadyIDs) {
				t.Errorf("expected %d ready nodes, got %d", len(expReadyIDs), len(ready))
			}
			for _, n := range ready {
				if _, ok := expReadyIDs[n.Job.ID]; !ok {
					t.Errorf("unexpected ready node %d", n.Job.ID)
				}
			}
		})
	}
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
					{DependencyID: 1, When: models.CondOnSuccess},
				}),
				makeJob(3, []models.JobDependency{
					{DependencyID: 1, When: models.CondOnSuccess},
				}),
				makeJob(4, []models.JobDependency{
					{DependencyID: 2, When: models.CondOnSuccess},
					{DependencyID: 3, When: models.CondOnSuccess},
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
					{DependencyID: 2, When: models.CondOnSuccess},
				}),
				makeJob(2, []models.JobDependency{
					{DependencyID: 1, When: models.CondOnSuccess},
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
					{DependencyID: 1, When: models.CondOnSuccess},
				}),
				makeJob(3, []models.JobDependency{
					{DependencyID: 1, When: models.CondOnSuccess},
					{DependencyID: 2, When: models.CondOnSuccess},
					{DependencyID: 5, When: models.CondOnSuccess},
				}),
				makeJob(4, []models.JobDependency{
					{DependencyID: 2, When: models.CondOnSuccess},
				}),
				makeJob(5, []models.JobDependency{
					{DependencyID: 2, When: models.CondOnSuccess},
					{DependencyID: 3, When: models.CondOnSuccess},
				}),
			},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			dag := make(models.DAG)
			err := addNodes(ctx, dag, tc.jobs)
			err = linkNodes(ctx, dag)
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
