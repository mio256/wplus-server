package test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/stretchr/testify/require"
)

func CreateOffice(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.Office)) *rdb.Office {
	t.Helper()

	target := &rdb.Office{
		Name: faker.Username(),
	}

	if f != nil {
		f(target)
	}

	created, err := rdb.New(db).TestCreateOffice(ctx, target.Name)

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteOffice(ctx, created.ID))
	})

	return &created
}

func CheckDeletedOffice(t *testing.T, ctx context.Context, db rdb.DBTX, id int64) time.Time {
	t.Helper()

	deletedAt, err := rdb.New(db).TestCheckDeletedOffice(ctx, id)
	require.NoError(t, err)
	require.Equal(t, true, deletedAt.Valid)

	return deletedAt.Time
}
