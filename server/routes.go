package main

import "net/http"

func registerRoutes(mux *http.ServeMux, userRepository *userRepository, authenticator *authenticator) {
	requireAuth := authenticator.middleware

	mux.Handle("GET /", requireAuth(http.HandlerFunc(rootHandler)))
	mux.HandleFunc("GET /docs", docsHandler)
	mux.HandleFunc("GET /openapi.yaml", openAPISpecHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /users", createUserHandler(userRepository))
	mux.HandleFunc("POST /login", loginHandler(userRepository, authenticator))
	mux.Handle("GET /users", requireAuth(listUsersHandler(userRepository)))
	mux.Handle("GET /users/name/{name}", requireAuth(getUserByNameHandler(userRepository)))
	mux.Handle("GET /users/{id}", requireAuth(getUserHandler(userRepository)))
	mux.Handle("PUT /users/{id}", requireAuth(updateUserHandler(userRepository)))
	mux.Handle("DELETE /users/{id}", requireAuth(deleteUserHandler(userRepository)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello from net/http"))
}
