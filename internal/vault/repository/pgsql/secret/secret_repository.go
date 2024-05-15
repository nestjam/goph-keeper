package secret

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type secretRepository struct {
	pool       *pgxpool.Pool
	connString string
}

func NewSecretRepository(ctx context.Context, connString string) (*secretRepository, error) {
	const op = "new secret repository"

	var err error
	pool, err := initPool(ctx, connString)
	if err != nil {
		return nil, errors.Wrapf(err, op)
	}

	r := &secretRepository{
		pool,
		connString,
	}
	return r, nil
}

func (r *secretRepository) Close() {
	if r.pool == nil {
		return
	}
	r.pool.Close()
}

func (r *secretRepository) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	const op = "list secrets"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer conn.Release()

	const sql = "SELECT secret_id, name FROM secrets WHERE user_id=$1"
	rows, err := conn.Query(ctx, sql, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer rows.Close()

	var secrets []*model.Secret
	for rows.Next() {
		secret := &model.Secret{}
		err := rows.Scan(&secret.ID, &secret.Name)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}

		secrets = append(secrets, secret)
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(err, op)
	}

	return secrets, nil
}

func (r *secretRepository) AddSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) (uuid.UUID, error) {
	const op = "add secret"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}
	defer conn.Release()

	var txOptions pgx.TxOptions
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const sql = `INSERT INTO secrets (user_id, key_id, name, data) VALUES ($1, $2, $3, $4) RETURNING secret_id;`
	row := tx.QueryRow(ctx, sql, userID, secret.KeyID, secret.Name, secret.Data)
	var id uuid.UUID
	err = row.Scan(&id)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return id, nil
}

func (r *secretRepository) UpdateSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) error {
	const op = "update secret"

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

	const sql = `UPDATE secrets SET name=$1, data=$2 WHERE secret_id=$3;`
	tag, err := tx.Exec(ctx, sql, secret.Name, secret.Data, secret.ID)
	if tag.RowsAffected() == 0 {
		return vault.ErrSecretNotFound
	}
	if err != nil {
		return errors.Wrap(err, op)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *secretRepository) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	const op = "get secret"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer conn.Release()

	secret := &model.Secret{ID: secretID}
	const sql = `SELECT key_id, name, data FROM secrets WHERE secret_id=$1 AND user_id=$2`
	row := conn.QueryRow(ctx, sql, secretID, userID)
	err = row.Scan(&secret.KeyID, &secret.Name, &secret.Data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, vault.ErrSecretNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return secret, nil
}

func (r *secretRepository) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	const op = "delete secret"

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

	const sql = `DELETE FROM secrets WHERE secret_id=$1;`
	_, err = tx.Exec(ctx, sql, secretID)
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
