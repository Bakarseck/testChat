package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var Sessions = make(map[string]string)

func GenerateSessionID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func GetUsernameFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	sessionToken := cookie.Value

	username, ok := Sessions[sessionToken]
	if !ok {
		return "", fmt.Errorf("session invalide ou expirée")
	}

	return username, nil
}

