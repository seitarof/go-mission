package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var t model.TODO
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t.ID = lastInsertId
	err = s.db.QueryRowContext(ctx, confirm, lastInsertId).Scan(&t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &t, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todoList := make([]*model.TODO, 0, size)

	var rows *sql.Rows
	var err error

	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var t model.TODO
		err := rows.Scan(&t.ID, &t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		todoList = append(todoList, &t)
	}

	return todoList, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	num, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if num != 0 {
		var t model.TODO
		t.ID = id
		err = s.db.QueryRowContext(ctx, confirm, id).Scan(&t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &t, nil
	} else {
		err = &model.ErrNotFound{}
		log.Println(err)
		return nil, err
	}
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) != 0 {
		ids2 := make([]interface{}, 0)
		for _, id := range ids {
			ids2 = append(ids2, id)
		}

		result, err := s.db.ExecContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids2)-1)), ids2...)

		if err != nil {
			log.Println(err)
			return err
		}

		num, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
			return err
		}
		if num == 0 {
			err = &model.ErrNotFound{}
			log.Println(err)
			return err
		}
	}
	return nil
}
