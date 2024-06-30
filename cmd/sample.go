package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/spf13/cobra"
	"github.com/taxio/errors"
	"log"
	"time"
)

func sampleCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "sample",
	}
	cmd.AddCommand(
		createCmd(ctx),
	)
	return cmd
}

func createCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",
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

			o, err := repo.CreateOffice(ctx, "sample_office")
			if err != nil {
				return errors.Wrap(err)
			}
			log.Printf("office: id = %d, name = %s", o.ID, o.Name)

			wp, err := repo.CreateWorkplace(ctx, rdb.CreateWorkplaceParams{
				Name:     "sample_workplace_time",
				OfficeID: o.ID,
				WorkType: rdb.WorkTypeTime,
			})
			if err != nil {
				return errors.Wrap(err)
			}
			log.Printf("workplace: id = %d, name = %s, office_id = %d, work_type = %s", wp.ID, wp.Name, wp.OfficeID, wp.WorkType)

			var employees []*rdb.Employee
			for i := 0; i < 3; i++ {
				e, err := repo.CreateEmployee(ctx, rdb.CreateEmployeeParams{
					Name:        fmt.Sprintf("sample_employee_%d", i),
					WorkplaceID: wp.ID,
				})
				if err != nil {
					return errors.Wrap(err)
				}
				log.Printf("employee[%d]: id = %d, name = %s", i, e.ID, e.Name)
				employees = append(employees, &e)
			}

			var users []*rdb.User
			for i := 0; i < 3; i++ {
				hash, err := util.GeneratePasswordHash("password")
				if err != nil {
					return errors.Wrap(err)
				}
				u, err := repo.CreateUser(ctx, rdb.CreateUserParams{
					OfficeID: o.ID,
					EmployeeID: pgtype.Int8{
						Int64: employees[i].ID,
						Valid: true,
					},
					Role:     rdb.UserTypeEmployee,
					Password: hash,
				})
				if err != nil {
					return errors.Wrap(err)
				}
				log.Printf("user[%d]: id = %d, office_id = %d, employee_id = %d, role = %s", i, u.ID, u.OfficeID, u.EmployeeID.Int64, u.Role)
				users = append(users, &u)
			}

			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					we, err := repo.CreateWorkEntry(ctx, rdb.CreateWorkEntryParams{
						EmployeeID:  employees[i].ID,
						WorkplaceID: wp.ID,
						Date: pgtype.Date{
							Time:  time.Now().AddDate(0, 0, -j),
							Valid: true,
						},
						StartTime: pgtype.Time{
							Microseconds: int64((12 - j) * int(time.Hour) / int(time.Microsecond)),
							Valid:        true,
						},
						EndTime: pgtype.Time{
							Microseconds: int64((12 + j) * int(time.Hour) / int(time.Microsecond)),
							Valid:        true,
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					log.Printf("work_entry[%d][%d]: id = %d, employee_id = %d, workplace_id = %d, date = %s, start_time = %s, end_time = %s",
						i, j, we.ID, we.EmployeeID, we.WorkplaceID, we.Date.Time, we.StartTime.Microseconds, we.EndTime.Microseconds)
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
