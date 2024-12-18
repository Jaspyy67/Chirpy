package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourusername/chirpy/internal/auth"
	"github.com/yourusername/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	parts := strings.Split(authorizationHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := parts[1]
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var updateUserRequest UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	params := database.UpdateUserEmailAndPasswordParams{
		ID:             userID,
		Email:          updateUserRequest.Email,
		HashedPassword: string(hashedPassword),
	}

	err = cfg.db.UpdateUserEmailAndPassword(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// After successful update
	updatedUser, err := cfg.db.GetUserByID(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a response struct or map
	response := struct {
		Email string `json:"email"`
	}{
		Email: updatedUser.Email,
	}

	// Convert response to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)
}
