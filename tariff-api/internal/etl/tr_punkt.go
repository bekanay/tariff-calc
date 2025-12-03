package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"tariff-api/internal/db2"
)

// SyncTrPunkt copies transit points (tr_punkt) from DB2 into Postgres.
// It upserts on tr code to keep the table idempotent.
func SyncTrPunkt(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectTrPunkt = `
		SELECT KOD, COR_TIME, COR_TIP, DATE_ND
		FROM NSI_LAYER.FKI_STAN
		WHERE PR = 11 AND ADM = 27 AND DOR = 68
	`

	const upsertTrPunkt = `
		INSERT INTO tr_punkt (id, tr, updated_at, tip_corr, date_start)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tr) DO UPDATE
		SET updated_at = EXCLUDED.updated_at,
			tip_corr = EXCLUDED.tip_corr,
			date_start = EXCLUDED.date_start
	`

	rows, err := db2Client.QueryRows(ctx, selectTrPunkt)
	if err != nil {
		return 0, fmt.Errorf("query DB2 tr_punkt rows: %w", err)
	}
	defer rows.Close()

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertTrPunkt)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare tr_punkt upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			kod     sql.NullInt64
			corTime sql.NullTime
			corTip  sql.NullString
			dateNd  sql.NullTime
		)

		if err := rows.Scan(&kod, &corTime, &corTip, &dateNd); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 tr_punkt row: %w", err)
		}

		if !kod.Valid {
			continue
		}

		trCode := strconv.FormatInt(kod.Int64, 10)
		if trCode == "" {
			continue
		}

		tipCorr := strings.TrimSpace(corTip.String)

		if _, err := stmt.ExecContext(
			ctx,
			kod.Int64, // id
			trCode,    // tr
			corTime,   // updated_at
			tipCorr,   // tip_corr
			dateNd,    // date_start
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert tr_punkt %s: %w", trCode, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 tr_punkt rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit tr_punkt sync: %w", err)
	}

	return inserted, nil
}
