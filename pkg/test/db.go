package test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	dbConn, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		dbConn.Close()
	})

	return dbConn
}
