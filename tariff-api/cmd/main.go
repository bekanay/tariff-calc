package main

import (
	"context"
	"database/sql"
	"tariff-api/internal/config"
	"tariff-api/internal/handler"
	"tariff-api/internal/logger"
	mw "tariff-api/internal/middleware"
	"tariff-api/internal/repository"
	"tariff-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("configuration load failed")
	}
	logger.SetLevel(log, cfg.LogLevel)

	if cfg.GinMode == "" {
		cfg.GinMode = gin.ReleaseMode
	}
	gin.SetMode(cfg.GinMode)

	db, err := openDB(cfg)
	if err != nil {
		log.WithError(err).Fatal("database connection failed")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Warn("database close error")
		}
	}()

	log.Info("database connection established")

	authRepo := repository.NewAuthRepository(cfg)
	authService := service.NewAuthService(authRepo, cfg.KeyCloak.RequiredRole)
	authHandler := handler.NewAuthHandler(authService)

	stationRepo := repository.NewStationRepository(db)
	stationService := service.NewStationService(stationRepo)
	stationHandler := handler.NewStationHandler(stationService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(mw.RequestLogger(log), mw.CORS())

	r.POST("/login", authHandler.Login)
	r.POST("/refresh-token", authHandler.RefreshToken)
	r.POST("/logout", authHandler.Logout)
	r.GET("/roles/:role", authHandler.CheckRole)

	r.GET("/stations", mw.RequireRole(authService, cfg.KeyCloak.RequiredRole), stationHandler.GetStations)

	if err := r.Run(":8081"); err != nil {
		log.WithError(err).Fatal("server shutdown due to startup failure")
	}
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
