package main

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/spf13/cobra"
	"github.com/taxio/errors"
)

func userSubCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "user",
	}
	cmd.AddCommand(
		createUserCmd(ctx),
	)
	return cmd
}

func createUserCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "load",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			dbConn := infra.ConnectDB(ctx)
			defer dbConn.Close()

			tx, err := dbConn.Begin(ctx)
			if err != nil {
				return errors.Wrap(err)
			}
			defer util.DeferRollback(ctx, tx)

			repo := rdb.New(tx)

			officeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errors.Wrap(err)
			}
			userID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return errors.Wrap(err)
			}
			name := args[2]
			password, err := util.GeneratePasswordHash(args[3])
			if err != nil {
				return errors.Wrap(err)
			}
			role := rdb.UserType(args[4])

			if _, err := repo.LoadCreateUser(ctx, rdb.LoadCreateUserParams{
				OfficeID: int64(officeID),
				ID:       int64(userID),
				Name:     name,
				Password: password,
				Role:     role,
			}); err != nil {
				if !errors.Is(err, pgx.ErrNoRows) {
					return errors.Wrap(err)
				}
			}

			if err := tx.Commit(ctx); err != nil {
				return errors.Wrap(err)
			}

			return nil
		},
	}
	return cmd
}
