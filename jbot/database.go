package jbot

import (
	"database/sql"

	_ "github.com/lib/pq" // blank import to use PostgreSQL
)

// connected returns true if d is connected to a database
func connected(d *sql.DB) bool {
	return d.Ping() == nil
}
