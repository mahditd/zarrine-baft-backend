package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/utils"
)

type AuthService struct {
	userRepository repositories.UserRepository
	jwtSecret      string
	jwtExpireHours int
}

func NewAuthService(
	userRepository repositories.UserRepository,
	jwtSecret string,
	jwtExpireHours int,
) *AuthService {

	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

type RegisterInput struct {
	FullName        string
	Phone           string
	Email           string
	Password        string
	ConfirmPassword string
	CompanyName     string
	CompanyPhone    string
	Country         string
	Address         string
}

type LoginInput struct {
	Phone    string
	Password string
}

type LoginResult struct {
	User  *models.User
	Token string
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {

	input.Phone = utils.NormalizePhone(input.Phone)

	existingUser, err := s.userRepository.FindByPhone(input.Phone)

	if err == nil && existingUser != nil {
		return nil, errors.New("phone already exists")
	}

	if input.Email != "" {

		existingUser, err := s.userRepository.FindByEmail(input.Email)

		if err == nil && existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	if len(input.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	if input.Password != input.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:     input.FullName,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleCustomer,
	}

	if input.Email != "" {
		user.Email = &input.Email
	}

	if input.CompanyName != "" {
		user.CompanyName = &input.CompanyName
	}

	if input.CompanyPhone != "" {
		user.CompanyPhone = &input.CompanyPhone
	}

	if input.Country != "" {
		user.Country = &input.Country
	}

	if input.Address != "" {
		user.Address = &input.Address
	}

	err = s.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (*LoginResult, error) {

	input.Phone = utils.NormalizePhone(input.Phone)

	user, err := s.userRepository.FindByPhone(input.Phone)

	if err != nil {
		return nil, errors.New("invalid phone or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return nil, errors.New("invalid phone or password")
	}

	token, err := utils.GenerateToken(
		user.ID,
		string(user.Role),
		s.jwtSecret,
		s.jwtExpireHours,
	)

	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}
