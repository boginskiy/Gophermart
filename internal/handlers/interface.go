package handlers

import "github.com/go-chi/chi"

type Hdlrser interface {
	RegisterRoutes(r chi.Router)
}
