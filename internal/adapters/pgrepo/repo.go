package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"github.com/vlasashk/todo-manager/internal/models/todo"
	"time"
)

const (
	DateViolation = "23514"
)

const (
	InvalidIdErr = "invalid task id"
	DateErr      = "bad date"
)

const (
	createQry  = `INSERT INTO tasks (id, title, description, due_date, status) VALUES ($1, $2, $3, $4, $5)`
	deleteQry  = `UPDATE tasks SET deleted_at = NOW() WHERE id = $1`
	getByIDQry = `SELECT id, title, description, due_date, status 
					FROM tasks 
					WHERE id = $1 AND deleted_at IS NULL`
	updateQry = `UPDATE tasks 
					SET title = $1, description = $2, due_date = $3, status = $4
        			WHERE id = $5 AND deleted_at IS NULL`
)

func (db Repo) CreateTask(taskReq todo.TaskReq) (todo.Task, error) {
	newTask := todo.New(taskReq)
	conn, err := db.DB.Acquire(db.ctx)
	if err != nil {
		return todo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(db.ctx)
	if err != nil {
		return todo.Task{}, fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(db.ctx, tx, err)
	}()

	if _, err = tx.Exec(db.ctx, createQry, newTask.ID, newTask.Title, newTask.Description, newTask.DueDate, newTask.Status); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return todo.Task{}, errorHandler(pgErr)
		}
		return todo.Task{}, fmt.Errorf("exec transaction fail: %v", err)
	}

	return newTask, nil
}

func (db Repo) DeleteTask(taskID string) error {
	conn, err := db.DB.Acquire(db.ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(db.ctx)
	if err != nil {
		return fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(db.ctx, tx, err)
	}()

	res, err := tx.Exec(db.ctx, deleteQry, taskID)
	if err != nil {
		return fmt.Errorf("exec transaction fail: %v", err)
	}

	if res.RowsAffected() == 0 {
		return errors.New(InvalidIdErr)
	}

	return nil
}

func (db Repo) GetTask(taskID string) (todo.Task, error) {
	conn, err := db.DB.Acquire(db.ctx)
	if err != nil {
		return todo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(db.ctx, getByIDQry, taskID)

	var task todo.Task
	var tempTime time.Time
	err = row.Scan(&task.ID, &task.Title, &task.Description, &tempTime, &task.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return todo.Task{}, errors.New(InvalidIdErr)
		}
		return todo.Task{}, fmt.Errorf("query execution fail: %v", err)
	}
	task.DueDate = tempTime.Format("2006-01-02")
	return task, nil
}

func (db Repo) UpdateTask(newData todo.TaskReq, taskID string) (todo.Task, error) {
	updTask := todo.Task{
		ID:      taskID,
		TaskReq: newData,
	}
	conn, err := db.DB.Acquire(db.ctx)
	if err != nil {
		return todo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(db.ctx)
	if err != nil {
		return todo.Task{}, fmt.Errorf("begin transaction fail: %v", err)
	}
	defer func() {
		txFinisher(db.ctx, tx, err)
	}()

	res, err := tx.Exec(db.ctx, updateQry, newData.Title, newData.Description, newData.DueDate, newData.Status, taskID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return todo.Task{}, errorHandler(pgErr)
		}
		return todo.Task{}, fmt.Errorf("executing update query fail: %v", err)
	}

	if res.RowsAffected() == 0 {
		return todo.Task{}, errors.New(InvalidIdErr)
	}

	return updTask, nil
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
	case DateViolation:
		return errors.New(DateErr)
	default:
		return pgErr
	}
}
