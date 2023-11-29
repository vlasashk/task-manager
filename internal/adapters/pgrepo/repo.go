package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/vlasashk/todo-manager/internal/models/todo"
)

func (db Repo) CreateTask(taskReq todo.Task) (todo.Task, error) {
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

	qry := `INSERT INTO tasks (id, title, description, due_date, status) VALUES ($1, $2, $3, $4, $5)`
	if _, err = tx.Exec(db.ctx, qry, newTask.ID, newTask.Title, newTask.Description, newTask.DueDate, newTask.Status); err != nil {
		return todo.Task{}, fmt.Errorf("exec transaction fail: %v", err)
	}

	return newTask, nil
}

func (db Repo) DeleteTaskByID(taskID string) error {
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

	qry := `UPDATE tasks SET deleted_at = NOW() WHERE id = $1`
	res, err := tx.Exec(db.ctx, qry, taskID)
	if err != nil {
		return fmt.Errorf("exec transaction fail: %v", err)
	}

	if res.RowsAffected() == 0 {
		return errors.New("invalid task id")
	}

	return nil
}

func (db Repo) GetTaskByID(taskID string) (todo.Task, error) {
	conn, err := db.DB.Acquire(db.ctx)
	if err != nil {
		return todo.Task{}, fmt.Errorf("connection acquire fail: %v", err)
	}
	defer conn.Release()

	qry := `SELECT id, title, description, due_date, status FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	row := conn.QueryRow(db.ctx, qry, taskID)

	var task todo.Task
	err = row.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return todo.Task{}, errors.New("not found")
		}
		return todo.Task{}, fmt.Errorf("query execution fail: %v", err)
	}

	return task, nil
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
