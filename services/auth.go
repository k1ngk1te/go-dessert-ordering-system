package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	appConstants "dessert-ordering-go-system/internal/app_constants"
	models "dessert-ordering-go-system/models"
)

type AuthService struct {
	UserModel *models.UserModel
	JWTSecret []byte
}

type AuthData struct {
	User  models.UserData `json:"user"`
	Token string          `json:"token"`
}

type UserClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(userModel *models.UserModel, jwtSecret string) *AuthService {
	return &AuthService{
		UserModel: userModel,
		JWTSecret: []byte(jwtSecret),
	}
}

func (a *AuthService) Authenticate(contact, password string) (*models.UserData, error) {
	var (
		err      error
		userData *models.UserData
	)

	// 1. Check if the contact is an email or username
	if strings.Contains(contact, "@") {
		userData, err = a.UserModel.AuthenticateByEmail(contact, password)
	} else {
		userData, err = a.UserModel.AuthenticateByUsername(contact, password)
	}

	return userData, err
}

func (a *AuthService) RegisterUser(username, email, password string) error {
	return a.UserModel.CreateUser(username, email, password)
}

func (a *AuthService) GetTokenExpiration() time.Duration {
	return appConstants.Jwt_Expiration
}

func (a *AuthService) GenerateAuthToken(userID int, username, email string) (string, error) {
	expirationTime := time.Now().Add(appConstants.Jwt_Expiration)

	claims := &UserClaims{
		ID:       userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (a *AuthService) CreateAuthData(userData models.UserData, token string) *AuthData {
	return &AuthData{
		User:  userData,
		Token: token,
	}
}
