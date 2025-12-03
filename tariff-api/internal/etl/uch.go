package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tariff-api/internal/db2"
)

// SyncUch copies uch records from DB2 into Postgres.
// It upserts on uch_num to keep the table idempotent.
func SyncUch(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectUch = `
		SELECT NUM, NAME, TIP
		FROM NSI_LAYER.FKI_UCH
		WHERE ADM = 27 AND DOR = 68
	`

	const upsertUch = `
		INSERT INTO uch (id, uch_num, uch_name, uch_tip)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uch_num) DO UPDATE
		SET uch_name = EXCLUDED.uch_name,
			uch_tip = EXCLUDED.uch_tip
	`

	rows, err := db2Client.QueryRows(ctx, selectUch)
	if err != nil {
		return 0, fmt.Errorf("query DB2 uch rows: %w", err)
	}
	defer rows.Close()

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertUch)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare uch upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			num  sql.NullInt64
			name sql.NullString
			tip  sql.NullInt64
		)

		if err := rows.Scan(&num, &name, &tip); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 uch row: %w", err)
		}

		if !num.Valid || !tip.Valid {
			continue
		}

		uchName := strings.TrimSpace(decodeMaybeCP1251(name.String))
		if uchName == "" {
			continue
		}

		if _, err := stmt.ExecContext(
			ctx,
			num.Int64, // id
			num.Int64, // uch_num
			uchName,
			tip.Int64,
		); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert uch %d: %w", num.Int64, err)
		}
		inserted++
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 uch rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit uch sync: %w", err)
	}

	return inserted, nil
}
