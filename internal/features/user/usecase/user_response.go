package userusecase

import "github.com/codepnw/stdlib-ticket-system/internal/features/user"

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *userUsecase) generateToken(usr user.User) (Response, error) {
	accessToken, err := u.token.GenerateAccessToken(usr)
	if err != nil {
		return Response{}, err
	}

	refreshToken, err := u.token.GenerateRefreshToken(usr)
	if err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
