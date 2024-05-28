package util

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, schema interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(schema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}
