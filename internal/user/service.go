package user

import (
	"maps"
	"slices"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (a *Service) Save(name string) User {
	newUser := New(name)
	user := a.repository.Save(*newUser)

	return user
}

func (a *Service) FindOne(id string) User {
	user := a.repository.FindOne(id)

	return user
}

func (a *Service) All() []User {
	allMap := a.repository.All()
	allIter := maps.Values(allMap)
	allSlices := slices.Collect(allIter)

	return allSlices
}
