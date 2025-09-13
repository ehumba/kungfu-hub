package main

import (
	"encoding/json"
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/ehumba/kungfu-hub/internal/database"
)

func (a *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type reqParams struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid password")
		return
	}

	newUserParams := database.CreateUserParams{
		Username:     params.UserName,
		Email:        params.Email,
		PasswordHash: hashedPassword,
	}
	newUserDB, err := a.dbQueries.CreateUser(r.Context(), newUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	newUser := User{
		ID:        newUserDB.ID,
		Username:  newUserDB.Username,
		Email:     newUserDB.Email,
		CreatedAt: newUserDB.CreatedAt,
		UpdatedAt: newUserDB.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}
