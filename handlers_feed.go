package main

import (
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/google/uuid"
)

func (a *apiConfig) fetchFeedHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "no authorization header")
		return
	}

	userID, err := auth.ValidateJWT(token, a.secret)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	userNullID := uuid.NullUUID{UUID: userID, Valid: true}

	feedItems, err := a.dbQueries.FetchFeed(r.Context(), userNullID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch feed")
		return
	}

	respondWithJSON(w, http.StatusOK, feedItems)
}
