package expense

import (
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"time"
)

type Expense struct {
	ID        string
	CreatedAt time.Time
	UserId    int
	User      user.DTO
	Value     float64
}
