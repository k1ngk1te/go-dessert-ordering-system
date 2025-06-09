package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func openDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn environment variable not set or loaded")
	}
	// sql.Open doesn't actually connect to the database yet; it just validates the DSN format.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
