package db

import (
	"context"
	"database/sql"
	"time"
)

// Payments is used for all payments of a single receipt
type Payments []Payment

// Implement the sort.Interface methods for the Payments type
func (p Payments) Len() int           { return len(p) }
func (p Payments) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Payments) Less(i, j int) bool { return p[i].PaymentDatetime.Before(p[j].PaymentDatetime) }

type Receipts []Receipt

// Implement the sort.Interface methods for the Receipts type
func (r Receipts) Len() int           { return len(r) }
func (r Receipts) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Receipts) Less(i, j int) bool { return r[i].ReceiptDatetime.Before(r[j].ReceiptDatetime) }

// ReceiptPayments is used for a single receipt and all its payments.
type ReceiptWithPayments struct {
	Receipt  Receipt  `json:"receipt"`
	Payments Payments `json:"payments"`
}

type ReceiptsWithPayments []ReceiptWithPayments

// Implement the sort.Interface methods for the ReceiptsWithPayments type
func (r ReceiptsWithPayments) Len() int      { return len(r) }
func (r ReceiptsWithPayments) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ReceiptsWithPayments) Less(i, j int) bool {
	return r[i].Receipt.ReceiptDatetime.Before(r[j].Receipt.ReceiptDatetime)
}

// StudentReceipts is used for all receipts and their payments, of a single student.
type StudentReceiptsWithPayments struct {
	StudentID            int64                `json:"student_id"`
	ReceiptsWithPayments ReceiptsWithPayments `json:"receipts_with_payments"`
}

// CreateReceiptTxPaymentParams contains the input paramaters of a single payment, for the CreateReceiptTx function.
type CreateReceiptTxPaymentParams struct {
	PaymentDatetime time.Time `json:"payment_datetime"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int64     `json:"payment_method_id"`
}

// CreateReceiptTxParams contains the input paramaters of a single reciept and its payments, for the CreateReceiptTx function.
type CreateReceiptTxParams struct {
	StudentID             int64                          `json:"student_id"`
	ReceiptDatetime       time.Time                      `json:"receipt_datetime"`
	Notes                 sql.NullString                 `json:"notes"`
	ReceiptPaymentsParams []CreateReceiptTxPaymentParams `json:"receipt_payments_params"`
}

// CreateReceiptTx creates a Receipt and all the Payments releated to it.
// Receipt amount is calculated end updated based on all payments.
func (store *Store) CreateReceiptTx(ctx context.Context, arg CreateReceiptTxParams) (ReceiptWithPayments, error) {
	var result ReceiptWithPayments

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createReceiptArg := CreateReceiptParams{
			StudentID:       arg.StudentID,
			ReceiptDatetime: arg.ReceiptDatetime,
			Amount:          0.0,
			Notes:           arg.Notes,
		}

		result.Receipt, err = q.CreateReceipt(ctx, createReceiptArg)
		if err != nil {
			return err
		}

		for _, paymentArg := range arg.ReceiptPaymentsParams {
			createPaymentArg := CreatePaymentParams{
				ReceiptID:       result.Receipt.ReceiptID,
				PaymentDatetime: paymentArg.PaymentDatetime,
				PaymentMethodID: paymentArg.PaymentMethodID,
			}

			payment, err := q.CreatePayment(ctx, createPaymentArg)
			if err != nil {
				return err
			}

			result.Receipt.Amount += payment.Amount
			result.Payments = append(result.Payments, payment)
		}

		err = q.UpdateReceiptAmount(ctx, UpdateReceiptAmountParams{
			ReceiptID: result.Receipt.ReceiptID,
			Amount:    result.Receipt.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// DeleteReceiptTx deletes a Receipt and all the Payments releated to it.
func (store *Store) DeleteReceiptTx(ctx context.Context, receiptID int64) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := q.DeletePayments(ctx, receiptID)
		if err != nil {
			return err
		}

		err = q.DeleteReceipt(ctx, receiptID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// GetReceiptTx gets a Receipt and all the Payments releated to it.

func (store *Store) GetReceiptTx(ctx context.Context, receiptID int64) (ReceiptWithPayments, error) {
	var result ReceiptWithPayments

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Receipt, err = q.GetReceipt(ctx, receiptID)
		if err != nil {
			return err
		}

		result.Payments, err = q.GetPayments(ctx, receiptID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// GetReceiptsTxByStudent gets all Receipts of a single student, and all the Payments releated to each receipt.
// limit is used to determine the number of rows (row_count) returned by the query.
// offset is used to skip a number of rows before beginning to return the rows.
func (store *Store) GetReceiptsTxByStudent(ctx context.Context, studentID int64, limit, offset int) (StudentReceiptsWithPayments, error) {
	var result StudentReceiptsWithPayments
	result.StudentID = studentID

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		arg := GetReceiptsByStudentParams{
			StudentID: studentID,
			Limit:     int32(limit),
			Offset:    int32(offset),
		}
		receipts, err := q.GetReceiptsByStudent(ctx, arg)
		if err != nil {
			return err
		}

		for _, receipt := range receipts {
			payments, err := q.GetPayments(ctx, receipt.ReceiptID)
			if err != nil {
				return err
			}

			result.ReceiptsWithPayments = append(result.ReceiptsWithPayments, ReceiptWithPayments{
				Receipt:  receipt,
				Payments: payments,
			})
		}

		return nil
	})

	return result, err
}
