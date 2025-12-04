package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tariff-api/internal/db2"
)

// SyncUchUzel copies uch_uzel records from DB2 into Postgres.
// It upserts on the (uch, stan1, stan2, stan3) key to keep the table idempotent.
func SyncUchUzel(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectUchUzel = `
		SELECT NUM, STAN1, STAN2, STAN3, COR_TIME, COR_TIP, DATE_ND
		FROM NSI_LAYER.FKI_UCH
		WHERE ADM = 27 AND DOR = 68
	`

	const upsertUchUzel = `
		INSERT INTO uch_uzel (uch, stan1, stan2, stan3, updated_at, tip_corr, date_start)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (uch, stan1, stan2, stan3) DO UPDATE
		SET updated_at = EXCLUDED.updated_at,
			tip_corr = EXCLUDED.tip_corr,
			date_start = EXCLUDED.date_start
	`

	rows, err := db2Client.QueryRows(ctx, selectUchUzel)
	if err != nil {
		return 0, fmt.Errorf("query DB2 uch_uzel rows: %w", err)
	}
	defer rows.Close()

	validStan := make(map[string]struct{})
	stanRows, err := pg.QueryContext(ctx, `SELECT stan_kod FROM stan`)
	if err != nil {
		return 0, fmt.Errorf("load stan codes: %w", err)
	}
	for stanRows.Next() {
		var code string
		if err := stanRows.Scan(&code); err != nil {
			stanRows.Close()
			return 0, fmt.Errorf("scan stan code: %w", err)
		}
		validStan[code] = struct{}{}
	}
	if err := stanRows.Close(); err != nil {
		return 0, fmt.Errorf("close stan codes: %w", err)
	}
	if err := stanRows.Err(); err != nil {
		return 0, fmt.Errorf("iterate stan codes: %w", err)
	}

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertUchUzel)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare uch_uzel upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			num     sql.NullInt64
			stan1   sql.NullString
			stan2   sql.NullString
			stan3   sql.NullString
			corTime sql.NullTime
			corTip  sql.NullString
			dateNd  sql.NullTime
		)

		if err := rows.Scan(&num, &stan1, &stan2, &stan3, &corTime, &corTip, &dateNd); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 uch_uzel row: %w", err)
		}

		if !num.Valid {
			continue
		}

		stan1Val := strings.TrimSpace(decodeMaybeCP1251(stan1.String))
		stan2Val := strings.TrimSpace(decodeMaybeCP1251(stan2.String))
		stan3Val := strings.TrimSpace(decodeMaybeCP1251(stan3.String))

		stan1NS := sql.NullString{}
		if stan1Val != "" {
			if _, ok := validStan[stan1Val]; ok {
				stan1NS = sql.NullString{String: stan1Val, Valid: true}
			}
		}
		stan2NS := sql.NullString{}
		if stan2Val != "" {
			if _, ok := validStan[stan2Val]; ok {
				stan2NS = sql.NullString{String: stan2Val, Valid: true}
			}
		}
		stan3NS := sql.NullString{}
		if stan3Val != "" {
			if _, ok := validStan[stan3Val]; ok {
				stan3NS = sql.NullString{String: stan3Val, Valid: true}
			}
		}
		tipCorr := strings.TrimSpace(corTip.String)

		if _, err := stmt.ExecContext(
			ctx,
			num.Int64, // uch
			stan1NS,   // stan1
			stan2NS,   // stan2
			stan3NS,   // stan3
			corTime,   // updated_at
			tipCorr,   // tip_corr
			dateNd,    // date_start
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert uch_uzel %d: %w", num.Int64, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 uch_uzel rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit uch_uzel sync: %w", err)
	}

	return inserted, nil
}
