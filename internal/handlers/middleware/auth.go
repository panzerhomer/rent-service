package middleware

import (
	"avito/internal/domain"
	"avito/internal/handlers"
	manager "avito/pkg/jwt"
	"net/http"
	"regexp"
)

func isModerator(path string) bool {
	return path == "/house/create" || path == "/flat/update"
}

func isClient(path string) bool {
	matched, _ := regexp.MatchString("/house/[0-9]+/subscribe", path)
	return path == "/flat/create" || matched
}

func AuthMiddleware(handler http.HandlerFunc, manager manager.TokenManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenJwt, err := r.Cookie("token")
		if err != nil || tokenJwt.Value == "" {
			respBoby := handlers.CreateErrorResponse(r.Context(), handlers.ErrorNotAuthorized, handlers.ErrorNotAuthorizedMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBoby)
			return
		}

		_, err = manager.ValidateJWT(tokenJwt.Value)
		if err != nil {
			respBoby := handlers.CreateErrorResponse(r.Context(), handlers.ErrorNotAuthorized, handlers.ErrorNotAuthorizedMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBoby)
			return
		}

		role, err := manager.ParseJWT(tokenJwt.Value, "role")
		if err != nil {
			respBody := handlers.CreateErrorResponse(r.Context(), handlers.ErrorExtractRoleFromToken, handlers.ErrorExtractRoleFromTokenMsg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(respBody)
			return
		}

		path := r.URL.Path
		if isModerator(path) && role != domain.Moderator {
			respBody := handlers.CreateErrorResponse(r.Context(), handlers.ErrorNoAuthorized, handlers.ErrorNoAccessMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBody)
			return
		}

		if isClient(path) && role != domain.Client {
			respBody := handlers.CreateErrorResponse(r.Context(), handlers.ErrorNoAuthorized, handlers.ErrorNoAccessMsg)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(respBody)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
