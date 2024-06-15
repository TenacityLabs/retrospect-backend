package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TenacityLabs/time-capsule-backend/config"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "userID"

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

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" && strings.HasPrefix(tokenAuth, "Bearer ") {
		return strings.TrimPrefix(tokenAuth, "Bearer ")
	}
	return ""
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
}

func WithJWTAuth(handlerFunc http.HandlerFunc, userStore types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from user request
		tokenString := getTokenFromRequest(r)

		// validate jwt
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			PermissionDenied(w)
			return
		}
		if !token.Valid {
			log.Println("token is invalid")
			PermissionDenied(w)
			return
		}

		// fetch user id from db if authenticated
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("error converting userID to int: %v", err)
			PermissionDenied(w)
			return
		}

		user, err := userStore.GetUserById(uint(userID))
		if err != nil {
			log.Printf("error fetching user: %v", err)
			PermissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func GetUserIdFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value(UserKey).(uint)
	if !ok {
		return 0
	}

	return userID
}
