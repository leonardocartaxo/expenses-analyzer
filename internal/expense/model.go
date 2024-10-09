package expense

import (
	"github.com/google/uuid"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"gorm.io/gorm"
	"time"
)

type Expense struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,-" gorm:"index"`
	Value     float64        `json:"value"`
	UserId    uuid.UUID      `json:"userId"`
	User      user.Model
}
