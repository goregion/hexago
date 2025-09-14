package jwt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret            []byte
	blocklistFilePath string
}

func NewTokenManager(secret string, blocklistFilePath string) *TokenManager {
	return &TokenManager{
		secret:            []byte(secret),
		blocklistFilePath: blocklistFilePath,
	}
}

func (tm *TokenManager) isBlocked(token string) (bool, error) {
	file, err := os.Open(tm.blocklistFilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			if line == token {
				return true, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func (tm *TokenManager) ParseToken(tokenString string) (client string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secret, nil
	})
	if err != nil {
		log.Fatal("Error parsing token:", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		isBlocked, err := tm.isBlocked(tokenString)
		if err != nil {
			return "", fmt.Errorf("failed to check token blocklist: %w", err)
		}
		if isBlocked {
			return "", fmt.Errorf("token is blocked")
		}

		client, ok := claims["client"].(string)
		if !ok {
			return "", fmt.Errorf("client claim is missing or invalid")
		}
		return client, nil
	}
	return "", fmt.Errorf("invalid token")
}

func (tm *TokenManager) GenerateToken(client string) (string, error) {
	claims := jwt.MapClaims{
		"client": client,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tm.secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}
