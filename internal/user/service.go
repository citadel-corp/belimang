package user

import (
	"context"
	"time"

	"github.com/citadel-corp/belimang/internal/common/id"
	"github.com/citadel-corp/belimang/internal/common/jwt"
	"github.com/citadel-corp/belimang/internal/common/password"
)

type Service interface {
	Create(ctx context.Context, req CreateUserPayload) (*UserAuthResponse, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) Create(ctx context.Context, req CreateUserPayload) (*UserAuthResponse, error) {
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user := &Users{
		UID:            id.GenerateStringID(16),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		UserType:       req.UserType,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Minute*2, user.UID, string(user.UserType))
	if err != nil {
		return nil, err
	}
	return &UserAuthResponse{
		AccessToken: accessToken,
	}, nil
}
