package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"backend/domain"
	"backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	accountRepo repository.AccountRepository
	secretKey   string
}

func NewAuthService(accountRepo repository.AccountRepository, secretKey string) *AuthService {
	return &AuthService{
		accountRepo: accountRepo,
		secretKey:   secretKey,
	}
}

func (s *AuthService) RegisterAccount(ctx context.Context, account domain.Account) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(hashedPassword)

	return s.accountRepo.CreateAccount(ctx, account)
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	account, err := s.accountRepo.GetAccountByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if account == nil {
		return "", errors.New("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid login or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    account.ID,
		"login": account.Login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
