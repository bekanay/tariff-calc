package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tariff-api/internal/db2"
)

// SyncTarifStan copies tariff distances between stations from DB2 into Postgres.
// It upserts on (uch, stan) to keep the table idempotent.
func SyncTarifStan(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectTarifStan = `
		SELECT NUM, STAN, KM1, KM2, KM3, COR_TIME, COR_TIP, DATE_ND
		FROM NSI_LAYER.FKI_STAN_UCH
		WHERE DOR = 68
	`

	const upsertTarifStan = `
		INSERT INTO tarif_stan (uch, stan, dist1, dist2, dist3, updated_at, tip_corr, date_start)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (uch, stan) DO UPDATE
		SET dist1 = EXCLUDED.dist1,
			dist2 = EXCLUDED.dist2,
			dist3 = EXCLUDED.dist3,
			updated_at = EXCLUDED.updated_at,
			tip_corr = EXCLUDED.tip_corr,
			date_start = EXCLUDED.date_start
	`

	rows, err := db2Client.QueryRows(ctx, selectTarifStan)
	if err != nil {
		return 0, fmt.Errorf("query DB2 tarif_stan rows: %w", err)
	}
	defer rows.Close()

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertTarifStan)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare tarif_stan upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			num     sql.NullInt64
			stan    sql.NullString
			km1     sql.NullFloat64
			km2     sql.NullFloat64
			km3     sql.NullFloat64
			corTime sql.NullTime
			corTip  sql.NullString
			dateNd  sql.NullTime
		)

		if err := rows.Scan(&num, &stan, &km1, &km2, &km3, &corTime, &corTip, &dateNd); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 tarif_stan row: %w", err)
		}

		if !num.Valid {
			continue
		}

		stanKod := strings.TrimSpace(decodeMaybeCP1251(stan.String))
		if stanKod == "" {
			continue
		}

		tipCorr := strings.TrimSpace(corTip.String)

		if _, err := stmt.ExecContext(
			ctx,
			num.Int64, // uch
			stanKod,   // stan
			km1,       // dist1
			km2,       // dist2
			km3,       // dist3
			corTime,   // updated_at
			tipCorr,   // tip_corr
			dateNd,    // date_start
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert tarif_stan %d: %w", num.Int64, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 tarif_stan rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit tarif_stan sync: %w", err)
	}

	return inserted, nil
}
