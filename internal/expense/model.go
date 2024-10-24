package expense

import (
	"github.com/google/uuid"
	"github.com/leonardocartaxo/expenses-analyzer/internal/collector"
	"github.com/leonardocartaxo/expenses-analyzer/internal/place"
	"github.com/leonardocartaxo/expenses-analyzer/internal/tag"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
	Name        string          `gorm:"not null"`
	Value       float64         `gorm:"not null"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User        user.Model      `gorm:"foreignKey:UserID;references:ID"`
	CollectorID uuid.UUID       `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Collector   collector.Model `gorm:"foreignKey:CollectorID;references:ID"`
	PlaceID     uuid.UUID       `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Place       place.Model     `gorm:"foreignKey:PlaceID;references:ID"`
	Tags        []tag.Model     `gorm:"many2many:expense_tags;"`
}

type DTO struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Value       float64   `json:"value"`
	UserId      uuid.UUID `json:"userId"`
	CollectorId uuid.UUID `json:"collectorId"`
	PlaceId     uuid.UUID `json:"placeId"`
}

type CreateDTO struct {
	Value       float64   `json:"value"`
	Name        string    `json:"name"`
	UserId      uuid.UUID `json:"userId"`
	CollectorId uuid.UUID `json:"collectorId"`
	PlaceId     uuid.UUID `json:"placeId"`
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
	return "expenses"
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
