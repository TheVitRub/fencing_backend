package application

import (
	"context"
	"strings"

	"fencing-club/internal/domain/apperrors"
	"fencing-club/internal/domain/dto"
	dbEntities "fencing-club/internal/domain/entities/db"

	"golang.org/x/crypto/bcrypt"
)

var roles = map[string]bool{
	"registered": true,
	"student":    true,
	"instructor": true,
	"admin":      true,
	"founder":    true,
}

func userToDTO(u *dbEntities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          u.ID,
		Login:       u.Login,
		Email:       u.Email,
		DisplayName: displayName(u),
		Role:        u.Role,
	}
}

func displayName(u *dbEntities.User) string {
	if strings.TrimSpace(u.DisplayName) != "" {
		return u.DisplayName
	}
	return u.Login
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, apperrors.Unauthorizedf("неверный логин или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.Unauthorizedf("неверный логин или пароль")
	}

	token, err := s.issueToken(user.ID, user.Login, user.Role)
	if err != nil {
		return nil, apperrors.New(apperrors.KindInternal, "не удалось выпустить токен", err)
	}

	return &dto.LoginResponse{Token: token, User: userToDTO(user)}, nil
}

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	login := strings.TrimSpace(req.Login)
	if login == "" {
		return nil, apperrors.Validationf("логин обязателен")
	}
	if len(req.Password) < 6 {
		return nil, apperrors.Validationf("пароль должен быть не короче 6 символов")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.New(apperrors.KindInternal, "не удалось создать пароль", err)
	}
	user := &dbEntities.User{
		Login:        login,
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: string(hash),
		DisplayName:  strings.TrimSpace(req.DisplayName),
		Role:         "registered",
	}
	if user.DisplayName == "" {
		user.DisplayName = login
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, wrapRepoError("Register", err)
	}
	token, err := s.issueToken(user.ID, user.Login, user.Role)
	if err != nil {
		return nil, apperrors.New(apperrors.KindInternal, "не удалось выпустить токен", err)
	}
	return &dto.LoginResponse{Token: token, User: userToDTO(user)}, nil
}

func (s *Service) Me(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, wrapRepoError("Me", err)
	}
	resp := userToDTO(user)
	return &resp, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, wrapRepoError("ListUsers", err)
	}
	resp := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		u := user
		resp = append(resp, userToDTO(&u))
	}
	return resp, nil
}

func (s *Service) UpdateUserRole(ctx context.Context, id int64, role string) error {
	if id <= 0 {
		return apperrors.Validationf("id пользователя должен быть больше нуля")
	}
	if !roles[role] {
		return apperrors.Validationf("неизвестная роль")
	}
	return wrapRepoError("UpdateUserRole", s.repo.UpdateUserRole(ctx, id, role))
}

func IsAdminRole(role string) bool {
	return role == "admin" || role == "founder"
}

func IsInstructorRole(role string) bool {
	return role == "instructor" || role == "admin" || role == "founder"
}

func IsStudentRole(role string) bool {
	return role == "student" || IsInstructorRole(role)
}
