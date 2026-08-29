package main

import (
	"log"
	"os"

	dbPooler "github.com/wksn753/umosan-backend/internal/database"
	regHandler "github.com/wksn753/umosan-backend/internal/registrations/handler"
	regInfra "github.com/wksn753/umosan-backend/internal/registrations/infrastructure"
	regModel "github.com/wksn753/umosan-backend/internal/registrations/models"
	registrationRouter "github.com/wksn753/umosan-backend/internal/registrations/router"
	regService "github.com/wksn753/umosan-backend/internal/registrations/services"
	"github.com/wksn753/umosan-backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}
	// 1. Initialize Database connection first via environment DSN
	db, err := dbPooler.NewDatabase(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Run GORM Auto-Migrations for your models
	if err := db.AutoMigrate(
		&regModel.Member{},
		&regModel.RegisterRecord{},
	); err != nil {
		log.Fatalf("Failed to run auto-migrations: %v", err)
	}

	// 3. Setup Gin Engine and Core App Router
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	app := router.NewApp(engine)
	v1Group := app.SetUpV1()

	// 4. Instantiate Dependency Injection Chain (Infra -> Service -> Handler -> Router)
	registrationInfrastructure := regInfra.NewInfrastructure(db)
	registrationService := regService.NewRegistrationService(registrationInfrastructure)
	registrationHandler := regHandler.NewRegistrationHandler(registrationService)

	// Mount modular routes onto the V1 group
	registrationRouter.RegisterRouter(v1Group, registrationHandler)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000" // Fallback default port for local development
	} else if port[0] != ':' {
		port = ":" + port // Ensure leading colon for Vercel's dynamic port numbers
	}

	log.Printf("Starting server on port %s...", port)
	if err := app.Engine.Run(port); err != nil {
		log.Fatalf("Server stopped abruptly: %v", err)
	}
}
