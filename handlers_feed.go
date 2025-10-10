package main

import (
	"fmt"
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/ehumba/kungfu-hub/internal/feed"
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

func (cfg *apiConfig) importHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	importer := feed.NewImporter(*cfg.dbQueries)

	if err := importer.ImportIJFNews(ctx, 1); err != nil {
		http.Error(w, fmt.Sprintf("import failed: %v", err), http.StatusInternalServerError)
		return
	}
	if err := importer.ImportWTNews(ctx, 2); err != nil {
		http.Error(w, fmt.Sprintf("import failed: %v", err), http.StatusInternalServerError)
		return
	}
	if err := importer.ImportIBJJFNews(ctx, 3); err != nil {
		http.Error(w, fmt.Sprintf("import failed: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Import completed successfully"))
}
