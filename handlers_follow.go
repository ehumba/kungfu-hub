package main

import (
	"encoding/json"
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) followHandler(w http.ResponseWriter, r *http.Request) {
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

	var params struct {
		FollowedID string `json:"followed_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.FollowedID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	followedUUID, err := uuid.Parse(params.FollowedID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid followed_id")
		return
	}

	followedNullID := uuid.NullUUID{UUID: followedUUID, Valid: true}
	userNullID := uuid.NullUUID{UUID: userID, Valid: true}

	if userNullID.UUID == followedNullID.UUID {
		respondWithError(w, http.StatusBadRequest, "cannot follow yourself")
		return
	}

	newFollow, err := a.dbQueries.CreateFollow(r.Context(), database.CreateFollowParams{
		Column1: userNullID,
		Column2: followedNullID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create follow")
		return
	}

	respondWithJSON(w, http.StatusOK, newFollow)
}

func (a *apiConfig) unfollowHandler(w http.ResponseWriter, r *http.Request) {
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

	var params struct {
		FollowedID string `json:"followed_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.FollowedID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	followedUUID, err := uuid.Parse(params.FollowedID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid followed_id")
		return
	}

	followedNullID := uuid.NullUUID{UUID: followedUUID, Valid: true}
	userNullID := uuid.NullUUID{UUID: userID, Valid: true}

	err = a.dbQueries.DeleteFollow(r.Context(), database.DeleteFollowParams{
		Column1: userNullID,
		Column2: followedNullID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete follow")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "unfollowed successfully"})
}
