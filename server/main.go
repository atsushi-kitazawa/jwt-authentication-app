package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := openDatabase("app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	userRepository := newUserRepository(db)
	authenticator := newAuthenticator(jwtSecret())

	registerRoutes(mux, userRepository, authenticator)

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func jwtSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}

	return "dev-secret-change-me"
}
