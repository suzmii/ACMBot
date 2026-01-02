package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suzmii/ACMBot/database/sqlc"
	"github.com/suzmii/ACMBot/util/logx"
)

var logger = logx.New("database")

type dbStore struct {
	sqlc.Querier
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Fatalf("无法连接到数据库: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("无法连接至数据库: %w", err)
	}

	logger.Info("成功连接到数据库")

	return &dbStore{
		Querier: sqlc.New(pool),
		pool:    pool,
	}, nil
}

func (s *dbStore) StartTx(
	ctx context.Context,
	fn func(q Store) error,
) error {
	if s == nil || s.pool == nil {
		return errors.New("transaction can only be started from root store")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	txStore := &dbStore{
		Querier: sqlc.New(tx),
		pool:    nil,
	}

	if err := fn(txStore); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %w, rollback err: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
