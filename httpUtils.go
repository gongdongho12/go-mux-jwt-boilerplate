package main

import (
	"net/http"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"os"
)

type ErrorResponse struct {
	Type string `json:"type"`
	Message string `json:"message"`
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, response string){
	error := ErrorResponse{
		"error",
		response,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&error)
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}