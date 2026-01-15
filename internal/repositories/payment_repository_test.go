package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/andre/project-app-bioskop-golang/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRepository_CreatePayment_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO payments").
		WithArgs(1, 1, 150000.0, "credit_card", "pending", "TXN123456").
		WillReturnRows(rows)

	// Execute
	payment := &models.Payment{
		BookingID:     1,
		UserID:        1,
		Amount:        150000.0,
		PaymentMethod: "credit_card",
		Status:        "pending",
		TransactionID: "TXN123456",
	}
	err = repo.CreatePayment(context.Background(), payment)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, payment.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_GetPaymentByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "booking_id", "user_id", "amount", "payment_method", "status", "transaction_id", "created_at", "updated_at"}).
		AddRow(1, 1, 1, 150000.0, "credit_card", "completed", "TXN123456", now, now)

	mock.ExpectQuery("SELECT id, booking_id, user_id").
		WithArgs(1).
		WillReturnRows(rows)

	// Execute
	payment, err := repo.GetPaymentByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, 1, payment.ID)
	assert.Equal(t, "credit_card", payment.PaymentMethod)
	assert.Equal(t, "completed", payment.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_GetPaymentByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	mock.ExpectQuery("SELECT id, booking_id, user_id").
		WithArgs(999).
		WillReturnError(pgx.ErrNoRows)

	// Execute
	payment, err := repo.GetPaymentByID(context.Background(), 999)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, payment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_GetPaymentByBookingID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "booking_id", "user_id", "amount", "payment_method", "status", "transaction_id", "created_at", "updated_at"}).
		AddRow(1, 10, 1, 150000.0, "credit_card", "completed", "TXN123456", now, now)

	mock.ExpectQuery("SELECT id, booking_id, user_id").
		WithArgs(10).
		WillReturnRows(rows)

	// Execute
	payment, err := repo.GetPaymentByBookingID(context.Background(), 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, 10, payment.BookingID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_UpdatePaymentStatus_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	mock.ExpectExec("UPDATE payments SET status").
		WithArgs("completed", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Execute
	err = repo.UpdatePaymentStatus(context.Background(), 1, "completed")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_GetPaymentMethods_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "name", "type", "is_active", "created_at", "updated_at"}).
		AddRow(1, "Credit Card", "card", true, now, now).
		AddRow(2, "Bank Transfer", "transfer", true, now, now).
		AddRow(3, "E-Wallet", "ewallet", true, now, now)

	mock.ExpectQuery("SELECT id, name, type").
		WillReturnRows(rows)

	// Execute
	methods, err := repo.GetPaymentMethods(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, methods, 3)
	assert.Equal(t, "Credit Card", methods[0].Name)
	assert.Equal(t, "Bank Transfer", methods[1].Name)
	assert.Equal(t, "E-Wallet", methods[2].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentRepository_GetUserPayments_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPaymentRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "booking_id", "user_id", "amount", "payment_method", "status", "transaction_id", "created_at", "updated_at"}).
		AddRow(1, 1, 5, 150000.0, "credit_card", "completed", "TXN123456", now, now).
		AddRow(2, 2, 5, 200000.0, "bank_transfer", "pending", "TXN123457", now, now)

	mock.ExpectQuery("SELECT id, booking_id, user_id").
		WithArgs(5).
		WillReturnRows(rows)

	// Execute
	payments, err := repo.GetUserPayments(context.Background(), 5)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, payments, 2)
	assert.Equal(t, "completed", payments[0].Status)
	assert.Equal(t, "pending", payments[1].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}
