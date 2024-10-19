package user

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
func (a *API) Save(c *gin.Context) {
	createDTO := &CreateDTO{}
	err := c.ShouldBindJSON(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	newUser, err := a.service.Save(createDTO)
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
func (a *API) FindOne(c *gin.Context) {
	id := c.Param("id")
	user, err := a.service.FindOne(id)
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
func (a *API) UpdateOne(c *gin.Context) {
	id := c.Param("id")
	updateDTO := &UpdateDTO{}
	err := c.ShouldBindJSON(updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	user, err := a.service.UpdateOne(id, updateDTO)
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
// @Param        start   query      string  false  "User createdAt start date"
// @Param        end   query      string  false  "User createdAt end date"
// @Success      200  {object}  []DTO
// @Failure      400
// @Failure      500
// @Router       /users [get]
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
