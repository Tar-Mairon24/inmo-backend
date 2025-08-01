package repository_test

import (
	"regexp"
	"testing"
	"time"

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

    // Mock the INSERT query with transaction
    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`username`,`email`,`password`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)")).
        WithArgs(
            user.Username,
            user.Email,
            user.Password,
            sqlmock.AnyArg(), // created_at
            sqlmock.AnyArg(), // updated_at
        ).WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // Execute test
    err := repo.Create(user)

    // Assertions
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAll(t *testing.T) {
    _, mock, cleanup := helpers.SetupMockDB(t)
    defer cleanup()

    repo := repository.NewUserRepository()

    // Mock the SELECT query
    now := time.Now()
    rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
        AddRow(1, "user1", "user1@example.com", "password1", now, now).
        AddRow(2, "user2", "user2@example.com", "password2", now.Add(time.Hour), now.Add(time.Hour))

    mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
        WillReturnRows(rows)

    // Execute test
    users, err := repo.GetAll()

    // Assertions
    assert.NoError(t, err)
    assert.Len(t, users, 2)
    assert.Equal(t, "user1", users[0].Username)
    assert.Equal(t, "user2", users[1].Username)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
    _, mock, cleanup := helpers.SetupMockDB(t)
    defer cleanup()

    repo := repository.NewUserRepository()

    t.Run("Success", func(t *testing.T) {
        mock.ExpectationsWereMet()
        
		now := time.Now()
        rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
            AddRow(1, "testuser", "test@example.com", "password123", now, now)

        mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?")).
            WithArgs(1,1).
            WillReturnRows(rows)

        user, err := repo.GetByID(1)

        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, uint(1), user.ID)
        assert.Equal(t, "testuser", user.Username)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("NotFound", func(t *testing.T) {
        // Setup new mock for this subtest
        _, newMock, newCleanup := helpers.SetupMockDB(t)
        defer newCleanup()

        newMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?")).
            WithArgs(999,1).
            WillReturnError(gorm.ErrRecordNotFound)

        user, err := repo.GetByID(999)

        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Equal(t, gorm.ErrRecordNotFound, err)
        assert.NoError(t, newMock.ExpectationsWereMet())
    })
}

func TestUserRepository_Update(t *testing.T) {
    _, mock, cleanup := helpers.SetupMockDB(t)
    defer cleanup()

    repo := repository.NewUserRepository()

    t.Run("Success", func(t *testing.T) {
        user := &models.User{
            ID:       1,
            Username: "updated_user",
            Email:    "updated@example.com",
        }

		now := time.Now()
        rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
            AddRow(1, "old_user", "old@example.com", "password", now, now)

        mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?")).
            WithArgs(1, 1).
            WillReturnRows(rows)

        // Mock the update query
        mock.ExpectBegin()
        mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `id`=?, `username`=?,`email`=?,`updated_at`=? WHERE id = ?")).
            WithArgs(1, "updated_user", "updated@example.com", sqlmock.AnyArg(), 1).
            WillReturnResult(sqlmock.NewResult(0, 1))
        mock.ExpectCommit()

        err := repo.Update(user)

        assert.NoError(t, err)
        assert.NoError(t, mock.ExpectationsWereMet()) 
    })

    t.Run("UserNotFound", func(t *testing.T) {
        _, newMock, newCleanup := helpers.SetupMockDB(t)
        defer newCleanup()

        user := &models.User{ID: 999, Username: "updated"}

        newMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
            WithArgs(999).
            WillReturnError(gorm.ErrRecordNotFound)

        err := repo.Update(user)

        assert.Error(t, err)
        assert.Equal(t, gorm.ErrRecordNotFound, err)
        assert.NoError(t, newMock.ExpectationsWereMet())
    })
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

        err := repo.Delete(1)

        assert.NoError(t, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("UserNotFound", func(t *testing.T) {
        // Setup new mock for this subtest
        _, newMock, newCleanup := helpers.SetupMockDB(t)
        defer newCleanup()

        newMock.ExpectBegin()
        newMock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
            WithArgs(sqlmock.AnyArg(), 999).
            WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
        newMock.ExpectCommit()

        err := repo.Delete(999)
		assert.NoError(t, err) // Depending on your implementation, this might not error

        // This might not error depending on your implementation
        // If your Delete method checks affected rows, it should error
        // If not, it might succeed even for non-existent IDs
        assert.NoError(t, newMock.ExpectationsWereMet())
    })
}