package tests

import "github.com/go-chi/chi"

type TestHandlers struct {
}

func NewTestHandlers() *TestHandlers {
	return &TestHandlers{}
}

func (th *TestHandlers) RegisterRoutes(r chi.Router) {

}
