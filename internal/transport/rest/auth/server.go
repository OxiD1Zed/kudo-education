package auth_server

import (
	"context"
	"encoding/json"
	"errors"
	auth_service "kode-education/internal/service/auth"
	"kode-education/internal/transport/rest"
	"net/http"
	"strings"
)

type AuthServer interface {
	Authenticate(w http.ResponseWriter, req *http.Request)
	AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type serverAPI struct {
	authService auth_service.AuthService
}

func New(
	authService auth_service.AuthService,
) AuthServer {
	return &serverAPI{
		authService: authService,
	}
}

func (s *serverAPI) Authenticate(w http.ResponseWriter, req *http.Request) {
	type RequastUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type ResponseToken struct {
		Token string `json:"token"`
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var reqUser RequastUser
	if err := dec.Decode(&reqUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if reqUser.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	if reqUser.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	token, err := s.authService.Authenticate(reqUser.Username, reqUser.Password)
	if err != nil {
		http.Error(w, "the username or password is incorrect", http.StatusNotFound)
		return
	}

	rest.RenderJSON(w, ResponseToken{Token: token})
}

func (s *serverAPI) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

		idUser, err := s.authService.Authorization(token)
		if err != nil {
			if errors.Is(err, auth_service.ErrorInvalidToken) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			http.Error(w, "failed to authorize", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(req.Context(), "idUser", idUser)
		next.ServeHTTP(w, req.WithContext(ctx))
	}
}
