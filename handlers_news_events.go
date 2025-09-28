package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func (a *apiConfig) getEventsByMartialArtHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		MartialArtID string `json:"martial_art_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.MartialArtID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	id64, err := strconv.ParseInt(params.MartialArtID, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid martial_art_id")
		return
	}

	id32 := int32(id64)
	art_id := sql.NullInt32{Int32: id32, Valid: true}

	events, err := a.dbQueries.GetEventsByMartialArtID(r.Context(), art_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch events")
		return
	}
	respondWithJSON(w, http.StatusOK, events)
}

func (a *apiConfig) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := a.dbQueries.GetAllEvents(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch events")
		return
	}
	respondWithJSON(w, http.StatusOK, events)
}

func (a *apiConfig) getNewsByMartialArtHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		MartialArtID string `json:"martial_art_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.MartialArtID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	id64, err := strconv.ParseInt(params.MartialArtID, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid martial_art_id")
		return
	}

	id32 := int32(id64)
	art_id := sql.NullInt32{Int32: id32, Valid: true}

	news, err := a.dbQueries.GetNewsByMartialArt(r.Context(), art_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch news")
		return
	}
	respondWithJSON(w, http.StatusOK, news)
}

func (a *apiConfig) getNewsByEventHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		EventID string `json:"event_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.EventID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	id64, err := strconv.ParseInt(params.EventID, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid martial_art_id")
		return
	}

	id32 := int32(id64)
	event_id := sql.NullInt32{Int32: id32, Valid: true}

	news, err := a.dbQueries.GetNewsByEvent(r.Context(), event_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch news")
		return
	}
	respondWithJSON(w, http.StatusOK, news)
}

func (a *apiConfig) getAllNewsHandler(w http.ResponseWriter, r *http.Request) {
	news, err := a.dbQueries.GetAllNews(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch news")
		return
	}
	respondWithJSON(w, http.StatusOK, news)
}
