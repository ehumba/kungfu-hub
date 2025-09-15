package main

import (
	"encoding/json"
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/ehumba/kungfu-hub/internal/database"
)

func (a *apiConfig) updateUserDataHandler(w http.ResponseWriter, r *http.Request) {
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

	defer r.Body.Close()
	type reqParams struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if params.UserName == "" || params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid password")
		return
	}

	updateUserParams := database.UpdateUserDataParams{
		ID:           userID,
		Username:     params.UserName,
		Email:        params.Email,
		PasswordHash: hashedPassword,
	}
	err = a.dbQueries.UpdateUserData(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user data")
		return
	}

	updatedUserDB, err := a.dbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve updated user data")
		return
	}

	updatedUser := User{
		ID:        updatedUserDB.ID,
		Username:  updatedUserDB.Username,
		Email:     updatedUserDB.Email,
		CreatedAt: updatedUserDB.CreatedAt,
		UpdatedAt: updatedUserDB.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, updatedUser)

}

func (a *apiConfig) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err = a.dbQueries.DeleteUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete user")
		return
	}

	w.WriteHeader(204)
}
