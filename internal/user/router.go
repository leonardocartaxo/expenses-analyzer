package user

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{db: db}
}

func (l *Router) Route(r chi.Router) {
	repo := NewRepository(l.db)
	service := NewService(repo)
	api := NewApi(service)
	r.Post("/", api.Save)
	r.Get("/{id}", api.FindOne)
	r.Post("/{id}", api.UpdateOne)
	r.Get("/", api.All)
}
