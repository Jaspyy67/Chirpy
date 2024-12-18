package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/yourusername/chirpy/internal/auth"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or malformed JWT")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired JWT")
		return
	}

	chirpID := r.PathValue("chirpID")
	fmt.Printf("Raw chirpID: %s\n", chirpID)

	chirpIDUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpIDUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.UserID.String() != userID.String() {
		respondWithError(w, http.StatusForbidden, "You can't delete someone else's chirp")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpIDUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp")
		return
	}
	w.WriteHeader(204)
}
