package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/infrastructure/db"
)

type TokenRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewTokenRepository(db *sql.DB) ports.TokenRepository {
	return &TokenRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *TokenRepository) SaveToken(token *models.RefreshToken) error {
	query := r.qb.Insert("refresh_tokens").
		Columns("id", "user_id", "token", "expires_at", "created_at").
		Values(token.ID, token.UserID, token.Token, token.ExpiresAt, time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to save refresh token")
		return err
	}

	logrus.Info("New refresh token saved successfully")
	return nil
}

func (r *TokenRepository) DeleteToken(tokenID string) error {
	query := r.qb.Delete("refresh_tokens").
		Where(squirrel.Eq{"id": tokenID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			logrus.Warn("No refresh token found with the provided ID")
			return nil
		}
		logrus.WithError(err).Error("Failed to delete refresh token")
		return err
	}

	logrus.Info("Refresh token deleted successfully")
	return nil
}

func (r *TokenRepository) GetTokenIDByUserID(userID uint) (string, error)  {
	query := r.qb.Select("id").
		From("refresh_tokens").
		Where(squirrel.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	var refreshTokenID string
	ctx := context.Background()
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&refreshTokenID)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			logrus.Warn("No refresh token found with the provided user ID")
			return "", errors.New("no token found")
		}
		return "", err
	}

	return refreshTokenID, nil
}

func (r *TokenRepository) GetTokenByUserID(userID uint) (*models.RefreshToken, error) {
	query := r.qb.Select("id", "user_id", "token", "expires_at", "created_at").
		From("refresh_tokens").
		Where(squirrel.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query")
		return nil, err
	}

	var token models.RefreshToken
	ctx := context.Background()
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			logrus.Warn("No refresh token found for the provided user ID")
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to execute query")
		return nil, err
	}

	return &token, nil
}
