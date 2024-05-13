package key

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type dataKeyRepository struct {
	pool       *pgxpool.Pool
	connString string
}

func NewDataKeyRepository(ctx context.Context, connString string) (*dataKeyRepository, error) {
	const op = "new user repository"

	var err error
	pool, err := initPool(ctx, connString)
	if err != nil {
		return nil, errors.Wrapf(err, op)
	}

	r := &dataKeyRepository{
		pool,
		connString,
	}
	return r, nil
}

func (r *dataKeyRepository) Close() {
	if r.pool == nil {
		return
	}
	r.pool.Close()
}

func (r *dataKeyRepository) RotateKey(ctx context.Context, key *model.DataKey) (*model.DataKey, error) {
	const op = "rotate key"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer conn.Release()

	var txOptions pgx.TxOptions
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `UPDATE keys SET is_disposed = 'true'`)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	const sql = `INSERT INTO keys (key_data) VALUES ($1) RETURNING key_id;`
	row := tx.QueryRow(ctx, sql, key.Key)
	err = row.Scan(&key.ID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return key, nil
}

func (r *dataKeyRepository) GetKey(ctx context.Context) (*model.DataKey, error) {
	const op = "get key"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer conn.Release()

	key := &model.DataKey{}
	const sql = `SELECT key_id, key_data, COALESCE(encriptions_count, 0), COALESCE(encrypted_data_size, 0)
FROM keys WHERE is_disposed='false'`
	row := conn.QueryRow(ctx, sql)
	err = row.Scan(&key.ID, &key.Key, &key.EncryptionsCount, &key.EncryptedDataSize)
	if errors.Is(err, pgx.ErrNoRows) {
		var k *model.DataKey
		return k, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return key, nil
}

func (r *dataKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.DataKey, error) {
	const op = "get by id"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer conn.Release()

	key := &model.DataKey{}
	const sql = `SELECT key_id, key_data, COALESCE(encriptions_count, 0), COALESCE(encrypted_data_size, 0)
FROM keys WHERE key_id=$1`
	row := conn.QueryRow(ctx, sql, id)
	err = row.Scan(&key.ID, &key.Key, &key.EncryptionsCount, &key.EncryptedDataSize)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, vault.ErrKeyNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return key, nil
}

func (r *dataKeyRepository) UpdateStats(ctx context.Context, id uuid.UUID, dataSize int64) error {
	const op = "update stats"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer conn.Release()

	var txOptions pgx.TxOptions
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var count int
	var size int64
	const sql = `SELECT COALESCE(encriptions_count, 0), COALESCE(encrypted_data_size, 0)
FROM keys WHERE key_id=$1`
	row := tx.QueryRow(ctx, sql, id)
	err = row.Scan(&count, &size)
	if errors.Is(err, pgx.ErrNoRows) {
		return vault.ErrKeyNotFound
	}
	if err != nil {
		return errors.Wrap(err, op)
	}

	count++
	size += dataSize
	_, err = tx.Exec(ctx, `UPDATE keys SET encriptions_count=$1, encrypted_data_size=$2 WHERE key_id=$3`, count, size, id)
	if err != nil {
		return errors.Wrap(err, op)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func initPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	const op = "init pool"

	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, op)
	}

	return pool, nil
}
