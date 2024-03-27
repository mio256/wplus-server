package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/gocarina/gocsv"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/spf13/cobra"
	"github.com/taxio/errors"
)

func outputCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "output",
	}
	cmd.AddCommand(
		outputCSVCmd(ctx),
	)
	return cmd
}

// Entries TODO: switch to use entry type
type Entries struct {
	Date       string `json:"date"`
	Hours      int    `json:"hours"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Attendance bool   `json:"attendance"`
	Comment    string `json:"comment"`
}

func outputCSVCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "csv",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("invalid args: <employee_id> <output file name>")
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errors.Wrap(err)
			}

			arg := struct {
				EmployeeID uint64
				FileName   string
			}{
				EmployeeID: id,
				FileName:   args[1],
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

			file, err := os.Create(arg.FileName)
			if err != nil {
				return errors.Wrap(err)
			}
			defer func() {
				if closeErr := file.Close(); closeErr != nil {
					err = fmt.Errorf("original error: %v, defer close error: %v", err, closeErr)
				}
			}()

			f := func(in io.Writer) *gocsv.SafeCSVWriter {
				r := gocsv.NewSafeCSVWriter(csv.NewWriter(transform.NewWriter(in, japanese.ShiftJIS.NewEncoder())))
				r.UseCRLF = true
				return r
			}
			gocsv.SetCSVWriter(f)

			entries, err := repo.GetWorkEntriesByEmployee(ctx, int64(arg.EmployeeID))
			if err != nil {
				return errors.Wrap(err)
			}

			utc, err := time.LoadLocation("UTC")
			if err != nil {
				return errors.Wrap(err)
			}

			var data []Entries
			for _, e := range entries {
				rec := Entries{
					Date:       e.Date.Time.Format("2006-01-02"),
					Hours:      int(e.Hours.Int16),
					StartTime:  time.UnixMicro(e.StartTime.Microseconds).In(utc).Format("15:04:05"),
					EndTime:    time.UnixMicro(e.EndTime.Microseconds).In(utc).Format("15:04:05"),
					Attendance: e.Attendance.Bool,
					Comment:    e.Comment.String,
				}
				data = append(data, rec)
			}

			if err := gocsv.MarshalFile(&data, file); err != nil {
				return errors.Wrap(err)
			}

			if err := tx.Commit(ctx); err != nil {
				return errors.Wrap(err)
			}

			return nil
		},
	}
	return cmd
}
