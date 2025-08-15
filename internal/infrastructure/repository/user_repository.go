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
		Where(squirrel.Eq{"email": email}).
		Where(squirrel.Expr("deleted_at IS NULL")) // Ensure deleted_at is NULL

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
	query := r.qb.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Where(squirrel.Expr("deleted_at IS NULL")) // Ensure deleted_at is NULL

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for consulting user by email")
		return nil, err
	}

	var User models.UserResponse
	err = r.db.QueryRow(sqlStr, args...).Scan(
		&User.ID, &User.Username, &User.Email, &User.CreatedAt, &User.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("No user found with the provided email")
			return nil, errors.New("user not found")
		}
		logrus.WithError(err).Error("Failed to execute query for consulting user by email")
		return nil, err
	}

	logrus.Infof("User found successfully with email: %s", email)
	return &User, nil
}

func (r *UserRepository) Create(user *models.User) (*models.UserResponse, error) {
	query := r.qb.Insert("users").
		Columns("username", "email", "password", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, time.Now(), time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query to create user")
		return nil, err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
        logrus.WithError(err).Error("Failed to execute query to create user")
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
        logrus.WithError(err).Error("Failed to retrieve last insert ID")
		return nil, err
	}

	user.ID = uint(id)
    logrus.Infof("User created successfully with ID: %d", user.ID)

	return user.ToUserResponse(), nil
}

func (r *UserRepository) GetAll() ([]models.UserResponse, error) {
	logrus.Info("Retrieving all users from the database")
	if r.db == nil {
		logrus.Error("Database connection is nil")
		return nil, errors.New("database connection is not initialized")
	}
	query := r.qb.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Expr("deleted_at IS NULL")).
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
	defer func() {
        if err := rows.Close(); err != nil {
            logrus.WithError(err).Error("Failed to close database rows")
        }
    }()

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
		Where(squirrel.And{
			squirrel.Eq{"id": id},
			squirrel.Expr("deleted_at IS NULL"),
	})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting user by ID")
		return nil, err
	}

	var user models.UserResponse
	err = r.db.QueryRow(sqlStr, args...).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warnf("No user found with ID: %d", id)
			return nil, errors.New("user not found")
		}
		logrus.WithError(err).Error("Failed to execute query for getting user by ID")
		return nil, err
	}

	logrus.Infof("User found successfully with ID: %d", id)
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) (*models.UserResponse, error) {
	query := r.qb.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": user.ID}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return user.ToUserResponse(), nil
}

func (r *UserRepository) Delete(id uint) error {
	query := r.qb.Update("users").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sql, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for deleting user")
		return err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for deleting user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve rows affected for deleting user")
		return err
	}

	if rowsAffected == 0 {
		logrus.Warnf("No user found with ID: %d or already deleted", id)
		return errors.New("user not found or already deleted")
	}

	logrus.Infof("User with ID: %d deleted successfully", id)
	return nil
}


