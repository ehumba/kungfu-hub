package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
		log.Printf("JSON decode error: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
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

func (a *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	userDB, err := a.dbQueries.GetUserByEmail(r.Context(), params.Email)
	hashErr := auth.CheckPasswordHash(params.Password, userDB.PasswordHash)
	if err != nil || hashErr != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	authToken, err := auth.GenerateJWT(userDB.ID, a.secret, 30*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not generate authentication token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not generate refresh token")
		return
	}

	refreshParams := database.GenerateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userDB.ID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, err = a.dbQueries.GenerateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not save refresh token")
		return
	}

	user := User{
		ID:        userDB.ID,
		Username:  userDB.Username,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email:     userDB.Email,
	}

	resStruct := struct {
		User         `json:",inline"`
		AuthToken    string `json:"auth_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		User:         user,
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resStruct)
}
