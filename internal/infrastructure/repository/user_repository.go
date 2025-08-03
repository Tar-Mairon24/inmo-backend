package repository

import (
	"errors"
	"time"

	"database/sql"
	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct{
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepositoryImpl{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
    logrus.Info("Login method called in UserRepositoryImpl")

    query := r.qb.Select("id", "username", "email", "password").
        From("users").
        Where(squirrel.Eq{"email": email})

    sqlStr, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for login")
        return nil, err
    }

    var dbUser models.User
    err = r.db.QueryRow(sqlStr, args...).Scan(
        &dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            logrus.Warn("No user found with the provided email")
            return nil,errors.New("user not found")
        }
        logrus.WithError(err).Error("Failed to execute query for login")
        return nil,err
    }

    logrus.Info("User found successfully in UserRepositoryImpl")
    return &dbUser, nil
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
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

func (r *UserRepositoryImpl) GetAll() ([]models.User, error) {
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

	var users []models.User
	for rows.Next() {
		var user models.User
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

func (r *UserRepositoryImpl) GetByID(id uint) (*models.User, error) {
    query := r.qb.Select("id", "username", "email", "password", "created_at", "updated_at").
        From("users").
        Where(squirrel.Eq{"id": id})

    sqlStr, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }

    var user models.User
    err = r.db.QueryRow(sqlStr, args...).Scan(
        &user.ID, &user.Username, &user.Email, 
        &user.Password, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    

    return &user, nil
}

func (r *UserRepositoryImpl) Update(user *models.User) error {
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

func (r *UserRepositoryImpl) Delete(id uint) error {
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