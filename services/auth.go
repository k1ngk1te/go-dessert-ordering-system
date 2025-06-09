package services

import (
	"dessert-ordering-go-system/models"
	"strings"
)

type AuthService struct {
	UserModel *models.UserModel
}

func NewAuthService(userModel *models.UserModel) *AuthService {
	return &AuthService{
		UserModel: userModel,
	}
}

func (a *AuthService) Authenticate(contact, password string) (*models.UserData, error) {
	var (
		err error
		userData *models.UserData
	)
	
	// 1. Check if the contact is an email or username
	if (strings.Contains(contact, "@")) {
		userData, err = a.UserModel.AuthenticateByEmail(contact, password)
	} else {
		userData, err = a.UserModel.AuthenticateByUsername(contact, password)
	}

	return userData, err
}

func (a *AuthService) RegisterUser(username, email, password string) error {
	return a.UserModel.CreateUser(username, email, password)
}