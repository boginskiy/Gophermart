package handlers

import "github.com/go-chi/chi"

type Handler interface {
	RegisterRoutes(r chi.Router)
}
