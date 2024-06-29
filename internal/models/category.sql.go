// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: category.sql

package models

import (
	"context"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO category (name) VALUES ($1) RETURNING id, name
`

func (q *Queries) CreateCategory(ctx context.Context, name string) (Category, error) {
	row := q.db.QueryRow(ctx, createCategory, name)
	var i Category
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteCategory, id)
	return err
}

const getCategoryByFuzzy = `-- name: GetCategoryByFuzzy :one
SELECT id, name from category WHERE name ILIKE $1
`

func (q *Queries) GetCategoryByFuzzy(ctx context.Context, name string) (Category, error) {
	row := q.db.QueryRow(ctx, getCategoryByFuzzy, name)
	var i Category
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getCategoryByName = `-- name: GetCategoryByName :one
SELECT id, name from category WHERE name = $1
`

func (q *Queries) GetCategoryByName(ctx context.Context, name string) (Category, error) {
	row := q.db.QueryRow(ctx, getCategoryByName, name)
	var i Category
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listCategory = `-- name: ListCategory :many
SELECT id, name from category
`

func (q *Queries) ListCategory(ctx context.Context) ([]Category, error) {
	rows, err := q.db.Query(ctx, listCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
