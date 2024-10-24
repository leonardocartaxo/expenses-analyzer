package place

import (
	"github.com/google/uuid"
	"github.com/leonardocartaxo/expenses-analyzer/internal/shared"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt           `gorm:"index"`
	Name      string                   `gorm:"uniqueIndex"`
	Expenses  []shared.ExpenseRefModel `gorm:"foreignKey:PlaceID;references:ID"`
}

type DTO struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
}

type CreateDTO struct {
	Name string `json:"name"`
}

type UpdateDTO struct {
	CreateDTO
}

type FindOptions struct {
	Name  string
	Start string
	End   string
}

type Models []*Model

type Tabler interface {
	TableName() string
}

func (Model) TableName() string {
	return "places"
}

func (m *Model) ToDTO() *DTO {
	return &DTO{
		ID:        m.ID.String(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *Models) ToDTO() []*DTO {
	dtos := []*DTO{}
	for _, model := range *m {
		dtos = append(dtos, model.ToDTO())
	}

	return dtos
}
