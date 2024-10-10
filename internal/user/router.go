package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	db *gorm.DB
	rg *gin.RouterGroup
}

func NewRouter(db *gorm.DB, rg *gin.RouterGroup) *Router {
	return &Router{db: db, rg: rg}
}

func (r *Router) Route() {
	repo := NewRepository(r.db)
	service := NewService(repo)
	api := NewApi(service)

	r.rg.POST("/", api.Save)
	r.rg.GET("/:id", api.FindOne)
	r.rg.POST("/:id", api.UpdateOne)
	r.rg.GET("/", api.All)
}
