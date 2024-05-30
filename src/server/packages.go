package server

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func hashPassword(password string, login string) string {
	hash := sha512.Sum512([]byte(login + password))
	hashedPassword := hex.EncodeToString(hash[:])

	return hashedPassword
}

func checkAuth(w http.ResponseWriter, r *http.Request, signingKey string) jwt.MapClaims {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"reason": "error: user havent got auth token"}`))
		return nil
	}

	token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"reason": "invalid token"}`))
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"reason": "error while check token claims"}`))
		return nil
	}
	return claims
}
