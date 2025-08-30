package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"

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
		Values(token.ID, token.UserID, token.ExpiresAt, token.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql, args...)
	return err
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

func (r *TokenRepository) GetByID(tokenID string) (*models.RefreshToken, error)  {
	query := r.qb.Select("id", "user_id", "token", "expires_at", "created_at").
		From("refresh_tokens").
		Where(squirrel.Eq{"id": tokenID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var token models.RefreshToken
	err = r.db.QueryRow(sql, args...).Scan(&token.ID, &token.UserID, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
