package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type DTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type Model struct {
	//gorm.Model
	uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Name      string
}

func (m *Model) ToDTO() *DTO {
	return &DTO{
		ID:        m.UUID.String(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt.Time,
	}
}

//func (m []*Model) ToDTO() []*DTO  {}

type CreateDTO struct {
	Name string `json:"name"`
}

type UpdateDTO struct {
	CreateDTO
}

func (u DTO) Bind(r *http.Request) error {
	//TODO implement me
	panic("implement me")
}

func New(name string) *DTO {
	return &DTO{ID: uuid.New().String(), Name: name, CreatedAt: time.Now()}
}

type myObg struct {
	f1 string
	f2 int
}

func New2(obg myObg, myVar string) int {
	return obg.f2
}
