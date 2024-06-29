// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: item.sql

package models

import (
	"context"
)

const createItem = `-- name: CreateItem :one
INSERT INTO item (name, description, rarity, cost) VALUES ($1, $2, $3, $4) RETURNING id, name, description, rarity, cost
`

type CreateItemParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rarity      Rarity `json:"rarity"`
	Cost        int32  `json:"cost"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, createItem,
		arg.Name,
		arg.Description,
		arg.Rarity,
		arg.Cost,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rarity,
		&i.Cost,
	)
	return i, err
}

const deleteItem = `-- name: DeleteItem :exec
DELETE FROM item WHERE id = $1
`

func (q *Queries) DeleteItem(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteItem, id)
	return err
}

const getItem = `-- name: GetItem :one
SELECT id, name, description, rarity, cost from item WHERE id = $1
`

func (q *Queries) GetItem(ctx context.Context, id int32) (Item, error) {
	row := q.db.QueryRow(ctx, getItem, id)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rarity,
		&i.Cost,
	)
	return i, err
}

const getItemByFuzzy = `-- name: GetItemByFuzzy :one
SELECT id, name, description, rarity, cost from item WHERE name ILIKE $1
`

func (q *Queries) GetItemByFuzzy(ctx context.Context, name string) (Item, error) {
	row := q.db.QueryRow(ctx, getItemByFuzzy, name)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rarity,
		&i.Cost,
	)
	return i, err
}

const getItemByName = `-- name: GetItemByName :one
SELECT id, name, description, rarity, cost from item WHERE name = $1
`

func (q *Queries) GetItemByName(ctx context.Context, name string) (Item, error) {
	row := q.db.QueryRow(ctx, getItemByName, name)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rarity,
		&i.Cost,
	)
	return i, err
}

const listItem = `-- name: ListItem :many
SELECT id, name, description, rarity, cost from item
`

func (q *Queries) ListItem(ctx context.Context) ([]Item, error) {
	rows, err := q.db.Query(ctx, listItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Rarity,
			&i.Cost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
