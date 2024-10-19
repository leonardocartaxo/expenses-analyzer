package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/leonardocartaxo/expenses-analyzer/docs"
	"github.com/leonardocartaxo/expenses-analyzer/internal"
	"github.com/leonardocartaxo/expenses-analyzer/internal/expense"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"github.com/leonardocartaxo/expenses-analyzer/internal/utils/logger"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// @title           Expenses Analyser
// @version         0.1
// @description     Pet project for learning Golang and maybe a boilerplate for future projects
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
func main() {
	c := internal.NewConfig()
	l := logger.NewLogger(c.LogLevel)
	addr := fmt.Sprintf(":%d", c.Server.Port)
	const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
	dbString := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.User, c.DB.Pass, c.DB.Name, c.DB.Port)
	var logLevel gormlogger.LogLevel
	if c.DB.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		panic("failed to connect database")
	}
	// Auto migrate the schema
	if c.DB.AutoMigrate {
		err = db.AutoMigrate(&user.Model{})
		if err != nil {
			panic(err)
		}
		err = db.AutoMigrate(&expense.Expense{})
		if err != nil {
			panic(err)
		}
	}

	var ginMode string
	if c.Server.Debug {
		ginMode = gin.DebugMode
	} else {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)
	r := gin.Default()
	// Use the custom middleware and pass the logger
	r.Use(logger.SetRequestIDMiddleware())
	r.Use(logger.LogRequestMiddleware(l))
	r.Use(logger.LogResponseMiddleware(l))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	users := r.Group("/users")
	user.NewRouter(db, users, l).Route()
	err = r.Run(addr)
	if err != nil {
		panic(err)
	}
}
