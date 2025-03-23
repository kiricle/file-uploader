package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request) // TODO
	SignIn(w http.ResponseWriter, r *http.Request) // TODO
}

func SetupRouter(ah AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", HealthHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", ah.SignUp)

			r.Post("/sign-in", ah.SignIn)
		})

	})

	return r
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
