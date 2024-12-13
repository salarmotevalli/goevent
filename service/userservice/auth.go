package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"event-manager/entity"
	"fmt"
	"log"
)

type UserRepo interface {
	GetUserByUsername(string) (entity.User, bool, error)
	CreateUser(entity.User) (entity.User, error)
}

type AuthService interface {
	CreateAccessToken(entity.User) (string, error)
	CreateRefreshToken(entity.User) (string, error)
}

type UserService struct {
	repo UserRepo
	auth AuthService
}

func New(r UserRepo, as AuthService) UserService {
	return UserService{
		repo: r,
		auth: as,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func (s *UserService) Register(req RegisterRequest) (RegisterResponse, error) {
	// check is there username in db
	_, exist, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		log.Println(err)
		return RegisterResponse{}, errors.New("unexpected")
	}

	if exist {
		return RegisterResponse{}, errors.New("user already exists")
	}

	// hash the password
	hashedPassword := getMD5Hash(req.Password)

	// create new user in db
	user := entity.User{}
	user.SetUsername(req.Username)
	user.SetPassword(hashedPassword)

	user, err = s.repo.CreateUser(user)

	return RegisterResponse{
		User: user,
	}, err
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *UserService) Login(req LoginRequest) (LoginResponse, error) {
	// check is there username and password in db
	user, exist, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return LoginResponse{}, errors.New("unexpected")
	}

	if !exist {
		return LoginResponse{}, errors.New("user not found")
	}

	if user.Password() != getMD5Hash(req.Password) {
		return LoginResponse{}, errors.New("password is incorrect")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
