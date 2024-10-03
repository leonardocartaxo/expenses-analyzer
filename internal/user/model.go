package user

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Form struct {
	Name string `json:"name"`
}

func (u User) Bind(r *http.Request) error {
	//TODO implement me
	panic("implement me")
}

func New(name string) *User {
	return &User{ID: uuid.New().String(), Name: name, CreatedAt: time.Now()}
}

type myObg struct {
	f1 string
	f2 int
}

func New2(obg myObg, myVar string) int {
	return obg.f2
}
