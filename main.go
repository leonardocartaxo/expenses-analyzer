package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonardocartaxo/expenses-analyzer/internal"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

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
