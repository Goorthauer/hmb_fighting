package usecase

import (
	"fmt"
	"hmb_fighting/server/dtos"
	"hmb_fighting/server/entities"
	"hmb_fighting/server/jwt"
	"hmb_fighting/server/types"
	"hmb_fighting/server/utils"
)

func (u *Usecase) RegisterUser(currentUser entities.User) (*dtos.RegisterUserResp, error) {
	user, err := u.db.GetUserByEmail(currentUser.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %v", err)
	}

	if user.Email != "" {
		return nil, fmt.Errorf("User exist")
	}

	user = currentUser
	hashPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("Failed hash password: %v", err)
	}
	user.Password = hashPass
	user.ID = generateClientID()

	tokenPair, err := jwt.GenerateTokenPair(user, types.UserRoleSpectator)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate tokens: %v", err)
	}

	err = u.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		return nil, fmt.Errorf("Failed to save user with refresh token: %v", err)
	}

	return &dtos.RegisterUserResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}, nil
}

func (u *Usecase) LoginUser(email, password string) (*dtos.RegisterUserResp, error) {
	user, err := u.db.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %v", err)
	}

	if user.Email == "" {
		return nil, fmt.Errorf("user not exist")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	tokenPair, err := jwt.GenerateTokenPair(user, types.UserRoleSpectator)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate tokens: %v", err)
	}

	err = u.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		return nil, fmt.Errorf("Failed to save user with refresh token: %v", err)
	}

	return &dtos.RegisterUserResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}, nil
}

func (u *Usecase) RefreshToken(refreshToken string) (*dtos.RegisterUserResp, error) {
	user, err := u.db.GetUserByRefresh(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid refresh token: %v", err)
	}
	if user.ID == "" {
		return nil, fmt.Errorf("Invalid refresh token")
	}

	tokenPair, err := jwt.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid refresh token: %v", err)
	}

	err = u.db.SetUser(tokenPair.RefreshToken, user)
	if err != nil {
		return nil, fmt.Errorf("Failed to update refresh token: %v", err)
	}

	return &dtos.RegisterUserResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ClientID:     user.ID,
	}, nil
}

func (u *Usecase) CheckClient(clientID, accessToken string) (bool, error) {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return false, fmt.Errorf("Invalid token: %v", err)
	}

	user, err := u.db.GetUserByEmail(claims.Email)
	if err != nil {
		return false, fmt.Errorf("Failed to get user: %v", err)
	}

	return user.ID == clientID && user.ID == claims.ClientID, nil
}
