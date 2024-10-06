package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/leonardocartaxo/expenses-analyzer/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/leonardocartaxo/expenses-analyzer/internal"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
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

// @host      localhost:3000
// @BasePath  /

// @externalDocs.description  OpenAPI
func main() {
	c := internal.NewServerConfig()
	addr := fmt.Sprintf(":%d", c.Port)
	const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
	dbString := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.User, c.DB.Pass, c.DB.Name, c.DB.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	err = db.AutoMigrate(&user.Model{})
	if err != nil {
		panic("failed to migrate database")
	}

	userRouter := user.NewRouter(db)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/docs/doc.json", c.Port)), //The url pointing to API definition
	))
	r.Route("/users", userRouter.Route)
	fmt.Printf("Starting GO API service at %s ...\n", addr)
	fmt.Println(`
		 ______     ______        ______     ______   __    
		/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
		\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
		 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
		  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}
