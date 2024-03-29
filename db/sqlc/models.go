// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql"
	"time"
)

type College struct {
	CollegeID int64  `json:"college_id"`
	Name      string `json:"name"`
}

type Funnel struct {
	FunnelID int64  `json:"funnel_id"`
	Name     string `json:"name"`
}

type Invoice struct {
	InvoiceID       int64     `json:"invoice_id"`
	StudentID       int64     `json:"student_id"`
	LessonID        int64     `json:"lesson_id"`
	InvoiceDatetime time.Time `json:"invoice_datetime"`
	// hourly fee for the lesson
	HourlyFee float64 `json:"hourly_fee"`
	// lesson duration in minutes
	Duration int64   `json:"duration"`
	Discount float64 `json:"discount"`
	// total amount based on lesson duration, hourly fee and discount
	Amount float64        `json:"amount"`
	Notes  sql.NullString `json:"notes"`
}

type Lesson struct {
	LessonID       int64     `json:"lesson_id"`
	LessonDatetime time.Time `json:"lesson_datetime"`
	// lesson duration in minutes
	Duration   int64          `json:"duration"`
	LocationID int64          `json:"location_id"`
	SubjectID  int64          `json:"subject_id"`
	Notes      sql.NullString `json:"notes"`
}

type LessonLocation struct {
	LocationID int64  `json:"location_id"`
	Name       string `json:"name"`
}

type LessonSubject struct {
	SubjectID int64  `json:"subject_id"`
	Name      string `json:"name"`
}

type Payment struct {
	PaymentID       int64     `json:"payment_id"`
	ReceiptID       int64     `json:"receipt_id"`
	PaymentDatetime time.Time `json:"payment_datetime"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int64     `json:"payment_method_id"`
}

type PaymentMethod struct {
	PaymentMethodID int64  `json:"payment_method_id"`
	Name            string `json:"name"`
}

type Receipt struct {
	ReceiptID       int64     `json:"receipt_id"`
	StudentID       int64     `json:"student_id"`
	ReceiptDatetime time.Time `json:"receipt_datetime"`
	// total amount of all payments
	Amount float64        `json:"amount"`
	Notes  sql.NullString `json:"notes"`
}

type Student struct {
	StudentID   int64          `json:"student_id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Email       sql.NullString `json:"email"`
	PhoneNumber sql.NullString `json:"phone_number"`
	Address     sql.NullString `json:"address"`
	CollegeID   sql.NullInt64  `json:"college_id"`
	FunnelID    sql.NullInt64  `json:"funnel_id"`
	// hourly fee for the student
	HourlyFee sql.NullFloat64 `json:"hourly_fee"`
	Notes     sql.NullString  `json:"notes"`
	CreatedAt time.Time       `json:"created_at"`
}
