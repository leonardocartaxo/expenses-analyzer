package user

import (
	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router) {
	usersArr := []User{
		*New("John Doe"),
		*New("Alice"),
		*New("Bob"),
	}
	usersMap := map[string]User{}
	for _, it := range usersArr {
		usersMap[it.ID] = it
	}
	repo := NewRepository(&usersMap)
	service := NewService(repo)
	api := NewApi(service)
	r.Post("/", api.Save)
	r.Get("/{id}", api.FindOne)
	r.Get("/", api.All)
}
