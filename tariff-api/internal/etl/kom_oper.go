package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tariff-api/internal/db2"
)

// SyncKomOper copies commercial operation codes from DB2 into Postgres.
// It upserts on opr_code to keep the table idempotent.
func SyncKomOper(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectKomOper = `
		SELECT OBJECT_KOD, OBJECT_KODSTR, OBJECT_NAME, COR_TIP, DATE_ND, DATE_KD, COR_TIME
		FROM NSI_LAYER.FKI_DIC_OBJECTS
		WHERE CLASS_ID = 124
	`

	const upsertKomOper = `
		INSERT INTO kom_oper (place_code, opr_code, opr_name, tip_corr, date_start, date_end, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (opr_code) DO UPDATE
		SET place_code = EXCLUDED.place_code,
			opr_name = EXCLUDED.opr_name,
			tip_corr = EXCLUDED.tip_corr,
			date_start = EXCLUDED.date_start,
			date_end = EXCLUDED.date_end,
			updated_at = EXCLUDED.updated_at
	`

	rows, err := db2Client.QueryRows(ctx, selectKomOper)
	if err != nil {
		return 0, fmt.Errorf("query DB2 kom_oper rows: %w", err)
	}
	defer rows.Close()

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, upsertKomOper)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare kom_oper upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			objectKod    sql.NullInt64
			objectKodStr sql.NullString
			objectName   sql.NullString
			corTip       sql.NullString
			dateNd       sql.NullTime
			dateKd       sql.NullTime
			corTime      sql.NullTime
		)

		if err := rows.Scan(&objectKod, &objectKodStr, &objectName, &corTip, &dateNd, &dateKd, &corTime); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 kom_oper row: %w", err)
		}

		if !objectKod.Valid {
			continue
		}

		oprCode := strings.TrimSpace(decodeMaybeCP1251(objectKodStr.String))
		oprName := strings.TrimSpace(decodeMaybeCP1251(objectName.String))
		tipCorr := strings.TrimSpace(corTip.String)

		if oprName == "" {
			continue
		}

		if _, err := stmt.ExecContext(
			ctx,
			objectKod,    // place_code (carry original kod if present)
			oprCode,      // opr_code
			oprName,      // opr_name
			tipCorr,      // tip_corr
			dateNd.Time,  // date_start
			dateKd.Time,  // date_end
			corTime.Time, // updated_at
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert kom_oper %d: %w", objectKod.Int64, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 kom_oper rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit kom_oper sync: %w", err)
	}

	return inserted, nil
}
