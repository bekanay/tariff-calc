package db2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/ibmdb/go_ibm_db" // DB2 driver
	"github.com/sirupsen/logrus"
)

// Config represents minimal settings required to talk to DB2.
// Values may be left zero to fall back to sensible defaults.
type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type Client struct {
	db  *sql.DB
	log *logrus.Entry
}

// DB exposes the underlying connection pool for shared use.
// Caller must not close it directly; use Client.Close instead.
func (c *Client) DB() *sql.DB {
	return c.db
}

// New creates a pooled DB2 client using the provided DSN.
// The caller owns the returned client and should Close it when done.
func New(cfg Config, log *logrus.Logger) (*Client, error) {
	if cfg.DSN == "" {
		return nil, errors.New("db2 DSN is required")
	}

	db, err := sql.Open("go_ibm_db", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open db2 connection: %w", err)
	}

	applyPoolSettings(db, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping db2: %w", err)
	}

	entry := log.WithField("component", "db2")
	entry.Info("DB2 connection pool initialized")

	return &Client{
		db:  db,
		log: entry,
	}, nil
}

func applyPoolSettings(db *sql.DB, cfg Config) {
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 20
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 5
	}
	if cfg.ConnMaxLifetime <= 0 {
		cfg.ConnMaxLifetime = 30 * time.Minute
	}
	if cfg.ConnMaxIdleTime <= 0 {
		cfg.ConnMaxIdleTime = 10 * time.Minute
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}

// Close frees all resources held by the client.
func (c *Client) Close() error {
	return c.db.Close()
}

// Query runs a SELECT-like statement and returns each row as a map of column
// name to value. Use QueryRows for streaming large results.
func (c *Client) Query(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db2 query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("db2 query columns: %w", err)
	}

	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range dest {
			dest[i] = &values[i]
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("db2 scan: %w", err)
		}

		row := make(map[string]any, len(columns))
		for i, col := range columns {
			row[col] = normalizeValue(values[i])
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db2 rows: %w", err)
	}

	c.log.WithField("rows", len(results)).Debug("query completed")
	return results, nil
}

// QueryRows returns a raw *sql.Rows for callers that want to stream or scan
// into concrete structs. The caller must close the returned rows.
func (c *Client) QueryRows(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db2 query rows: %w", err)
	}
	return rows, nil
}

// Exec runs a single statement (INSERT/UPDATE/DELETE) and returns the number of
// affected rows.
func (c *Client) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("db2 exec: %w", err)
	}
	count, _ := res.RowsAffected()
	c.log.WithFields(logrus.Fields{"rows": count}).Debug("statement executed")
	return count, nil
}

// ExecBatch executes a series of statements sequentially. It stops at the first
// error and returns the number of statements successfully executed.
func (c *Client) ExecBatch(ctx context.Context, queries []string) (int, error) {
	successful := 0
	for i, query := range queries {
		if query == "" {
			continue
		}
		if _, err := c.Exec(ctx, query); err != nil {
			return successful, fmt.Errorf("db2 exec batch at index %d: %w", i, err)
		}
		successful++
	}
	return successful, nil
}

// ExportCSV calls DB2 ADMIN_CMD with a single command string (e.g. EXPORT TO ...).
func (c *Client) ExportCSV(ctx context.Context, command string) error {
	stmt, err := c.db.PrepareContext(ctx, "CALL SYSPROC.ADMIN_CMD(?)")
	if err != nil {
		return fmt.Errorf("prepare export command: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, command); err != nil {
		return fmt.Errorf("execute export command: %w", err)
	}

	c.log.Info("DB2 export command executed")
	return nil
}

func normalizeValue(v any) any {
	switch t := v.(type) {
	case []byte:
		return string(t)
	default:
		return t
	}
}
