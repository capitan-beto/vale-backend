package handlers

import (
	"github.com/capitan-beto/vale-backend/pkg/middleware"
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.RateLimiter)

	r.Route("/", func(router chi.Router) {
		router.Get("/", WelcomeMessage)

		router.Post("/contestant", AddContestant)

		router.Post("/payment", Payment)
	})
}
