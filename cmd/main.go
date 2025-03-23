package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/kiricle/file-uploader/internal/handlers"
	"github.com/kiricle/file-uploader/internal/router"
	"github.com/kiricle/file-uploader/internal/service"
	"github.com/kiricle/file-uploader/internal/storage/postgres"
	_ "github.com/lib/pq"
)

func main() {
	DB_URL := os.Getenv("DATABASE_URL")
	JWT_SECRET := os.Getenv("JWT_SECRET")

	storage := postgres.NewStorage(DB_URL)
	jwtService := service.NewJWTService(JWT_SECRET)
	authService := service.NewAuthService(storage, jwtService)
	validate := validator.New()
	authHandler := handlers.NewAuthHandler(validate, authService)

	r := router.SetupRouter(authHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
