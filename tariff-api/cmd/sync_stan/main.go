package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"tariff-api/internal/db2"
	"tariff-api/internal/etl"
	"tariff-api/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	log := logger.New()

	db2DSN := os.Getenv("DB2_DSN")
	pgDSN := os.Getenv("PG_DSN")
	if db2DSN == "" || pgDSN == "" {
		log.Fatal("DB2_DSN and PG_DSN environment variables are required")
	}

	db2Client, err := db2.New(db2.Config{DSN: db2DSN}, log)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to DB2")
	}
	defer db2Client.Close()

	pgDB, err := sql.Open("pgx", pgDSN)
	if err != nil {
		log.WithError(err).Fatal("failed to open postgres connection")
	}
	pgDB.SetMaxOpenConns(10)
	pgDB.SetMaxIdleConns(5)
	pgDB.SetConnMaxLifetime(30 * time.Minute)
	pgDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := pingWithTimeout(pgDB, 5*time.Second); err != nil {
		log.WithError(err).Fatal("postgres ping failed")
	}
	defer pgDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	rows, err := etl.SyncStan(ctx, db2Client, pgDB)
	if err != nil {
		log.WithError(err).Fatal("stan sync failed")
	}
	log.Infof("stan sync complete: upserted %d rows", rows)

	tipUchRows, err := etl.SyncTipUch(ctx, db2Client, pgDB)
	if err != nil {
		log.WithError(err).Fatal("tip_uch sync failed")
	}
	log.Infof("tip_uch sync complete: upserted %d rows", tipUchRows)

	uchRows, err := etl.SyncUch(ctx, db2Client, pgDB)
	if err != nil {
		log.WithError(err).Fatal("uch sync failed")
	}
	log.Infof("uch sync complete: upserted %d rows", uchRows)

	trPunktRows, err := etl.SyncTrPunkt(ctx, db2Client, pgDB)
	if err != nil {
		log.WithError(err).Fatal("tr_punkt sync failed")
	}
	log.Infof("tr_punkt sync complete: upserted %d rows", trPunktRows)
}

func pingWithTimeout(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.PingContext(ctx)
}
