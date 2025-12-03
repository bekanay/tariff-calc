package etl

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"tariff-api/internal/db2"

	"golang.org/x/text/encoding/charmap"
)

// SyncStan copies stations from DB2 into the Postgres stan table.
// It streams rows to avoid loading everything into memory and performs
// an upsert on stan_kod to keep the table idempotent.
func SyncStan(ctx context.Context, db2Client *db2.Client, pg *sql.DB) (int, error) {
	const selectStan = `
		SELECT STAN_ID, KOD, NAME, PARAGRAF, PR
		FROM FKI_STAN WHERE ADM = 27 AND DOR = 68
	`
	const upsertStan = `
		INSERT INTO stan (id, stan_kod, stan_name, stan_priznak, paragraph)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (stan_kod) DO UPDATE
		SET stan_name = EXCLUDED.stan_name,
			stan_priznak = EXCLUDED.stan_priznak,
			paragraph = EXCLUDED.paragraph
	`

	rows, err := db2Client.QueryRows(ctx, selectStan)
	if err != nil {
		return 0, fmt.Errorf("query DB2 stan rows: %w", err)
	}
	defer rows.Close()

	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, upsertStan)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("prepare upsert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for rows.Next() {
		var (
			id        sql.NullInt64
			kod       sql.NullInt64
			name      sql.NullString
			paragraph sql.NullString
			priznak   sql.NullInt64
		)
		if err := rows.Scan(&id, &kod, &name, &paragraph, &priznak); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("scan DB2 stan row: %w", err)
		}

		stanKod := ""
		if kod.Valid {
			stanKod = strconv.FormatInt(kod.Int64, 10)
		}
		stanName := strings.TrimSpace(decodeMaybeCP1251(name.String))
		decodedParagraph := strings.TrimSpace(decodeMaybeCP1251(paragraph.String))
		stanParagraph := sql.NullString{
			String: decodedParagraph,
			Valid:  paragraph.Valid && decodedParagraph != "",
		}
		stanPriznak := 0
		if priznak.Valid {
			stanPriznak = int(priznak.Int64)
		}

		if stanKod == "" || stanName == "" {
			continue
		}

		if _, err := stmt.ExecContext(ctx, id.Int64, stanKod, stanName, stanPriznak, stanParagraph); err != nil {
			_ = tx.Rollback()
			return inserted, fmt.Errorf("upsert stan %s: %w", stanKod, err)
		}
		inserted++
	}
	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("iterate DB2 stan rows: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		SELECT setval(
			pg_get_serial_sequence('stan', 'id'),
			COALESCE((SELECT MAX(id) FROM stan), 0),
			true
		)`); err != nil {
		_ = tx.Rollback()
		return inserted, fmt.Errorf("refresh stan id sequence: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("commit stan sync: %w", err)
	}
	return inserted, nil
}

// decodeMaybeCP1251 returns UTF-8 text, only decoding when input is not valid UTF-8.
func decodeMaybeCP1251(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	decoded, err := charmap.Windows1251.NewDecoder().String(s)
	if err != nil {
		return s
	}
	return decoded
}
