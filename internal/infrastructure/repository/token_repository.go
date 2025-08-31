package repository

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
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

	_, err = r.db.Exec(sql, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to save refresh token")
		return err
	}

	logrus.Infof("Refresh token saved successfully: %s", token.ID)
	return nil
}

func (r *TokenRepository) DeleteToken(tokenID string) error {
	query := r.qb.Delete("refresh_tokens").
		Where(squirrel.Eq{"id": tokenID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql, args...)
	return err
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
	err = r.db.QueryRow(sql, args...).Scan(&refreshTokenID)
	if err != nil {
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
	err = r.db.QueryRow(sql, args...).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query or scan result")
		return nil, err
	}

	return &token, nil
}
