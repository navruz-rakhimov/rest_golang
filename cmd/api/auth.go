package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/navruz-rakhimov/sarkortelecom/internal/data"
	"time"
)

var jwtKey = []byte("secret_key")

func (app *application) GenerateJwtToken(user data.User) (string, error) {

	payload := jwt.MapClaims{
		"sub":     user.Login,
		"login":   user.Login,
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("failed to sign JWT token")
	}

	return accessToken, nil
}

func (app *application) IsAccessTokenValid(signedToken string) (bool, error) {
	var mapClaims jwt.MapClaims

	token, err := jwt.ParseWithClaims(
		signedToken,
		&mapClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("гnexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		},
	)
	if err != nil {
		return false, err
	}

	if err = token.Claims.Valid(); err != nil {
		return false, nil
	}

	return true, nil
}

func (app *application) GetCurrentUserId(signedToken string) (int, error) {
	var mapClaims jwt.MapClaims

	token, err := jwt.ParseWithClaims(
		signedToken,
		&mapClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		},
	)
	if err != nil {
		return -1, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return -1, errors.New("couldn't parse claims")
	}

	userId := int((*claims)["user_id"].(float64))
	return userId, nil
}
