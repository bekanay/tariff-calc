package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"tariff-api/internal/db2"
)

// SyncStanZd copies stan_zd records from DB2 into Postgres.
// It upserts on id (KOD) to keep the table idempotent.
func SyncStanZd(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectStanZd = `
        SELECT KOD, STAN1, KM1, STAN2, KM2, STAN3, KM3, STAN4, KM4, COR_TIP, DATE_ND, COR_TIME
        FROM NSI_LAYER.FKI_STAN
        WHERE ADM = 27 AND DOR = 68
    `

	const upsertStanZd = `
        INSERT INTO stan_zd (stan, tr1, dist1, tr2, dist2, tr3, dist3, tr4, dist4, tip_corr, date_start, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        ON CONFLICT (stan) DO UPDATE
        SET tr1 = EXCLUDED.tr1,
            dist1 = EXCLUDED.dist1,
            tr2 = EXCLUDED.tr2,
            dist2 = EXCLUDED.dist2,
            tr3 = EXCLUDED.tr3,
            dist3 = EXCLUDED.dist3,
            tr4 = EXCLUDED.tr4,
            dist4 = EXCLUDED.dist4,
            tip_corr = EXCLUDED.tip_corr,
            date_start = EXCLUDED.date_start,
            updated_at = EXCLUDED.updated_at
    `

	rows, err := db2Client.QueryRows(ctx, selectStanZd)
	if err != nil {
		return 0, fmt.Errorf("query DB2 stan_zd rows: %w", err)
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
	stmt, err := tx.PrepareContext(ctx, upsertStanZd)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare stan_zd upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			kod     sql.NullInt64
			stan1   sql.NullString
			km1     sql.NullFloat64
			stan2   sql.NullString
			km2     sql.NullFloat64
			stan3   sql.NullString
			km3     sql.NullFloat64
			stan4   sql.NullString
			km4     sql.NullFloat64
			corTip  sql.NullString
			dateNd  sql.NullTime
			corTime sql.NullTime
		)

		if err := rows.Scan(&kod, &stan1, &km1, &stan2, &km2, &stan3, &km3, &stan4, &km4, &corTip, &dateNd, &corTime); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 stan_zd row: %w", err)
		}

		if !kod.Valid {
			continue
		}

		stanCode := strconv.FormatInt(kod.Int64, 10)
		if stanCode == "" {
			continue
		}

		tr1 := strings.TrimSpace(decodeMaybeCP1251(stan1.String))
		tr2 := strings.TrimSpace(decodeMaybeCP1251(stan2.String))
		tr3 := strings.TrimSpace(decodeMaybeCP1251(stan3.String))
		tr4 := strings.TrimSpace(decodeMaybeCP1251(stan4.String))
		tipCorr := strings.TrimSpace(corTip.String)

		tr1NS := sql.NullString{}
		if tr1 != "" && tr1 != "0" {
			if _, ok := validTr[tr1]; ok {
				tr1NS = sql.NullString{String: tr1, Valid: true}
			}
		}
		tr2NS := sql.NullString{}
		if tr2 != "" && tr2 != "0" {
			if _, ok := validTr[tr2]; ok {
				tr2NS = sql.NullString{String: tr2, Valid: true}
			}
		}
		tr3NS := sql.NullString{}
		if tr3 != "" && tr3 != "0" {
			if _, ok := validTr[tr3]; ok {
				tr3NS = sql.NullString{String: tr3, Valid: true}
			}
		}
		tr4NS := sql.NullString{}
		if tr4 != "" && tr4 != "0" {
			if _, ok := validTr[tr4]; ok {
				tr4NS = sql.NullString{String: tr4, Valid: true}
			}
		}

		if _, err := stmt.ExecContext(
			ctx,
			stanCode, // stan
			tr1NS,    // tr1
			km1,      // dist1
			tr2NS,    // tr2
			km2,      // dist2
			tr3NS,    // tr3
			km3,      // dist3
			tr4NS,    // tr4
			km4,      // dist4
			tipCorr,  // tip_corr
			dateNd,   // date_start
			corTime,  // updated_at
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert stan_zd %d: %w", kod.Int64, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 stan_zd rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit stan_zd sync: %w", err)
	}

	return inserted, nil
}
