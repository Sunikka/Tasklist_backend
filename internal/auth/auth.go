package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sunikka/tasklist-backendGo/internal/db"
	"github.com/sunikka/tasklist-backendGo/internal/utils"
)

type AuthHandler func(http.ResponseWriter, *http.Request, utils.User)

var jwtKey = os.Getenv("JWT_KEY")

func GenerateToken(userID uuid.UUID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtKey), nil
	})
}

// JWT auth middleware
func MiddlewareJWT(handlerFunc http.HandlerFunc, s db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := GetTokenString(r)
		if err != nil {
			utils.ResponsePermDenied(w)
			return
		}

		token, err := validateJWT(tokenStr)
		if err != nil {
			utils.ResponsePermDenied(w)
			return
		}

		userID, err := utils.GetUserID(r)
		if err != nil {
			utils.ResponsePermDenied(w)
			return
		}

		user, err := s.GetUserById(userID)
		if err != nil {
			utils.ResponsePermDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if user.ID.String() != claims["user_id"].(string) {
			utils.ResponsePermDenied(w)
			return
		}

		handlerFunc(w, r)

	}
}

func GetTokenString(r *http.Request) (string, error) {
	authHeaderContent := r.Header.Get("Authorization")

	if authHeaderContent == "" {
		return "", errors.New("authentication failed")
	}

	values := strings.Split(authHeaderContent, " ")
	if len(values) != 2 {
		return "", errors.New("authentication failed")
	}

	if values[0] != "JWT" {
		return "", errors.New("authentication failed")
	}

	return values[1], nil

}
