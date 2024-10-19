package user_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"github.com/leonardocartaxo/expenses-analyzer/internal/utils/logger"
	testUtil "github.com/leonardocartaxo/expenses-analyzer/internal/utils/test"
	mockDB "github.com/leonardocartaxo/expenses-analyzer/mock"
	"testing"
)

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	l := logger.NewLogger(-4)
	repo := user.NewRepository(db, l)

	mockRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(uuid.New(), "Bob").
		AddRow(uuid.New(), "Alice").
		AddRow(uuid.New(), "Doug").
		AddRow(uuid.New(), "Charlie")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"").WillReturnRows(mockRows)

	users, err := repo.All()
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(users), 4)
}
