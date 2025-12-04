package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tariff-api/internal/db2"
)

// SyncTipUch copies tip_uch dictionary rows from DB2 into Postgres.
// It upserts on kod (OBJECT_KOD) to keep the table idempotent.
func SyncTipUch(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectTipUch = `
        SELECT OBJECT_KOD, OBJECT_NAME, COR_TIP, DATE_ND, DATE_KD, COR_TIME
        FROM NSI_LAYER.FKI_DIC_OBJECTS
        WHERE CLASS_ID = 123
    `

	const upsertTipUch = `
        INSERT INTO tip_uch (kod, name, tip_corr, date_start, date_end, upd_time)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (kod) DO UPDATE
        SET name = EXCLUDED.name,
            tip_corr = EXCLUDED.tip_corr,
            date_start = EXCLUDED.date_start,
            date_end = EXCLUDED.date_end,
            upd_time = EXCLUDED.upd_time
    `

	rows, err := db2Client.QueryRows(ctx, selectTipUch)
	if err != nil {
		return 0, fmt.Errorf("query DB2 tip_uch rows: %w", err)
	}
	defer rows.Close()

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertTipUch)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare tip_uch upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			objectKod  sql.NullInt64
			objectName sql.NullString
			corTip     sql.NullString
			dateNd     sql.NullTime
			dateKd     sql.NullTime
			corTime    sql.NullTime
		)

		if err := rows.Scan(&objectKod, &objectName, &corTip, &dateNd, &dateKd, &corTime); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 tip_uch row: %w", err)
		}

		if !objectKod.Valid || !dateNd.Valid || !dateKd.Valid || !corTime.Valid {
			// Skip rows missing required fields for the target schema.
			continue
		}

		name := strings.TrimSpace(decodeMaybeCP1251(objectName.String))
		if name == "" {
			continue
		}

		tipCorr := strings.TrimSpace(corTip.String)

		if _, err := stmt.ExecContext(
			ctx,
			objectKod.Int64, // kod
			name,
			tipCorr,
			dateNd.Time,
			dateKd.Time,
			corTime.Time,
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert tip_uch %d: %w", objectKod.Int64, err)
		}

		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 tip_uch rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit tip_uch sync: %w", err)
	}

	return inserted, nil
}
