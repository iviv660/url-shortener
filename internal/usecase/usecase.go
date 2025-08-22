package usecase

import (
	"app/internal/validate"
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type UseCase struct {
	repo       RepoPostgres
	cache      RepoRedis
	cacheTTL   time.Duration
	secretCode string
}

func New(repo RepoPostgres, cache RepoRedis, cacheTTL time.Duration, secretCode string) *UseCase {
	return &UseCase{
		repo:       repo,
		cache:      cache,
		cacheTTL:   cacheTTL,
		secretCode: secretCode,
	}
}

func (u *UseCase) CreateShortURL(ctx context.Context, original string) (string, error) {
	// 1) Валидируем + нормализуем входной URL
	normURL, err := validate.Validate(ctx, original)
	if err != nil {
		return "", err
	}

	// 2) Идемпотентность: если уже есть short для этого normURL — вернём его
	if existing, err := u.repo.FindByOriginal(ctx, normURL); err == nil && existing != "" {
		_ = u.cache.Set(ctx, existing, normURL, u.cacheTTL)
		return existing, nil
	}

	// 3) Генерируем короткий код
	short := generateShort(normURL, u.secretCode)

	// 4) Сохраняем новую запись (теперь нужно чтобы SaveURL принимал short и original)
	if err := u.repo.SaveURL(ctx, short, normURL); err != nil {
		return "", err
	}

	// 5) Пишем в кэш (best-effort)
	_ = u.cache.Set(ctx, short, normURL, u.cacheTTL)

	return short, nil
}

func (u *UseCase) OpenUrl(ctx context.Context, short string) (string, error) {
	// 1) Пытаемся из кэша
	if long, err := u.cache.Get(ctx, short); err == nil {
		return long, nil
	} else if !errors.Is(err, ErrCacheMiss) {
		// реальная ошибка Redis — отдаём вверх
		return "", err
	}

	// 2) Из БД
	long, err := u.repo.FindByShort(ctx, short)
	if err != nil {
		return "", err
	}

	// 3) Backfill в кэш (best-effort)
	_ = u.cache.Set(ctx, short, long, u.cacheTTL)

	return long, nil
}
