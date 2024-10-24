package expense

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type API struct {
	service *Service
	l       *slog.Logger
}

func NewApi(service *Service, l *slog.Logger) *API {
	return &API{service: service, l: l}
}

// SaveExpense godoc
// @Summary      Save Expense
// @Description  Save a expense by giver form
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param		 expense	body	CreateDTO	true	"Add Expense"
// @Success      201  {object}  DTO
// @Failure      500
// @Router       /expenses [post]
func (a *API) Save(c *gin.Context) {
	createDTO := &CreateDTO{}
	err := c.ShouldBindJSON(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	newExpense, err := a.service.Save(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, newExpense)
}

// FindOneExpense godoc
// @Summary      Find one an Expense
// @Description  Find one Expense by ID
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Expense ID"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /expenses/{id} [get]
func (a *API) FindOne(c *gin.Context) {
	id := c.Param("id")
	expense, err := a.service.FindOne(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateOneExpense godoc
// @Summary      Update one an Expense
// @Description  Update one Expense by ID
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Expense ID"
// @Param        expense   body	UpdateDTO  true  "Update Expense"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /expenses/{id} [post]
func (a *API) UpdateOne(c *gin.Context) {
	id := c.Param("id")
	updateDTO := &UpdateDTO{}
	err := c.ShouldBindJSON(updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	expense, err := a.service.UpdateOne(id, updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, expense)
}

// FilterExpenses godoc
// @Summary      Filter Expenses
// @Description  Filter Expenses by query paramenters
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        name   query      string  false  "Expense name"
// @Param        start   query      string  false  "Expense createdAt start date"
// @Param        end   query      string  false  "Expense createdAt end date"
// @Success      200  {object}  []DTO
// @Failure      400
// @Failure      500
// @Router       /expenses [get]
func (a *API) Find(c *gin.Context) {
	name := c.Query("name")
	start := c.Param("start")
	end := c.Param("end")
	dtos, err := a.service.Find(FindOptions{Name: name, Start: start, End: end})
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, dtos)
}
