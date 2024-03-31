package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
)

func ConnectDB(ctx context.Context) *pgxpool.Pool {
	var (
		dbUser         = os.Getenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = os.Getenv("DB_PASS")              // e.g. 'my-db-password'
		unixSocketPath = os.Getenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
		dbName         = os.Getenv("DB_NAME")              // e.g. 'my-database'
	)

	dbConn, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s database=%s host=%s", dbUser, dbPwd, dbName, unixSocketPath))
	if err != nil {
		panic(err)
	}

	return dbConn
}

func CheckConnectDB(ctx context.Context) {
	dbConn := ConnectDB(ctx)
	repo := rdb.New(dbConn)
	if err := repo.Ping(ctx); err != nil {
		panic(err)
	}
	defer dbConn.Close()
}
