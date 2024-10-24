package tag

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
)

type Router struct {
	db *gorm.DB
	rg *gin.RouterGroup
	l  *slog.Logger
}

func NewRouter(db *gorm.DB, rg *gin.RouterGroup, l *slog.Logger) *Router {
	return &Router{db: db, rg: rg, l: l}
}

func (r *Router) Route() {
	repo := NewRepository(r.db, r.l)
	service := NewService(repo, r.l)
	api := NewApi(service, r.l)

	r.rg.POST("/", api.Save)
	r.rg.GET("/:id", api.FindOne)
	r.rg.POST("/:id", api.UpdateOne)
	r.rg.GET("/", api.Find)
}
