package main

import (
	"context"
	"log"
	"marketgrid/user/app/admin"
	"marketgrid/user/app/user"
	"marketgrid/user/config"
	"marketgrid/user/handler/http"
	"marketgrid/user/infrastructure/persistence/postgres"
	"marketgrid/user/infrastructure/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}

	// Initialize PostgreSQL
	userRepo, err := postgres.NewPostgresUserRepository(cfg.GetDSN())
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := userRepo.Init(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	// Initialize Redis
	redisConfig := redis.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	redisClient := redis.NewRedisClient(redisConfig)

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	emailExistenceService := redis.NewRedisEmailExistenceService(redisClient)

	// Initialize Bloom Filter with existing emails
	emailSyncService := user.NewEmailSyncService(userRepo, emailExistenceService)
	if err := emailSyncService.InitializeBloomFilter(ctx); err != nil {
		log.Fatalf("could not initialize Bloom Filter: %v", err)
	}

	userApp := user.NewUserApp(userRepo, emailExistenceService)
	adminApp := admin.NewAdminApp(userRepo)

	userHandler := http.NewUserHandler(userApp)
	adminHandler := http.NewAdminHandler(adminApp)

	router := gin.Default()
	userHandler.RegisterRoutes(router)
	adminHandler.RegisterRoutes(router)

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
