package user

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardocartaxo/expenses-analyzer/internal/shared"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type API struct {
	service *Service
	l       *slog.Logger
}

func NewApi(service *Service, l *slog.Logger) *API {
	return &API{service: service, l: l}
}

// SaveUser godoc
// @Summary      Save User
// @Description  Save a user by giver form
// @Tags         users
// @Accept       json
// @Produce      json
// @Param		 user	body	CreateDTO	true	"Add User"
// @Success      201  {object}  DTO
// @Failure      500
// @Router       /users [post]
func (a *API) Create(c *gin.Context) {
	createDTO := &CreateDTO{}
	err := c.ShouldBindJSON(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	newUser, err := a.service.Create(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, newUser)
}

// FindOneUser godoc
// @Summary      Find one an User
// @Description  Find one User by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /users/{id} [get]
func (a *API) FindById(c *gin.Context) {
	id := c.Param("id")
	user, err := a.service.GetById(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// UpdateOneUser godoc
// @Summary      Update one an User
// @Description  Update one User by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        user   body	UpdateDTO  true  "Update User"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /users/{id} [post]
func (a *API) UpdateById(c *gin.Context) {
	id := c.Param("id")
	updateDTO := &UpdateDTO{}
	err := c.ShouldBindJSON(updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	user, err := a.service.Update(id, updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, user)
}

// FilterUsers godoc
// @Summary      Filter Users
// @Description  Filter Users by query paramenters
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name   query      string  false  "User name"
// @Param        email   query      string  false  "User email"
// @Param        start   query      string  false  "User createdAt start date"
// @Param        end   query      string  false  "User createdAt end date"
// @Param        populate   query      string  false  "User populate properties"
// @Param        limit   query      int  false  "User pagination limit"
// @Param        offset   query      int  false  "User pagination limit"
// @Success      200  {object}  []DTO
// @Failure      400
// @Failure      500
// @Router       /users [get]
func (a *API) Find(c *gin.Context) {
	name := c.Query("name")
	email := c.Query("email")
	start := c.Query("start")
	end := c.Query("end")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	populateFieldsStr := c.Query("populateFields")

	// Query options
	var conditions []shared.BaseFindCondition
	if name != "" {
		conditions = append(conditions, shared.BaseFindCondition{Field: "name", Comparator: "=", Value: name})
	}
	if email != "" {
		conditions = append(conditions, shared.BaseFindCondition{Field: "email", Comparator: "=", Value: email})
	}
	if start != "" {
		conditions = append(conditions, shared.BaseFindCondition{Field: "created_at", Comparator: ">=", Value: start})
	}
	if end != "" {
		conditions = append(conditions, shared.BaseFindCondition{Field: "created_at", Comparator: "<=", Value: end})
	}

	var populateFields []string
	if populateFieldsStr != "" {
		populateFields = strings.Split(populateFieldsStr, ",")
	}

	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	offset := 0
	if limitStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	orderBy := "created_at desc"

	dtos, err := a.service.Find(conditions, populateFields, orderBy, limit, offset)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, dtos)
}
