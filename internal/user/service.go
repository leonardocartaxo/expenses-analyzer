package user

import (
	"github.com/leonardocartaxo/expenses-analyzer/internal/shared"
)

type Service struct {
	shared.BaseService[Model, DTO, CreateDTO, UpdateDTO]
}
