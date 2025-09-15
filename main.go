package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	dbQueries *database.Queries
	secret    string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL_LOCAL")
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("could not load database: %v", err)
		return
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := apiConfig{
		dbQueries: dbQueries,
		secret:    secret,
	}

	// user creation
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	// login
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	// refresh token
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	// revoke token
	mux.HandleFunc("DELETE /api/revoke", apiCfg.handlerRevoke)

	// update user data
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserDataHandler)

	// delete user
	mux.HandleFunc("DELETE /api/users", apiCfg.deleteUserHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
