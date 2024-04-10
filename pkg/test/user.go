package test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.User)) (*rdb.User, string) {
	t.Helper()

	o := CreateOffice(t, ctx, db, nil)

	target := &rdb.User{
		ID:         rand.Int63(),
		OfficeID:   o.ID,
		Name:       faker.Username(),
		Password:   faker.Password(),
		Role:       rdb.UserTypeAdmin,
		EmployeeID: pgtype.Int8{},
	}

	if f != nil {
		f(target)
	}

	hash, err := util.GeneratePasswordHash(target.Password)
	require.NoError(t, err)

	created, err := rdb.New(db).TestCreateUser(ctx, rdb.TestCreateUserParams{
		ID:         target.ID,
		OfficeID:   target.OfficeID,
		Name:       target.Name,
		Password:   hash,
		Role:       target.Role,
		EmployeeID: target.EmployeeID,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteUser(ctx, created.ID))
	})

	return &created, target.Password
}

func CreateUserWithToken(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.User)) (*rdb.User, string, string) {
	t.Helper()

	user, plain := CreateUser(t, ctx, db, f)

	userClaims := util.UserClaims{
		UserID:   uint64(user.ID),
		OfficeID: uint64(user.OfficeID),
		Name:     user.Name,
		Role:     string(user.Role),
	}
	if user.EmployeeID.Valid {
		employee, err := rdb.New(db).TestGetEmployee(ctx, user.EmployeeID.Int64)
		require.NoError(t, err)
		userClaims.EmployeeID = uint64(employee.ID)
		userClaims.WorkplaceID = uint64(employee.WorkplaceID)
	}

	token, err := util.GenerateToken(userClaims)
	require.NoError(t, err)

	return user, token, plain
}
