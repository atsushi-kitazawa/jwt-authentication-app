package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func createUserHandler(userRepository *userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := parseUserRequest(w, r)
		if !ok {
			return
		}

		passwordHash, err := hashPassword(req.Password)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		createdUser, err := userRepository.create(req.Name, req.Email, passwordHash)
		if err == errUserNameAlreadyExists {
			writeJSON(w, http.StatusConflict, errorResponse{Message: "user name already exists"})
			return
		}
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusCreated, createdUser)
	}
}

func updateUserHandler(userRepository *userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := parseUserRequest(w, r)
		if !ok {
			return
		}

		passwordHash, err := hashPassword(req.Password)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		id := r.PathValue("id")
		updatedUser, err := userRepository.update(id, req.Name, req.Email, passwordHash)
		if err == errUserNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Message: "user not found"})
			return
		}
		if err == errUserNameAlreadyExists {
			writeJSON(w, http.StatusConflict, errorResponse{Message: "user name already exists"})
			return
		}
		if err != nil {
			http.Error(w, "failed to update user", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, updatedUser)
	}
}

func listUsersHandler(userRepository *userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepository.list()
		if err != nil {
			http.Error(w, "failed to list users", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, users)
	}
}

func deleteUserHandler(userRepository *userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		err := userRepository.delete(id)
		if err == errUserNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Message: "user not found"})
			return
		}
		if err != nil {
			http.Error(w, "failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getUserHandler(userRepository *userRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		foundUser, err := userRepository.get(id)
		if err == errUserNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Message: "user not found"})
			return
		}
		if err != nil {
			http.Error(w, "failed to fetch user", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, foundUser)
	}
}

func parseUserRequest(w http.ResponseWriter, r *http.Request) (createUserRequest, bool) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return createUserRequest{}, false
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "name, email and password are required"})
		return createUserRequest{}, false
	}

	return req, true
}

func parseLoginRequest(w http.ResponseWriter, r *http.Request) (loginRequest, bool) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return loginRequest{}, false
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	if req.Name == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "name and password are required"})
		return loginRequest{}, false
	}

	return req, true
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
