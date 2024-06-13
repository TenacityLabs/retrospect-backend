package auth

import (
	"strconv"
	"time"

	"github.com/TenacityLabs/time-capsule-backend/config"
	"github.com/golang-jwt/jwt"
)

func CreateJWT(secret []byte, userID uint) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.FormatUint(uint64(userID), 10),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
