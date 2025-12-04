package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tariff-api/internal/db2"
)

// SyncTarifTr copies tariff distances between transit points from DB2 into Postgres.
// It upserts on the (tr_start, tr_end) pair to keep the table idempotent.
func SyncTarifTr(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectTarifTr = `
		SELECT STAN1, STAN2, KM, COR_TIME, COR_TIP, DATE_ND
		FROM NSI_LAYER.FKI_KM km
		WHERE km.STAN1 IN (SELECT KOD FROM NSI_LAYER.FKI_STAN WHERE ADM = 27 AND DOR = 68)
	`

	const upsertTarifTr = `
		INSERT INTO tarif_tr (tr_start, tr_end, dist_tr, updated_at, tip_corr, date_start)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tr_start, tr_end) DO UPDATE
		SET dist_tr = EXCLUDED.dist_tr,
			updated_at = EXCLUDED.updated_at,
			tip_corr = EXCLUDED.tip_corr,
			date_start = EXCLUDED.date_start
	`

	rows, err := db2Client.QueryRows(ctx, selectTarifTr)
	if err != nil {
		return 0, fmt.Errorf("query DB2 tarif_tr rows: %w", err)
	}
	defer rows.Close()

	validTr := make(map[string]struct{})
	trRows, err := pg.QueryContext(ctx, `SELECT tr FROM tr_punkt`)
	if err != nil {
		return 0, fmt.Errorf("load tr_punkt codes: %w", err)
	}
	for trRows.Next() {
		var code string
		if err := trRows.Scan(&code); err != nil {
			trRows.Close()
			return 0, fmt.Errorf("scan tr_punkt code: %w", err)
		}
		code = strings.TrimSpace(code)
		if code != "" {
			validTr[code] = struct{}{}
		}
	}
	if err := trRows.Close(); err != nil {
		return 0, fmt.Errorf("close tr_punkt codes: %w", err)
	}
	if err := trRows.Err(); err != nil {
		return 0, fmt.Errorf("iterate tr_punkt codes: %w", err)
	}

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertTarifTr)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare tarif_tr upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			stan1   sql.NullString
			stan2   sql.NullString
			km      sql.NullFloat64
			corTime sql.NullTime
			corTip  sql.NullString
			dateNd  sql.NullTime
		)

		if err := rows.Scan(&stan1, &stan2, &km, &corTime, &corTip, &dateNd); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 tarif_tr row: %w", err)
		}

		start := strings.TrimSpace(decodeMaybeCP1251(stan1.String))
		end := strings.TrimSpace(decodeMaybeCP1251(stan2.String))
		if start == "" || end == "" {
			continue
		}
		if _, ok := validTr[start]; !ok {
			continue
		}
		if _, ok := validTr[end]; !ok {
			continue
		}

		tipCorr := strings.TrimSpace(corTip.String)

		if _, err := stmt.ExecContext(
			ctx,
			start,
			end,
			km,
			corTime,
			tipCorr,
			dateNd,
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert tarif_tr %s-%s: %w", start, end, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 tarif_tr rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit tarif_tr sync: %w", err)
	}

	return inserted, nil
}
