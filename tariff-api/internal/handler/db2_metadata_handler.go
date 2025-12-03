package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DB2QueryHandler executes ad-hoc queries against DB2.
// Intended for manual inspection; restrict or remove in production.
type DB2QueryHandler struct {
	db *sql.DB
}

func NewDB2QueryHandler(db *sql.DB) *DB2QueryHandler {
	return &DB2QueryHandler{db: db}
}

func (h *DB2QueryHandler) Query(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide query"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, req.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []map[string]any
	for rows.Next() {
		values := make([]any, len(cols))
		dests := make([]any, len(cols))
		for i := range dests {
			dests[i] = &values[i]
		}
		if err := rows.Scan(dests...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		row := make(map[string]any, len(cols))
		for i, col := range cols {
			row[col] = normalize(values[i])
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rows": result})
}
