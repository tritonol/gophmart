package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tritonol/gophmart.git/internal/models/user"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req credentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "can't pasrde body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	token, err := s.auth.Register(ctx, toModel(req))
	if err != nil {
		var alreadyExists user.UserAlreadyExistsError
		if errors.As(err, alreadyExists) {
			http.Error(w, "login already taken", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	saveAuthCookie(w, token)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req credentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "can't pasrde body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "wrong data", http.StatusBadRequest)
		return
	}

	token, err := s.auth.Login(ctx, toModel(req))
	if err != nil {
		var notFound user.UserNotFoundError
		if errors.As(err, notFound) {
			http.Error(w, "wrong credentials", http.StatusUnauthorized)
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	saveAuthCookie(w, token)

	w.WriteHeader(http.StatusOK)
}

func saveAuthCookie(w http.ResponseWriter, token string) {
	bearerToken := fmt.Sprintf("Bearer %s", token)

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    bearerToken,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func toModel(c credentials) user.UserCredentials {
	return user.UserCredentials{
		Login:    c.Login,
		Password: c.Password,
	}
}
