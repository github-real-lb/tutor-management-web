// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: payment.sql

package db

import (
	"context"
	"time"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (
  receipt_id, payment_datetime, amount, payment_method_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING payment_id, receipt_id, payment_datetime, amount, payment_method_id
`

type CreatePaymentParams struct {
	ReceiptID       int64     `json:"receipt_id"`
	PaymentDatetime time.Time `json:"payment_datetime"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int64     `json:"payment_method_id"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.ReceiptID,
		arg.PaymentDatetime,
		arg.Amount,
		arg.PaymentMethodID,
	)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.ReceiptID,
		&i.PaymentDatetime,
		&i.Amount,
		&i.PaymentMethodID,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payments
WHERE payment_id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, paymentID int64) error {
	_, err := q.db.ExecContext(ctx, deletePayment, paymentID)
	return err
}

const deletePaymentsByReceipt = `-- name: DeletePaymentsByReceipt :exec
DELETE FROM payments
WHERE receipt_id = $1
`

func (q *Queries) DeletePaymentsByReceipt(ctx context.Context, receiptID int64) error {
	_, err := q.db.ExecContext(ctx, deletePaymentsByReceipt, receiptID)
	return err
}

const getPayment = `-- name: GetPayment :one
SELECT payment_id, receipt_id, payment_datetime, amount, payment_method_id FROM payments
WHERE payment_id = $1 LIMIT 1
`

func (q *Queries) GetPayment(ctx context.Context, paymentID int64) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPayment, paymentID)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.ReceiptID,
		&i.PaymentDatetime,
		&i.Amount,
		&i.PaymentMethodID,
	)
	return i, err
}

const getPayments = `-- name: GetPayments :many
SELECT payment_id, receipt_id, payment_datetime, amount, payment_method_id FROM payments
WHERE receipt_id = $1 
ORDER BY receipt_id, payment_datetime
`

func (q *Queries) GetPayments(ctx context.Context, receiptID int64) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, getPayments, receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Payment{}
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.PaymentID,
			&i.ReceiptID,
			&i.PaymentDatetime,
			&i.Amount,
			&i.PaymentMethodID,
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

const listPayments = `-- name: ListPayments :many
SELECT payment_id, receipt_id, payment_datetime, amount, payment_method_id FROM payments
ORDER BY receipt_id, payment_datetime
LIMIT $1
OFFSET $2
`

type ListPaymentsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPayments(ctx context.Context, arg ListPaymentsParams) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, listPayments, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Payment{}
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.PaymentID,
			&i.ReceiptID,
			&i.PaymentDatetime,
			&i.Amount,
			&i.PaymentMethodID,
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

const updatePayment = `-- name: UpdatePayment :exec
UPDATE payments
  set   receipt_id = $2,
        payment_datetime = $3, 
        amount = $4,
        payment_method_id = $5
WHERE payment_id = $1
`

type UpdatePaymentParams struct {
	PaymentID       int64     `json:"payment_id"`
	ReceiptID       int64     `json:"receipt_id"`
	PaymentDatetime time.Time `json:"payment_datetime"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int64     `json:"payment_method_id"`
}

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) error {
	_, err := q.db.ExecContext(ctx, updatePayment,
		arg.PaymentID,
		arg.ReceiptID,
		arg.PaymentDatetime,
		arg.Amount,
		arg.PaymentMethodID,
	)
	return err
}
