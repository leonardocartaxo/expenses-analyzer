package internal

type CRUDFilter any

type CRUD[T any] interface {
	Save() *T
	Update(id string) *T
	Delete(id string) *T
	FindOne(id string) *T
	Find(filter CRUDFilter) *[]T
}
