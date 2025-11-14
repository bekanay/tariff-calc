package main

import (
	"context"
	"database/sql"
	"log"
	"tariff-api/internal/config"
	"tariff-api/internal/handler"
	mw "tariff-api/internal/middleware"
	"tariff-api/internal/repository"
	"tariff-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.LoadConfig()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	defer db.Close()

	log.Println("DB connected ✅")

	authRepo := repository.NewAuthRepository(cfg)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	stationRepo := repository.NewStationRepository(db)
	stationService := service.NewStationService(stationRepo)
	stationHandler := handler.NewStationHandler(stationService)

	r := gin.Default()
	r.Use(mw.CORS())

	r.POST("/login", authHandler.Login)
	r.POST("/refresh-token", authHandler.RefreshToken)
	r.POST("/logout", authHandler.Logout)

	r.GET("/stations", stationHandler.GetStations)

	r.Run(":8081")
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
