// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: receipt.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createReceipt = `-- name: CreateReceipt :one
INSERT INTO receipts (
  student_id, receipt_datetime, amount, notes
) VALUES (
  $1, $2, $3, $4
)
RETURNING receipt_id, student_id, receipt_datetime, amount, notes
`

type CreateReceiptParams struct {
	StudentID       int64          `json:"student_id"`
	ReceiptDatetime time.Time      `json:"receipt_datetime"`
	Amount          float64        `json:"amount"`
	Notes           sql.NullString `json:"notes"`
}

func (q *Queries) CreateReceipt(ctx context.Context, arg CreateReceiptParams) (Receipt, error) {
	row := q.db.QueryRowContext(ctx, createReceipt,
		arg.StudentID,
		arg.ReceiptDatetime,
		arg.Amount,
		arg.Notes,
	)
	var i Receipt
	err := row.Scan(
		&i.ReceiptID,
		&i.StudentID,
		&i.ReceiptDatetime,
		&i.Amount,
		&i.Notes,
	)
	return i, err
}

const deleteReceipt = `-- name: DeleteReceipt :exec
DELETE FROM receipts
WHERE receipt_id = $1
`

func (q *Queries) DeleteReceipt(ctx context.Context, receiptID int64) error {
	_, err := q.db.ExecContext(ctx, deleteReceipt, receiptID)
	return err
}

const getReceipt = `-- name: GetReceipt :one
SELECT receipt_id, student_id, receipt_datetime, amount, notes FROM receipts
WHERE receipt_id = $1 LIMIT 1
`

func (q *Queries) GetReceipt(ctx context.Context, receiptID int64) (Receipt, error) {
	row := q.db.QueryRowContext(ctx, getReceipt, receiptID)
	var i Receipt
	err := row.Scan(
		&i.ReceiptID,
		&i.StudentID,
		&i.ReceiptDatetime,
		&i.Amount,
		&i.Notes,
	)
	return i, err
}

const getReceiptsByStudent = `-- name: GetReceiptsByStudent :many
SELECT receipt_id, student_id, receipt_datetime, amount, notes FROM receipts
WHERE student_id = $1
ORDER BY receipt_datetime
LIMIT $2
OFFSET $3
`

type GetReceiptsByStudentParams struct {
	StudentID int64 `json:"student_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) GetReceiptsByStudent(ctx context.Context, arg GetReceiptsByStudentParams) ([]Receipt, error) {
	rows, err := q.db.QueryContext(ctx, getReceiptsByStudent, arg.StudentID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Receipt{}
	for rows.Next() {
		var i Receipt
		if err := rows.Scan(
			&i.ReceiptID,
			&i.StudentID,
			&i.ReceiptDatetime,
			&i.Amount,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listReceipts = `-- name: ListReceipts :many
SELECT receipt_id, student_id, receipt_datetime, amount, notes FROM receipts
ORDER BY student_id, receipt_datetime
LIMIT $1
OFFSET $2
`

type ListReceiptsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListReceipts(ctx context.Context, arg ListReceiptsParams) ([]Receipt, error) {
	rows, err := q.db.QueryContext(ctx, listReceipts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Receipt{}
	for rows.Next() {
		var i Receipt
		if err := rows.Scan(
			&i.ReceiptID,
			&i.StudentID,
			&i.ReceiptDatetime,
			&i.Amount,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateReceipt = `-- name: UpdateReceipt :exec
UPDATE receipts
  set   student_id = $2,
        receipt_datetime = $3, 
        amount = $4,
        notes = $5
WHERE receipt_id = $1
`

type UpdateReceiptParams struct {
	ReceiptID       int64          `json:"receipt_id"`
	StudentID       int64          `json:"student_id"`
	ReceiptDatetime time.Time      `json:"receipt_datetime"`
	Amount          float64        `json:"amount"`
	Notes           sql.NullString `json:"notes"`
}

func (q *Queries) UpdateReceipt(ctx context.Context, arg UpdateReceiptParams) error {
	_, err := q.db.ExecContext(ctx, updateReceipt,
		arg.ReceiptID,
		arg.StudentID,
		arg.ReceiptDatetime,
		arg.Amount,
		arg.Notes,
	)
	return err
}

const updateReceiptAmount = `-- name: UpdateReceiptAmount :exec
UPDATE receipts
  set   amount = $2
WHERE receipt_id = $1
`

type UpdateReceiptAmountParams struct {
	ReceiptID int64   `json:"receipt_id"`
	Amount    float64 `json:"amount"`
}

func (q *Queries) UpdateReceiptAmount(ctx context.Context, arg UpdateReceiptAmountParams) error {
	_, err := q.db.ExecContext(ctx, updateReceiptAmount, arg.ReceiptID, arg.Amount)
	return err
}
