package user

import (
	"github.com/leonardocartaxo/expenses-analyzer/internal/shared"
)

type Repository struct {
	shared.BaseRepository[Model]
}
