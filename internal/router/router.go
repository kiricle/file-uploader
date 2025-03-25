package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiricle/file-uploader/internal/middleware"
	"github.com/kiricle/file-uploader/internal/services"
)

type AuthHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request) // TODO
	SignIn(w http.ResponseWriter, r *http.Request) // TODO
}

func SetupRouter(ah AuthHandler, jwtService services.JWTService) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", HealthHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", ah.SignUp)

			r.Post("/sign-in", ah.SignIn)
		})

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.JwtMiddleware(jwtService))

			r.Get("/upload", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.Context().Value("user_id"))
				fmt.Println("Protected endpoint")
				w.Write([]byte("Protected endpoint"))
			})
		})
	})

	return r
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
