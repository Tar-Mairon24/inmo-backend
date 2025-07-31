package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/infrastructure/repository"
	"inmo-backend/test/helpers"
)

func TestUserRepository_Create(t *testing.T) {
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()
	user := helpers.CreateTestUser()

	// Mock the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`username`,`email`,`password`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)"),
	).WithArgs(
		user.Username,
		user.Email,
		user.Password,
		sqlmock.AnyArg(), // created_at
		sqlmock.AnyArg(), // updated_at
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test the method
	err := repo.Create(user)

	// Assertions
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAll(t *testing.T) {
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "user1", "user1@example.com", "password1", "2024-01-01 10:00:00", "2024-01-01 10:00:00").
		AddRow(2, "user2", "user2@example.com", "password2", "2024-01-01 11:00:00", "2024-01-01 11:00:00")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	// Test the method
	users, err := repo.GetAll()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()

	t.Run("Success", func(t *testing.T) {
		// Mock the SELECT query
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
			AddRow(1, "testuser", "test@example.com", "password123", "2024-01-01 10:00:00", "2024-01-01 10:00:00")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
			WithArgs(1).
			WillReturnRows(rows)

		// Test the method
		user, err := repo.GetByID(1)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "testuser", user.Username)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
			WithArgs(999).
			WillReturnError(gorm.ErrRecordNotFound)

		// Test the method
		user, err := repo.GetByID(999)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Test the method
		err := repo.Delete(1)

		// Assertions
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
			WithArgs(sqlmock.AnyArg(), 999).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		// Test the method
		err := repo.Delete(999)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAll_Success(t *testing.T) {
	// Setup
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()

	// Mock data
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, "user1", "user1@test.com", "password1", "2024-01-01", "2024-01-01").
		AddRow(2, "user2", "user2@test.com", "password2", "2024-01-01", "2024-01-01")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	// Execute
	users, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()

	// Mock expectation
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	// Execute
	user, err := repo.GetByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create_Success(t *testing.T) {
	// Setup
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()
	user := helpers.CreateTestUser()

	// Mock expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(user.Username, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Execute
	err := repo.Create(user)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_UserNotFound(t *testing.T) {
	// Setup
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()
	user := &models.User{ID: 999, Username: "updated"}

	// Mock user not found
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?")).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	// Execute
	err := repo.Update(user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete_Success(t *testing.T) {
	// Setup
	_, mock, cleanup := helpers.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository()

	// Mock soft delete
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ?")).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Execute
	err := repo.Delete(1)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
