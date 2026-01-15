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

func TestUserRepository_CreateUser_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "test@example.com", "hashedpassword", false).
		WillReturnRows(rows)

	// Execute
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err = repo.CreateUser(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByUsername_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "username", "email", "password", "is_verified", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", "hashedpassword", true, now, now)

	mock.ExpectQuery("SELECT id, username, email").
		WithArgs("testuser").
		WillReturnRows(rows)

	// Execute
	user, err := repo.GetUserByUsername(context.Background(), "testuser")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.True(t, user.IsVerified)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByUsername_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	mock.ExpectQuery("SELECT id, username, email").
		WithArgs("nonexistent").
		WillReturnError(pgx.ErrNoRows)

	// Execute
	user, err := repo.GetUserByUsername(context.Background(), "nonexistent")

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "username", "email", "password", "is_verified", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", "hashedpassword", true, now, now)

	mock.ExpectQuery("SELECT id, username, email").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Execute
	user, err := repo.GetUserByEmail(context.Background(), "test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "username", "email", "password", "is_verified", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", "hashedpassword", true, now, now)

	mock.ExpectQuery("SELECT id, username, email").
		WithArgs(1).
		WillReturnRows(rows)

	// Execute
	user, err := repo.GetUserByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateSession_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	rows := pgxmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, now)

	mock.ExpectQuery("INSERT INTO user_sessions").
		WithArgs(1, "session-token-123", pgxmock.AnyArg()).
		WillReturnRows(rows)

	// Execute
	session := &models.UserSession{
		UserID:    1,
		Token:     "session-token-123",
		ExpiresAt: expiresAt,
	}
	err = repo.CreateSession(context.Background(), session)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, session.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetSessionByToken_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	rows := pgxmock.NewRows([]string{"id", "user_id", "token", "expires_at", "created_at"}).
		AddRow(1, 1, "session-token-123", expiresAt, now)

	mock.ExpectQuery("SELECT id, user_id, token").
		WithArgs("session-token-123").
		WillReturnRows(rows)

	// Execute
	session, err := repo.GetSessionByToken(context.Background(), "session-token-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "session-token-123", session.Token)
	assert.Equal(t, 1, session.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteSession_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(&mockDB{pool: mock})

	mock.ExpectExec("DELETE FROM user_sessions WHERE token").
		WithArgs("session-token-123").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Execute
	err = repo.DeleteSession(context.Background(), "session-token-123")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}


