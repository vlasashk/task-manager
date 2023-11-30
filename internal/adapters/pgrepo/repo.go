package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"github.com/vlasashk/task-manager/internal/models/tasktodo"
	"time"
)

const (
	defaultLimit   = 10
	defaultTimeout = 10 * time.Second
)

const (
	dateViolation = "23514"
)

const (
	InvalidIdErr = "invalid task id"
	DateErr      = "bad date"
)

const (
	createQry  = `INSERT INTO tasks (id, title, description, due_date, status) VALUES ($1, $2, $3, $4, $5)`
	deleteQry  = `UPDATE tasks SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	getByIDQry = `SELECT id, title, description, due_date, status 
					FROM tasks 
					WHERE id = $1 AND deleted_at IS NULL`
	updateQry = `UPDATE tasks 
					SET title = $1, description = $2, due_date = $3, status = $4
        			WHERE id = $5 AND deleted_at IS NULL`
)

func (db Repo) CreateTask(taskReq tasktodo.Request) (tasktodo.Task, error) {
	ctx, cancel := context.WithTimeout(db.ctx, defaultTimeout)
	defer cancel()
	newTask := tasktodo.New(taskReq)
	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return tasktodo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return tasktodo.Task{}, fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(ctx, tx, err)
	}()

	if _, err = tx.Exec(ctx, createQry, newTask.ID, newTask.Title, newTask.Description, newTask.DueDate, newTask.Status); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return tasktodo.Task{}, errorHandler(pgErr)
		}
		return tasktodo.Task{}, fmt.Errorf("exec transaction fail: %v", err)
	}

	return newTask, nil
}

func (db Repo) DeleteTask(taskID string) error {
	ctx, cancel := context.WithTimeout(db.ctx, defaultTimeout)
	defer cancel()
	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(ctx, tx, err)
	}()

	res, err := tx.Exec(ctx, deleteQry, taskID)
	if err != nil {
		return fmt.Errorf("exec transaction fail: %v", err)
	}

	if res.RowsAffected() == 0 {
		return errors.New(InvalidIdErr)
	}

	return nil
}

func (db Repo) GetTask(taskID string) (tasktodo.Task, error) {
	ctx, cancel := context.WithTimeout(db.ctx, defaultTimeout)
	defer cancel()
	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return tasktodo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, getByIDQry, taskID)

	var task tasktodo.Task
	var tempTime time.Time
	err = row.Scan(&task.ID, &task.Title, &task.Description, &tempTime, &task.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return tasktodo.Task{}, errors.New(InvalidIdErr)
		}
		return tasktodo.Task{}, fmt.Errorf("query execution fail: %v", err)
	}
	task.DueDate = tempTime.Format("2006-01-02")
	return task, nil
}

func (db Repo) UpdateTask(newData tasktodo.Request, taskID string) (tasktodo.Task, error) {
	updTask := tasktodo.Task{
		ID:      taskID,
		Request: newData,
	}
	ctx, cancel := context.WithTimeout(db.ctx, defaultTimeout)
	defer cancel()
	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return tasktodo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return tasktodo.Task{}, fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(ctx, tx, err)
	}()

	res, err := tx.Exec(ctx, updateQry, newData.Title, newData.Description, newData.DueDate, newData.Status, taskID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return tasktodo.Task{}, errorHandler(pgErr)
		}
		return tasktodo.Task{}, fmt.Errorf("executing update query fail: %v", err)
	}

	if res.RowsAffected() == 0 {
		return tasktodo.Task{}, errors.New(InvalidIdErr)
	}

	return updTask, nil
}

func (db Repo) ListTasks(page uint, date string, status string) ([]tasktodo.Task, error) {
	ctx, cancel := context.WithTimeout(db.ctx, defaultTimeout)
	defer cancel()
	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()
	var args []any

	qry := `SELECT id, title, description, due_date, status FROM tasks WHERE deleted_at IS NULL`
	args = []any{}

	if date != "" {
		qry += fmt.Sprintf(` AND due_date = $%d`, len(args)+1)
		args = append(args, date)
	}
	if status != "" {
		qry += fmt.Sprintf(` AND status = $%d`, len(args)+1)
		args = append(args, status)
	}

	qry += fmt.Sprintf(` ORDER BY due_date LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)

	args = append(args, defaultLimit, page*defaultLimit)

	rows, err := conn.Query(ctx, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query fail: %v", err)
	}
	defer rows.Close()

	tasks := make([]tasktodo.Task, 0, defaultLimit)

	for rows.Next() {
		var task tasktodo.Task
		var tempTime time.Time
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &tempTime, &task.Status); err != nil {
			return nil, fmt.Errorf("scanning rows fail: %v", err)
		}
		task.DueDate = tempTime.Format("2006-01-02")
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return tasks, nil
}

func txFinisher(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error().Err(err).Msg("transaction rollback for register fail")
		}
	} else {
		err = tx.Commit(ctx)
		if err != nil {
			log.Error().Err(err).Msg("transaction commit for register fail")
		}
	}
}

func errorHandler(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	case dateViolation:
		return errors.New(DateErr)
	default:
		return pgErr
	}
}
