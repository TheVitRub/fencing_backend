package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"fencing-club/internal/domain/apperrors"
	"fencing-club/internal/domain/dto"
	dbEntities "fencing-club/internal/domain/entities/db"

	"golang.org/x/crypto/bcrypt"
)

type oauthProfile struct {
	Provider       string
	ProviderUserID string
	Email          string
	DisplayName    string
	Login          string
}

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

func (s *Service) OAuthStartURL(provider, redirectURI string) (string, error) {
	switch provider {
	case "vk":
		clientID := os.Getenv("VK_CLIENT_ID")
		if clientID == "" {
			return "", apperrors.Validationf("VK OAuth не настроен: нужен VK_CLIENT_ID")
		}
		q := url.Values{}
		q.Set("client_id", clientID)
		q.Set("display", "page")
		q.Set("redirect_uri", redirectURI)
		q.Set("scope", "email")
		q.Set("response_type", "code")
		q.Set("v", "5.131")
		return "https://oauth.vk.com/authorize?" + q.Encode(), nil
	case "google":
		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		if clientID == "" {
			return "", apperrors.Validationf("Google OAuth не настроен: нужен GOOGLE_CLIENT_ID")
		}
		q := url.Values{}
		q.Set("client_id", clientID)
		q.Set("redirect_uri", redirectURI)
		q.Set("response_type", "code")
		q.Set("scope", "openid email profile")
		return "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode(), nil
	default:
		return "", apperrors.Validationf("неизвестный OAuth-провайдер")
	}
}

func (s *Service) OAuthCallback(ctx context.Context, provider, code, redirectURI string) (*dto.LoginResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, apperrors.Validationf("OAuth code не передан")
	}

	var profile oauthProfile
	var err error
	switch provider {
	case "vk":
		profile, err = fetchVKProfile(code, redirectURI)
	case "google":
		profile, err = fetchGoogleProfile(code, redirectURI)
	default:
		return nil, apperrors.Validationf("неизвестный OAuth-провайдер")
	}
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByIdentity(ctx, profile.Provider, profile.ProviderUserID)
	if err != nil {
		if profile.Email != "" {
			if byEmail, emailErr := s.repo.GetUserByLogin(ctx, profile.Email); emailErr == nil {
				user = byEmail
				_ = s.repo.CreateUserIdentity(ctx, user.ID, profile.Provider, profile.ProviderUserID)
			}
		}
		if user == nil {
			user = &dbEntities.User{
				Login:       profile.Login,
				Email:       profile.Email,
				DisplayName: profile.DisplayName,
				Role:        "registered",
			}
			if user.DisplayName == "" {
				user.DisplayName = user.Login
			}
			if err := s.repo.CreateUser(ctx, user); err != nil {
				return nil, wrapRepoError("OAuthCreateUser", err)
			}
			if err := s.repo.CreateUserIdentity(ctx, user.ID, profile.Provider, profile.ProviderUserID); err != nil {
				return nil, wrapRepoError("OAuthCreateIdentity", err)
			}
		}
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

func fetchVKProfile(code, redirectURI string) (oauthProfile, error) {
	clientID := os.Getenv("VK_CLIENT_ID")
	clientSecret := os.Getenv("VK_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return oauthProfile{}, apperrors.Validationf("VK OAuth не настроен: нужны VK_CLIENT_ID и VK_CLIENT_SECRET")
	}
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("redirect_uri", redirectURI)
	q.Set("code", code)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		UserID      int64  `json:"user_id"`
		Email       string `json:"email"`
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := getJSON("https://oauth.vk.com/access_token?"+q.Encode(), &tokenResp); err != nil {
		return oauthProfile{}, err
	}
	if tokenResp.Error != "" {
		return oauthProfile{}, apperrors.Validationf("VK OAuth: %s", tokenResp.Description)
	}

	displayName := fmt.Sprintf("VK %d", tokenResp.UserID)
	if tokenResp.AccessToken != "" {
		uq := url.Values{}
		uq.Set("access_token", tokenResp.AccessToken)
		uq.Set("user_ids", fmt.Sprintf("%d", tokenResp.UserID))
		uq.Set("fields", "screen_name")
		uq.Set("v", "5.131")
		var usersResp struct {
			Response []struct {
				ID         int64  `json:"id"`
				FirstName  string `json:"first_name"`
				LastName   string `json:"last_name"`
				ScreenName string `json:"screen_name"`
			} `json:"response"`
		}
		if err := getJSON("https://api.vk.com/method/users.get?"+uq.Encode(), &usersResp); err == nil && len(usersResp.Response) > 0 {
			u := usersResp.Response[0]
			displayName = strings.TrimSpace(u.FirstName + " " + u.LastName)
			if displayName == "" {
				displayName = u.ScreenName
			}
		}
	}

	providerID := fmt.Sprintf("%d", tokenResp.UserID)
	return oauthProfile{
		Provider: "vk", ProviderUserID: providerID, Email: tokenResp.Email,
		DisplayName: displayName, Login: "vk_" + providerID,
	}, nil
}

func fetchGoogleProfile(code, redirectURI string) (oauthProfile, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return oauthProfile{}, apperrors.Validationf("Google OAuth не настроен: нужны GOOGLE_CLIENT_ID и GOOGLE_CLIENT_SECRET")
	}
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return oauthProfile{}, apperrors.New(apperrors.KindInternal, "не удалось создать OAuth-запрос", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := doJSON(req, &tokenResp); err != nil {
		return oauthProfile{}, err
	}
	if tokenResp.Error != "" {
		return oauthProfile{}, apperrors.Validationf("Google OAuth: %s", tokenResp.Description)
	}

	req, err = http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return oauthProfile{}, apperrors.New(apperrors.KindInternal, "не удалось создать OAuth-запрос", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	var info struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := doJSON(req, &info); err != nil {
		return oauthProfile{}, err
	}
	return oauthProfile{
		Provider: "google", ProviderUserID: info.Sub, Email: info.Email,
		DisplayName: info.Name, Login: "google_" + info.Sub,
	}, nil
}

func getJSON(rawURL string, dest any) error {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return apperrors.New(apperrors.KindInternal, "не удалось создать OAuth-запрос", err)
	}
	return doJSON(req, dest)
}

func doJSON(req *http.Request, dest any) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return apperrors.New(apperrors.KindInternal, "OAuth-провайдер недоступен", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return apperrors.New(apperrors.KindInternal, "не удалось прочитать OAuth-ответ", err)
	}
	if resp.StatusCode >= 400 {
		return apperrors.Validationf("OAuth-провайдер вернул статус %d", resp.StatusCode)
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return apperrors.New(apperrors.KindInternal, "не удалось разобрать OAuth-ответ", err)
	}
	return nil
}
