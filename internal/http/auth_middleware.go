package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/tritonol/gophmart.git/internal/models/user"
)

const keyUserId = "userId"

func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var token string

		splitedCookie := strings.Split(cookie.Value, " ")
		if len(splitedCookie) == 2 {
			token = splitedCookie[1]
		}

		var userId user.UserID

		userId, err = s.auth.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), keyUserId, userId)
		r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
