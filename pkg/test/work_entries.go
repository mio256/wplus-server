package test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/stretchr/testify/require"
)

func CreateWorkEntries(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.WorkEntry)) *rdb.WorkEntry {
	t.Helper()

	employee := CreateEmployee(t, ctx, db, nil)
	wp := CreateWorkplace(t, ctx, db, nil)

	target := &rdb.WorkEntry{
		EmployeeID:  employee.ID,
		WorkplaceID: wp.ID,
		Date:        pgtype.Date{Time: time.Now(), Valid: true},
		Hours:       pgtype.Int2{Int16: int16(rand.Int31n(24)), Valid: true},
		StartTime:   pgtype.Time{},
		EndTime:     pgtype.Time{},
		Attendance:  pgtype.Bool{},
		Comment:     pgtype.Text{String: "test", Valid: true},
	}

	if f != nil {
		f(target)
	}

	created, err := rdb.New(db).TestCreateWorkEntry(ctx, rdb.TestCreateWorkEntryParams{
		EmployeeID:  target.EmployeeID,
		WorkplaceID: target.WorkplaceID,
		Date:        target.Date,
		Hours:       target.Hours,
		StartTime:   target.StartTime,
		EndTime:     target.EndTime,
		Attendance:  target.Attendance,
		Comment:     target.Comment,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteWorkEntry(ctx, created.ID))
	})

	return &created
}

func GetDeletedAtWorkEntry(t *testing.T, ctx context.Context, db rdb.DBTX, id int64) time.Time {
	t.Helper()

	deletedAt, err := rdb.New(db).TestGetDeletedAtWorkEntry(ctx, id)
	require.NoError(t, err)
	require.Equal(t, true, deletedAt.Valid)

	return deletedAt.Time
}
