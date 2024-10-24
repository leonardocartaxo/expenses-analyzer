package tag

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

// SaveTag godoc
// @Summary      Save Tag
// @Description  Save a tag by giver form
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param		 tag	body	CreateDTO	true	"Add Tag"
// @Success      201  {object}  DTO
// @Failure      500
// @Router       /tags [post]
func (a *API) Save(c *gin.Context) {
	createDTO := &CreateDTO{}
	err := c.ShouldBindJSON(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	newTag, err := a.service.Save(createDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, newTag)
}

// FindOneTag godoc
// @Summary      Find one an Tag
// @Description  Find one Tag by ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /tags/{id} [get]
func (a *API) FindOne(c *gin.Context) {
	id := c.Param("id")
	tag, err := a.service.FindOne(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, tag)
}

// UpdateOneTag godoc
// @Summary      Update one an Tag
// @Description  Update one Tag by ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"
// @Param        tag   body	UpdateDTO  true  "Update Tag"
// @Success      200  {object}  DTO
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /tags/{id} [post]
func (a *API) UpdateOne(c *gin.Context) {
	id := c.Param("id")
	updateDTO := &UpdateDTO{}
	err := c.ShouldBindJSON(updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	tag, err := a.service.UpdateOne(id, updateDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, tag)
}

// FilterTags godoc
// @Summary      Filter Tags
// @Description  Filter Tags by query paramenters
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        name   query      string  false  "Tag name"
// @Param        start   query      string  false  "Tag createdAt start date"
// @Param        end   query      string  false  "Tag createdAt end date"
// @Success      200  {object}  []DTO
// @Failure      400
// @Failure      500
// @Router       /tags [get]
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
