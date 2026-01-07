package userusecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/user"
	userrepo "github.com/codepnw/stdlib-ticket-system/internal/features/user/repo"
	"github.com/codepnw/stdlib-ticket-system/internal/helper"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
	jwttoken "github.com/codepnw/stdlib-ticket-system/pkg/jwt"
)

type UserUsecase interface {
	Register(ctx context.Context, input user.User) (Response, error)
	Login(ctx context.Context, input user.User) (Response, error)
}

type userUsecase struct {
	tx    database.TxManager
	token jwttoken.JWTToken
	repo  userrepo.UserRepository
}

func NewUserUsecase(tx database.TxManager, token jwttoken.JWTToken, repo userrepo.UserRepository) UserUsecase {
	return &userUsecase{
		tx:    tx,
		token: token,
		repo:  repo,
	}
}

func (u *userUsecase) Register(ctx context.Context, input user.User) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	hashedPassword, err := helper.HashPassword(input.HashPassword)
	if err != nil {
		return Response{}, err
	}
	input.HashPassword = hashedPassword

	// Transaction
	var response Response
	err = u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		// Create User
		created, err := u.repo.CreateUser(ctx, input)
		if err != nil {
			return err
		}

		// Generate Token
		resp, err := u.generateToken(created)
		if err != nil {
			return err
		}

		// Save Refresh Token
		if err := u.repo.SaveRefreshToken(ctx, user.Auth{
			UserID:       created.ID,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    time.Now().Add(config.RefreshTokenDuration),
		}); err != nil {
			return err
		}

		response = resp
		return nil
	})
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func (u *userUsecase) Login(ctx context.Context, input user.User) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	foundUser, err := u.repo.FindUsername(ctx, input.Username)
	if err != nil {
		return Response{}, err
	}

	ok := helper.ComparePassword(input.HashPassword, foundUser.HashPassword)
	if !ok {
		return Response{}, errs.ErrInvalidCredentials
	}

	// Transaction
	var response Response
	err = u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		// Generate Token
		resp, err := u.generateToken(foundUser)
		if err != nil {
			return err
		}
		
		// Save Refresh Token
		if err := u.repo.SaveRefreshToken(ctx, user.Auth{
			UserID: foundUser.ID,
			RefreshToken: resp.RefreshToken,
			ExpiresAt: time.Now().Add(config.RefreshTokenDuration),
		}); err != nil {
			return err
		}
		
		response = resp
		return nil
	})
	
	if err != nil {
		return Response{}, err
	}
	return response, nil
}
