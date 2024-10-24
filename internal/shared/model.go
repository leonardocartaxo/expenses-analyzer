package shared

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ExpenseRefModel struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"not null"`
	Value       float64        `gorm:"not null"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CollectorID uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlaceID     uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Tabler interface {
	TableName() string
}

func (ExpenseRefModel) TableName() string {
	return "expenses"
}
