package test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T, ctx context.Context, db rdb.DBTX, f func(v *rdb.User)) *rdb.User {
	t.Helper()

	o := CreateOffice(t, ctx, db, nil)

	target := &rdb.User{
		ID:       rand.Int63(),
		OfficeID: o.ID,
		Name:     faker.Username(),
		Password: faker.Password(),
		Role:     rdb.UserTypeAdmin,
	}

	if f != nil {
		f(target)
	}

	created, err := rdb.New(db).TestCreateUser(ctx, rdb.TestCreateUserParams{
		ID:       target.ID,
		OfficeID: target.OfficeID,
		Name:     target.Name,
		Password: target.Password,
		Role:     target.Role,
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(db).TestDeleteUser(ctx, created.ID))
	})

	return &created
}