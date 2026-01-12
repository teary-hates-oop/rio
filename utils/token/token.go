package token

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user_id string) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "rio",
		Audience:  jwt.ClaimStrings{"rio-chat-client"},
		Subject:   user_id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(token_lifespan))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)

	if tokenString == "" {
		return errors.New("Token not found")
	}

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")

	if token != "" {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractTokenID(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			return sub, nil
		}
	}
	return "", errors.New("invalid token")
}
