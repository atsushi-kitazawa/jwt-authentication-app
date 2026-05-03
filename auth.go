package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var errInvalidCredentials = errors.New("invalid credentials")

type authenticator struct {
	secretKey []byte
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func newAuthenticator(secret string) *authenticator {
	return &authenticator{secretKey: []byte(secret)}
}

func (a *authenticator) authenticateUser(userRepository *userRepository, name, password string) (user, error) {
	foundUser, err := userRepository.findByName(name)
	if err != nil {
		return user{}, errInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password)); err != nil {
		return user{}, errInvalidCredentials
	}

	return user{
		ID:    foundUser.ID,
		Name:  foundUser.Name,
		Email: foundUser.Email,
	}, nil
}

func (a *authenticator) issueToken(appUser user) (string, error) {
	claims := jwt.MapClaims{
		"sub":  appUser.ID,
		"name": appUser.Name,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *authenticator) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Message: "missing authorization header"})
			return
		}

		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || tokenString == "" {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Message: "invalid authorization header"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}

			return a.secretKey, nil
		})
		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Message: "invalid token"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loginHandler(userRepository *userRepository, authenticator *authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := parseLoginRequest(w, r)
		if !ok {
			return
		}

		authenticatedUser, err := authenticator.authenticateUser(userRepository, req.Name, req.Password)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Message: "invalid name or password"})
			return
		}

		token, err := authenticator.issueToken(authenticatedUser)
		if err != nil {
			http.Error(w, "failed to issue token", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, tokenResponse{Token: token})
	}
}
