package repository

import (
	"avito/internal/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		log:  logger,
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	r.log.Info("creating user", slog.String("user_id", user.ID.String()))

	query := `insert into users(user_id, email, password_hash, user_role) values ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.Password, user.Role)
	if err != nil {
		r.log.Warn("creating user failed: " + err.Error())
		if isDuplicateError(err) {
			return domain.ErrUserExist
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	r.log.Info("getting user by id", slog.String("user_id", userID.String()))

	var user domain.User
	query := `select user_id, email, password_hash, user_role from users where user_id = $1`
	rows := r.pool.QueryRow(ctx, query, userID)
	err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		r.log.Warn("getting user by id error: " + err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, userEmail string) (domain.User, error) {
	r.log.Info("getting user by email", slog.String("email", userEmail))

	var user domain.User
	query := `select user_id, email, password_hash, user_role from users where email = $1`
	rows := r.pool.QueryRow(ctx, query, userEmail)
	err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		r.log.Warn("getting user by email error: " + err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, offset int, limit int) ([]domain.User, error) {
	r.log.Info("getting users", slog.Int("offset", offset), slog.Int("linit", limit))

	query := `select user_id, email, password_hash, user_role from users limit $1 offset $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		r.log.Warn("user repo: get all error: " + err.Error())
		return nil, fmt.Errorf("user repo: get all error: %v", err.Error())
	}
	defer rows.Close()

	var (
		users []domain.User
		user  domain.User
	)
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
		if err != nil {
			r.log.Warn("user get all error: scan user error")
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
