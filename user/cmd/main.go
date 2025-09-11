package main

import (
	"fmt"
	"log"
	"marketgrid/user/app/admin"
	"marketgrid/user/app/user"
	"marketgrid/user/handler/http"
	"marketgrid/user/infrastructure/persistence/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := fmt.Sprintf("host=localhost user=marketgrid_user password=marketgrid_password dbname=marketgrid_db port=5432 sslmode=disable")
	userRepo, err := postgres.NewPostgresUserRepository(dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := userRepo.Init(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	// Create app layer (use cases)
	userApp := user.NewUserApp(userRepo)
	adminApp := admin.NewAdminApp(userRepo)

	// Create handlers (primary adapters)
	userHandler := http.NewUserHandler(userApp)
	adminHandler := http.NewAdminHandler(adminApp)

	router := gin.Default()
	userHandler.RegisterRoutes(router)
	adminHandler.RegisterRoutes(router)

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
