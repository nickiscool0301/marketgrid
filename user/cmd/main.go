package main

import (
	"fmt"
	"log"
	"marketgrid/user/internal/adapter/driven/persistence"
	"marketgrid/user/internal/adapter/primary/http"
	"marketgrid/user/internal/application/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := fmt.Sprintf("host=localhost user=marketgrid_user password=marketgrid_password dbname=marketgrid_db port=5432 sslmode=disable")
	userRepo, err := persistence.NewPostgresUserRepository(dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := userRepo.Init(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	userService := service.NewUserService(userRepo)

	userHandler := http.NewUserHandler(userService)

	router := gin.Default()
	userHandler.RegisterRoutes(router)

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
