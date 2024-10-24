package collector

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

// SaveCollector godoc
// @Summary      Save Collector
// @Description  Save a collector by giver form
// @Tags         collectors
// @Accept       json
// @Produce      json
// @Param		 collector	body	CreateDTO	true	"Add Collector"
// @Success      201  {object}  DTO
// @Failure      500
// @Router       /collectors [post]
func (a *API) Save(c *gin.Context) {
	createDTO := &CreateDTO{}
	err := c.ShouldBindJSON(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	newCollector, err := a.service.Save(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, newCollector)
}

// FindOneCollector godoc
// @Summary      Find one an Collector
// @Description  Find one Collector by ID
// @Tags         collectors
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Collector ID"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /collectors/{id} [get]
func (a *API) FindOne(c *gin.Context) {
	id := c.Param("id")
	collector, err := a.service.FindOne(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, collector)
}

// UpdateOneCollector godoc
// @Summary      Update one an Collector
// @Description  Update one Collector by ID
// @Tags         collectors
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Collector ID"
// @Param        collector   body	UpdateDTO  true  "Update Collector"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /collectors/{id} [post]
func (a *API) UpdateOne(c *gin.Context) {
	id := c.Param("id")
	updateDTO := &UpdateDTO{}
	err := c.ShouldBindJSON(updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	collector, err := a.service.UpdateOne(id, updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, collector)
}

// FilterCollectors godoc
// @Summary      Filter Collectors
// @Description  Filter Collectors by query paramenters
// @Tags         collectors
// @Accept       json
// @Produce      json
// @Param        name   query      string  false  "Collector name"
// @Param        start   query      string  false  "Collector createdAt start date"
// @Param        end   query      string  false  "Collector createdAt end date"
// @Success      200  {object}  []DTO
// @Failure      400
// @Failure      500
// @Router       /collectors [get]
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
