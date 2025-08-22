package repository

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq" // драйвер для database/sql
)

type repoPostgres struct {
	db *sql.DB
}

func NewRepoPostgres(db *sql.DB) *repoPostgres {
	return &repoPostgres{db: db}
}

func (r *repoPostgres) SaveURL(ctx context.Context, short string, original string) error {
	const query = `
		INSERT INTO urls (short_url, original_url)
		VALUES ($1, $2)
		ON CONFLICT (short_url) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, short, original)
	return err
}

func (r *repoPostgres) FindByShort(ctx context.Context, short string) (original string, err error) {
	const query = `
		SELECT original_url
		FROM urls
		WHERE short_url = $1
		  AND (expires_at IS NULL OR expires_at > now())
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, query, short).Scan(&original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", sql.ErrNoRows // мапь на ErrNotFound в usecase
		}
		return "", err
	}
	return original, nil
}

func (r *repoPostgres) FindByOriginal(ctx context.Context, original string) (short string, err error) {
	const query = `
		SELECT short_url
		FROM urls
		WHERE original_url = $1
		  AND (expires_at IS NULL OR expires_at > now())
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, query, original).Scan(&short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", sql.ErrNoRows // мапь на ErrNotFound в usecase
		}
		return "", err
	}
	return short, nil
}
