package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomInvoice tests adding a new random invoice to the database, and returns the Invoice data type.
func createRandomInvoice(t *testing.T) Invoice {
	student := createRandomStudent(t)
	lesson := createRandomLesson(t)

	arg := CreateInvoiceParams{
		StudentID:       student.StudentID,
		LessonID:        lesson.LessonID,
		InvoiceDatetime: util.RandomDatetime(),
		HourlyFee:       util.RandomFloat64(85.0, 300.0),
		Amount:          util.RandomFloat64(85.0, 1200.0),
		Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
	}

	invoice, err := testQueries.CreateInvoice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	require.Equal(t, arg.StudentID, invoice.StudentID)
	require.Equal(t, arg.LessonID, invoice.LessonID)
	require.WithinDuration(t, arg.InvoiceDatetime, invoice.InvoiceDatetime, time.Second)
	require.Equal(t, arg.HourlyFee, invoice.HourlyFee)
	require.Equal(t, arg.Amount, invoice.Amount)
	require.Equal(t, arg.Notes, invoice.Notes)

	require.NotZero(t, invoice.InvoiceID)

	return invoice
}

func TestCreateInvoice(t *testing.T) {
	createRandomInvoice(t)
}

func TestGetInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	invoice2, err := testQueries.GetInvoice(context.Background(), invoice1.InvoiceID)
	require.NoError(t, err)
	require.NotEmpty(t, invoice2)

	require.Equal(t, invoice1.InvoiceID, invoice2.InvoiceID)
	require.Equal(t, invoice1.StudentID, invoice2.StudentID)
	require.Equal(t, invoice1.LessonID, invoice2.LessonID)
	require.WithinDuration(t, invoice1.InvoiceDatetime, invoice2.InvoiceDatetime, time.Second)
	require.Equal(t, invoice1.HourlyFee, invoice2.HourlyFee)
	require.Equal(t, invoice1.Amount, invoice2.Amount)
	require.Equal(t, invoice1.Notes, invoice2.Notes)
}

func TestUpdateInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	student := createRandomStudent(t)
	lesson := createRandomLesson(t)

	arg := UpdateInvoiceParams{
		InvoiceID:       invoice1.InvoiceID,
		StudentID:       student.StudentID,
		LessonID:        lesson.LessonID,
		InvoiceDatetime: util.RandomDatetime(),
		HourlyFee:       util.RandomFloat64(85.0, 300.0),
		Amount:          util.RandomFloat64(85.0, 1200.0),
		Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
	}
	err := testQueries.UpdateInvoice(context.Background(), arg)
	require.NoError(t, err)

	invoice2, err := testQueries.GetInvoice(context.Background(), arg.InvoiceID)
	require.NoError(t, err)
	require.NotEmpty(t, invoice2)

	require.Equal(t, arg.InvoiceID, invoice2.InvoiceID)
	require.Equal(t, arg.StudentID, invoice2.StudentID)
	require.Equal(t, arg.LessonID, invoice2.LessonID)
	require.WithinDuration(t, arg.InvoiceDatetime, invoice2.InvoiceDatetime, time.Second)
	require.Equal(t, arg.HourlyFee, invoice2.HourlyFee)
	require.Equal(t, arg.Amount, invoice2.Amount)
	require.Equal(t, arg.Notes, invoice2.Notes)
}

func TestDeleteInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	testQueries.DeleteInvoice(context.Background(), invoice1.InvoiceID)

	invoice2, err := testQueries.GetInvoice(context.Background(), invoice1.InvoiceID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, invoice2)
}

func TestListInvoices(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomInvoice(t)
	}

	arg := ListInvoicesParams{
		Limit:  5,
		Offset: 5,
	}

	invoices, err := testQueries.ListInvoices(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, invoices, 5)

	for _, invoice := range invoices {
		require.NotEmpty(t, invoice)
	}
}