package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"dessert-ordering-go-system/models"
)

type AuthService struct {
	UserModel *models.UserModel
	JWTSecret []byte
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

func (a *AuthService) GenerateAuthToken(userID int, username, email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

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
