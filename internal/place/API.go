package place

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

// SavePlace godoc
// @Summary      Save Place
// @Description  Save a place by giver form
// @Tags         places
// @Accept       json
// @Produce      json
// @Param		 place	body	CreateDTO	true	"Add Place"
// @Success      201  {object}  DTO
// @Failure      500
// @Router       /places [post]
func (a *API) Save(c *gin.Context) {
	createDTO := &CreateDTO{}
	err := c.ShouldBindJSON(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	newPlace, err := a.service.Save(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, newPlace)
}

// FindOnePlace godoc
// @Summary      Find one an Place
// @Description  Find one Place by ID
// @Tags         places
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Place ID"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /places/{id} [get]
func (a *API) FindOne(c *gin.Context) {
	id := c.Param("id")
	place, err := a.service.FindOne(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, place)
}

// UpdateOnePlace godoc
// @Summary      Update one an Place
// @Description  Update one Place by ID
// @Tags         places
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Place ID"
// @Param        place   body	UpdateDTO  true  "Update Place"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /places/{id} [post]
func (a *API) UpdateOne(c *gin.Context) {
	id := c.Param("id")
	updateDTO := &UpdateDTO{}
	err := c.ShouldBindJSON(updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	place, err := a.service.UpdateOne(id, updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, place)
}

// FilterPlaces godoc
// @Summary      Filter Places
// @Description  Filter Places by query paramenters
// @Tags         places
// @Accept       json
// @Produce      json
// @Param        name   query      string  false  "Place name"
// @Param        start   query      string  false  "Place createdAt start date"
// @Param        end   query      string  false  "Place createdAt end date"
// @Success      200  {object}  []DTO
// @Failure      400
// @Failure      500
// @Router       /places [get]
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
