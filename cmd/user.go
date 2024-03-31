package main

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		createPasswordCmd(ctx),
	)
	return cmd
}

func createPasswordCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "password",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("invalid args: password")
			}
			password, err := util.GeneratePasswordHash(args[0])
			if err != nil {
				return errors.Wrap(err)
			}
			cmd.Println(password)
			return nil
		},
	}
	return cmd
}

func createUserCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 6 {
				return errors.New("invalid args: officeID, userID, name, password, role, employeeID(null)")
			}
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
			employeeID := pgtype.Int8{}
			if args[5] != "null" {
				id, err := strconv.ParseUint(args[5], 10, 64)
				if err != nil {
					return errors.Wrap(err)
				}
				employeeOfficeID, err := repo.GetEmployeeOffice(ctx, int64(id))
				if err != nil {
					return errors.Wrap(err)
				}
				if officeID != uint64(employeeOfficeID) {
					return errors.Wrap(errors.New("invalid officeID for employeeID"))
				}
				employeeID = pgtype.Int8{Int64: int64(id), Valid: true}
			}

			params := rdb.LoadCreateUserParams{
				ID:         int64(userID),
				OfficeID:   int64(officeID),
				Name:       name,
				Password:   password,
				Role:       role,
				EmployeeID: employeeID,
			}

			if _, err := repo.LoadCreateUser(ctx, params); err != nil {
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
