package pgsql

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/migration"
)

type userRepository struct {
	pool       *pgxpool.Pool
	connString string
}

func NewUserRepository(ctx context.Context, connString string) (*userRepository, error) {
	const op = "new user repository"

	migrator := migration.NewDatabaseMigrator(connString)
	if err := migrator.Up(); err != nil {
		return nil, errors.Wrapf(err, op)
	}

	var err error
	pool, err := initPool(ctx, connString)
	if err != nil {
		return nil, errors.Wrapf(err, op)
	}

	r := &userRepository{
		pool,
		connString,
	}
	return r, nil
}

func (r *userRepository) Close() {
	if r.pool == nil {
		return
	}
	r.pool.Close()
}

func (r *userRepository) Register(ctx context.Context, user model.User) (model.User, error) {
	const op = "register"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return model.User{}, errors.Wrap(err, op)
	}
	defer conn.Release()

	var txOptions pgx.TxOptions
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		return model.User{}, errors.Wrap(err, op)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const sql = `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id;`
	row := conn.QueryRow(ctx, sql, user.Email, user.Password)
	err = row.Scan(&user.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation &&
			pgErr.ConstraintName == "users_email_key" {
			return model.User{}, auth.ErrUserWithEmailIsRegistered
		}
		if pgErr.Code == pgerrcode.CheckViolation &&
			pgErr.ConstraintName == "users_password_check" {
			return model.User{}, auth.ErrUserPasswordIsEmpty
		}
	}
	if err != nil {
		return model.User{}, errors.Wrap(err, op)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.User{}, errors.Wrap(err, op)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	const op = "find by email"

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return model.User{}, errors.Wrap(err, op)
	}
	defer conn.Release()

	var user model.User
	const sql = "SELECT user_id, email, password FROM users WHERE email=$1"
	row := conn.QueryRow(ctx, sql, email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, auth.ErrUserIsNotRegistered
	}

	return user, nil
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
