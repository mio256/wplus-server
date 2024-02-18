package test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/stretchr/testify/require"
)

func CreateWorkplace(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.Workplace)) *rdb.Workplace {
	t.Helper()

	office := CreateOffice(t, ctx, db, nil)

	target := &rdb.Workplace{
		Name:     faker.Username(),
		OfficeID: office.ID,
		WorkType: rdb.WorkTypeHours,
	}

	if f != nil {
		f(target)
	}

	created, err := rdb.New(db).TestCreateWorkplace(ctx, rdb.TestCreateWorkplaceParams{
		Name:     target.Name,
		OfficeID: target.OfficeID,
		WorkType: target.WorkType,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteWorkplace(ctx, created.ID))
	})

	return &created
}

func CheckDeletedWorkplace(t *testing.T, ctx context.Context, db rdb.DBTX, id int64) time.Time {
	t.Helper()

	deletedAt, err := rdb.New(db).TestCheckDeletedWorkplace(ctx, id)
	require.NoError(t, err)
	require.Equal(t, true, deletedAt.Valid)

	return deletedAt.Time
}
