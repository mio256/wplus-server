package test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/stretchr/testify/require"
)

func CreateEmployee(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.Employee)) *rdb.Employee {
	t.Helper()

	wp := CreateWorkplace(t, ctx, db, nil)

	target := &rdb.Employee{
		Name:        faker.Username(),
		WorkplaceID: wp.ID,
	}

	if f != nil {
		f(target)
	}

	created, err := rdb.New(db).TestCreateEmployee(ctx, rdb.TestCreateEmployeeParams{
		Name:        target.Name,
		WorkplaceID: target.WorkplaceID,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteEmployee(ctx, created.ID))
	})

	return &created
}

func CreateEmployeeWithWorkplace(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.Employee)) (*rdb.Employee, *rdb.Workplace) {
	t.Helper()

	wp := CreateWorkplace(t, ctx, db, nil)

	target := &rdb.Employee{
		Name:        faker.Username(),
		WorkplaceID: wp.ID,
	}

	if f != nil {
		f(target)
	}

	created, err := rdb.New(db).TestCreateEmployee(ctx, rdb.TestCreateEmployeeParams{
		Name:        target.Name,
		WorkplaceID: target.WorkplaceID,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteEmployee(ctx, created.ID))
	})

	return &created, wp
}

func GetDeletedAtEmployee(t *testing.T, ctx context.Context, db rdb.DBTX, id int64) time.Time {
	t.Helper()

	deletedAt, err := rdb.New(db).TestGetDeletedAtEmployee(ctx, id)
	require.NoError(t, err)
	require.Equal(t, true, deletedAt.Valid)

	return deletedAt.Time
}
