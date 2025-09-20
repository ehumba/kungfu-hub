package main

import (
	"encoding/json"
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) listMartialArtsHandler(w http.ResponseWriter, r *http.Request) {
	martialArts, err := a.dbQueries.GetMartialArts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch martial arts")
		return
	}

	respondWithJSON(w, http.StatusOK, martialArts)
}

func (a *apiConfig) subscribeHandler(w http.ResponseWriter, r *http.Request) {
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
		MartialArtID string `json:"martial_art_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.MartialArtID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	parsedMartialID, err := uuid.Parse(params.MartialArtID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid martial_art_id")
		return
	}

	subscribeParams := database.AddSubscriptionParams{
		UserID:       userID,
		MartialArtID: parsedMartialID,
	}
	newSub, err := a.dbQueries.AddSubscription(r.Context(), subscribeParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create subscription")
		return
	}
	respondWithJSON(w, http.StatusCreated, newSub)
}

func (a *apiConfig) unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
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
		MartialArtID string `json:"martial_art_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.MartialArtID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	parsedMartialID, err := uuid.Parse(params.MartialArtID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid martial_art_id")
		return
	}

	removeSubParams := database.RemoveSubscriptionParams{
		UserID:       userID,
		MartialArtID: parsedMartialID,
	}
	err = a.dbQueries.RemoveSubscription(r.Context(), removeSubParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not remove subscription")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *apiConfig) listUserSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
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

	subscriptions, err := a.dbQueries.GetUserSubscriptions(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch subscriptions")
		return
	}

	respondWithJSON(w, http.StatusOK, subscriptions)
}
