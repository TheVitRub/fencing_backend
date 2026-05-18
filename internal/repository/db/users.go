package db

import (
	"context"
	"fmt"

	dbEntities "fencing-club/internal/domain/entities/db"
)

func (r *DBRepository) GetUserByLogin(ctx context.Context, login string) (*dbEntities.User, error) {
	var u dbEntities.User
	err := r.db.Get(ctx, &u,
		`SELECT id, login, email, password_hash, display_name, role, created_at, updated_at
		 FROM users
		 WHERE login=$1 OR email=$1`,
		login,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserByLogin login=%s: %w", login, err)
	}
	return &u, nil
}

func (r *DBRepository) GetUserByID(ctx context.Context, id int64) (*dbEntities.User, error) {
	var u dbEntities.User
	err := r.db.Get(ctx, &u,
		`SELECT id, login, email, password_hash, display_name, role, created_at, updated_at
		 FROM users
		 WHERE id=$1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID id=%d: %w", id, err)
	}
	return &u, nil
}

func (r *DBRepository) GetUserByIdentity(ctx context.Context, provider, providerUserID string) (*dbEntities.User, error) {
	var u dbEntities.User
	err := r.db.Get(ctx, &u,
		`SELECT u.id, u.login, u.email, u.password_hash, u.display_name, u.role, u.created_at, u.updated_at
		 FROM users u
		 JOIN user_identities ui ON ui.user_id=u.id
		 WHERE ui.provider=$1 AND ui.provider_user_id=$2`,
		provider, providerUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserByIdentity provider=%s: %w", provider, err)
	}
	return &u, nil
}

func (r *DBRepository) ListUsers(ctx context.Context) ([]dbEntities.User, error) {
	var users []dbEntities.User
	err := r.db.Select(ctx, nil, &users,
		`SELECT id, login, email, password_hash, display_name, role, created_at, updated_at
		 FROM users
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListUsers: %w", err)
	}
	return users, nil
}

func (r *DBRepository) CreateUser(ctx context.Context, u *dbEntities.User) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&u.ID},
		`INSERT INTO users (login, email, password_hash, display_name, role)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		u.Login, u.Email, u.PasswordHash, u.DisplayName, u.Role,
	)
	if err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	return nil
}

func (r *DBRepository) CreateUserIdentity(ctx context.Context, userID int64, provider, providerUserID string) error {
	_, err := r.db.Exec(ctx, nil,
		`INSERT INTO user_identities (user_id, provider, provider_user_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (provider, provider_user_id) DO NOTHING`,
		userID, provider, providerUserID,
	)
	if err != nil {
		return fmt.Errorf("CreateUserIdentity: %w", err)
	}
	return nil
}

func (r *DBRepository) UpdateUserRole(ctx context.Context, id int64, role string) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`,
		role, id,
	)
	if err != nil {
		return fmt.Errorf("UpdateUserRole id=%d: %w", id, err)
	}
	return nil
}
