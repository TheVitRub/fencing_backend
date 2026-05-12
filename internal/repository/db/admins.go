package db

import (
	"context"
	"fmt"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// GetAdminByLogin возвращает администратора по логину.
// При отсутствии возвращает sql.ErrNoRows.
func (r *DBRepository) GetAdminByLogin(ctx context.Context, login string) (*dbEntities.Admin, error) {
	var a dbEntities.Admin
	err := r.db.Get(ctx, &a,
		`SELECT id, login, password_hash FROM admins WHERE login = $1`,
		login,
	)
	if err != nil {
		return nil, fmt.Errorf("GetAdminByLogin: %w", err)
	}
	return &a, nil
}
