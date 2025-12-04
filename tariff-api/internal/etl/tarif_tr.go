package etl

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"tariff-api/internal/db2"
)

// SyncTarifTr copies tariff distances between transit points from DB2 into Postgres.
// It upserts deterministically on a generated id derived from (tr_start, tr_end).
func SyncTarifTr(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectTarifTr = `
		SELECT STAN1, STAN2, KM, COR_TIME, COR_TIP, DATE_ND
		FROM NSI_LAYER.FKI_KM km
		WHERE km.STAN1 IN (SELECT KOD FROM NSI_LAYER.FKI_STAN WHERE ADM = 27 AND DOR = 68)
	`

	const upsertTarifTr = `
		INSERT INTO tarif_tr (id, tr_start, tr_end, dist_tr, updated_at, tip_corr, date_start)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET tr_start = EXCLUDED.tr_start,
			tr_end = EXCLUDED.tr_end,
			dist_tr = EXCLUDED.dist_tr,
			updated_at = EXCLUDED.updated_at,
			tip_corr = EXCLUDED.tip_corr,
			date_start = EXCLUDED.date_start
	`

	rows, err := db2Client.QueryRows(ctx, selectTarifTr)
	if err != nil {
		return 0, fmt.Errorf("query DB2 tarif_tr rows: %w", err)
	}
	defer rows.Close()

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

		id := hashPair(start, end)
		tipCorr := strings.TrimSpace(corTip.String)

		if _, err := stmt.ExecContext(
			ctx,
			id,
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

// hashPair returns a stable positive int64 hash for a pair of strings.
func hashPair(a, b string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(a))
	_, _ = h.Write([]byte{':'})
	_, _ = h.Write([]byte(b))
	return int64(h.Sum64() & math.MaxInt64)
}
