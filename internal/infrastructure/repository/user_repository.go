package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type UserRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *UserRepository) ConsultPassword(email string) (string, error) {
	query := r.qb.Select("password").
		From("users").
		Where(squirrel.Eq{"email": email})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for login")
		return "", err
	}

	var password string
	err = r.db.QueryRow(sqlStr, args...).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("No user found with the provided email")
			return "", errors.New("user not found")
		}
		logrus.WithError(err).Error("Failed to execute query for login")
		return "", err
	}

	return password, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.UserResponse, error) {
	query := r.qb.Select("id", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"email": email})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for consulting user by email")
		return nil, err
	}

	var dbUser models.UserResponse
	err = r.db.QueryRow(sqlStr, args...).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("No user found with the provided email")
			return nil, errors.New("user not found")
		}
		logrus.WithError(err).Error("Failed to execute query for consulting user by email")
		return nil, err
	}

	logrus.Info("User found successfully in UserRepositoryImpl")
	return &dbUser, nil
}

func (r *UserRepository) Create(user *models.User) error {
	query := r.qb.Insert("users").
		Columns("username", "email", "password", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, time.Now(), time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

func (r *UserRepository) GetAll() ([]models.UserResponse, error) {
	logrus.Info("Retrieving all users from the database")
	if r.db == nil {
		logrus.Error("Database connection is nil")
		return nil, errors.New("database connection is not initialized")
	}
	query := r.qb.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query")
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query to get all users")
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			logrus.WithError(err).Error("Failed to scan user row")
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("Error occurred while iterating over user rows")
		return nil, err
	}
	logrus.Infof("Retrieved %d users from the database", len(users))
	return users, nil
}

func (r *UserRepository) GetByID(id uint) (*models.UserResponse, error) {
	query := r.qb.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user models.UserResponse
	err = r.db.QueryRow(sqlStr, args...).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := r.qb.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": user.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) Delete(id uint) error {
	query := r.qb.Delete("users").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}

	return nil
}
