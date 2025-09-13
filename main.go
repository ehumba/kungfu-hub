package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	dbQueries *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL_LOCAL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("could not load database: %v", err)
		return
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := apiConfig{
		dbQueries: dbQueries,
	}

	// user creation
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

}
