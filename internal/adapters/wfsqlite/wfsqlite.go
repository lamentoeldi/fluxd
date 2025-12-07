package wfsqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lamentoeldi/fluxd/internal/domain/models"
)

type WorkflowRepo struct {
	db *sql.DB
}

func (wr *WorkflowRepo) runTx(
	ctx context.Context,
	f func(tx *sql.Tx) error,
) (err error) {
	tx, err := wr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			txErr = fmt.Errorf("rollback/commit error: %w", txErr)
			err = errors.Join(err, txErr)
		}
	}()

	return f(tx)
}

func (wr *WorkflowRepo) addMany(
	ctx context.Context,
	workflows []models.Workflow,
) error {
	builder := sq.
		Insert("workflows").
		Columns("id", "content")
	for _, workflow := range workflows {
		data, err := json.Marshal(workflow)
		if err != nil {
			return fmt.Errorf(
				"failed to marshal workflow %s: %w",
				workflow.ID.String(),
				err,
			)
		}

		builder = builder.
			Values(workflow.ID.String(), string(data))
	}

	return wr.runTx(ctx, func(tx *sql.Tx) error {
		_, err := builder.
			RunWith(tx).
			ExecContext(ctx)
		return err
	})
}

func (wr *WorkflowRepo) replaceMany(
	ctx context.Context,
	workflows []models.Workflow,
) error {
	ids := make([]string, len(workflows))
	for i, workflow := range workflows {
		ids[i] = workflow.ID.String()
	}

	delBuilder := sq.
		Delete("workflows").
		Where(sq.Eq{"id": ids})

	insBuilder := sq.
		Insert("workflows").
		Columns("id", "content")
	for _, workflow := range workflows {
		data, err := json.Marshal(workflow)
		if err != nil {
			return fmt.Errorf(
				"failed to marshal workflow %s: %w",
				workflow.ID.String(),
				err,
			)
		}

		insBuilder = insBuilder.
			Values(workflow.ID.String(), string(data))
	}

	return wr.runTx(ctx, func(tx *sql.Tx) error {
		_, err := delBuilder.
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		_, err = insBuilder.
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func (wr *WorkflowRepo) deleteMany(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) error {
	_, err := sq.
		Delete("workflows").
		Where(sq.Eq{"id": workflowIDs}).
		RunWith(wr.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func parseWorkflows(len int, rows *sql.Rows) ([]models.Workflow, error) {
	workflows := make([]models.Workflow, 0, len)
	for rows.Next() {
		var id uuid.UUID
		var content string

		if err := rows.Scan(&id, &content); err != nil {
			return nil, err
		}

		var workflow models.Workflow
		if err := json.Unmarshal([]byte(content), &workflow); err != nil {
			return nil, err
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

func (wr *WorkflowRepo) getMany(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Workflow, error) {
	rows, err := sq.
		Select("content").
		From("workflows").
		Where(sq.Eq{"id": ids}).
		RunWith(wr.db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseWorkflows(len(ids), rows)
}

func (wr *WorkflowRepo) getAll(
	ctx context.Context,
) ([]models.Workflow, error) {
	rows, err := sq.
		Select("content").
		From("workflows").
		RunWith(wr.db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseWorkflows(0, rows)
}

func (wr *WorkflowRepo) Add(
	ctx context.Context,
	workflow models.Workflow,
) error {
	return wr.addMany(ctx, []models.Workflow{workflow})
}

func (wr *WorkflowRepo) AddMany(
	ctx context.Context,
	workflow []models.Workflow,
) error {
	return wr.addMany(ctx, workflow)
}

func (wr *WorkflowRepo) Update(
	ctx context.Context,
	upd models.WorkflowUpdate,
) error {
	workflow, err := wr.Get(ctx, upd.ID)
	if err != nil {
		return err
	}

	if upd.Name != nil {
		workflow.Name = *upd.Name
	}

	if upd.Enabled != nil {
		workflow.Enabled = *upd.Enabled
	}

	if upd.LogDriver != nil {
		workflow.LogDriver = *upd.LogDriver
	}

	if upd.Schedules != nil && len(upd.Schedules) > 0 {
		workflow.Schedules = upd.Schedules
	}

	if upd.Jobs != nil && len(upd.Jobs) > 0 {
		workflow.Jobs = upd.Jobs
	}

	data, err := json.Marshal(workflow)
	if err != nil {
		return err
	}

	_, err = sq.
		Update("workflows").
		Where(sq.Eq{"id": workflow.ID.String()}).
		Set("content", string(data)).
		RunWith(wr.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (wr *WorkflowRepo) Replace(
	ctx context.Context,
	workflow models.Workflow,
) error {
	return wr.replaceMany(ctx, []models.Workflow{workflow})
}

func (wr *WorkflowRepo) ReplaceMany(
	ctx context.Context,
	workflow []models.Workflow,
) error {
	return wr.replaceMany(ctx, workflow)
}

func (wr *WorkflowRepo) Delete(
	ctx context.Context,
	workflowID uuid.UUID,
) error {
	return wr.deleteMany(ctx, []uuid.UUID{workflowID})
}

func (wr *WorkflowRepo) DeleteMany(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) error {
	return wr.deleteMany(ctx, workflowIDs)
}

func (wr *WorkflowRepo) Get(
	ctx context.Context,
	workflowID uuid.UUID,
) (models.Workflow, error) {
	workflows, err := wr.getMany(ctx, []uuid.UUID{workflowID})
	if err != nil {
		return models.Workflow{}, err
	}

	return workflows[0], nil
}

func (wr *WorkflowRepo) GetMany(
	ctx context.Context,
	workflowIDs []uuid.UUID,
) ([]models.Workflow, error) {
	return wr.getMany(ctx, workflowIDs)
}

func (wr *WorkflowRepo) GetAll(
	ctx context.Context,
) ([]models.Workflow, error) {
	return wr.getAll(ctx)
}
