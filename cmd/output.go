package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	"log"
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
	"github.com/xuri/excelize/v2"
)

const TEMPLATE = "resource/template.xlsx"
const SHEET = "Sheet1"

func outputCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "output",
	}
	cmd.AddCommand(
		outputCSVCmd(ctx),
		outputElsxCmd(ctx),
	)
	return cmd
}

func outputElsxCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "xlsx",
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

			if len(args) != 2 {
				return errors.New("invalid args: <workplace_id> <yyyy/mm>")
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errors.Wrap(err)
			}
			// HACK: ここで年月を分割しているが、本来は正規表現でチェックするべき
			year, err := strconv.Atoi(args[1][:4])
			if err != nil {
				return errors.Wrap(err)
			}
			month, err := strconv.Atoi(args[1][5:])
			if err != nil {
				return errors.Wrap(err)
			}

			f, err := excelize.OpenFile(TEMPLATE)
			if err != nil {
				return errors.Wrap(err)
			}
			defer func() {
				if f.WorkBook != nil && f.WorkBook.CalcPr != nil {
					f.WorkBook.CalcPr.FullCalcOnLoad = true
				}
				if err := f.SaveAs("output.xlsx"); err != nil {
					log.Fatal(err)
				}
			}()

			if err := f.SetCellValue(SHEET, "C4", year); err != nil {
				return errors.Wrap(err)
			}
			if err := f.SetCellValue(SHEET, "E4", month); err != nil {
				return errors.Wrap(err)
			}

			entries, err := repo.OutputWorkEntriesByWorkplaceAndDate(ctx, rdb.OutputWorkEntriesByWorkplaceAndDateParams{
				ID: int64(id),
				MinDate: pgtype.Date{
					Time:  time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC),
					Valid: true,
				},
				MaxDate: pgtype.Date{
					Time:  time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC),
					Valid: true,
				},
			})
			if err != nil {
				return errors.Wrap(err)
			}

			type WorkEntry struct {
				Date  time.Time
				Hours uint
			}

			employeeEntries := map[string][]WorkEntry{}
			for _, e := range entries {
				employeeEntries[e.EmployeeName] = append(employeeEntries[e.EmployeeName], WorkEntry{
					Date:  e.Date.Time,
					Hours: uint(e.EndTime.Microseconds-e.StartTime.Microseconds) / 3600 / 1000000,
				})
			}

			if err = f.SetCellValue(SHEET, "E4", month); err != nil {
				return err
			}

			i := 0
			for name, entries := range employeeEntries {
				fmt.Printf("%s:\n", name)
				if nameCell, err := excelize.CoordinatesToCellName(2, 7+int(i)*2); err == nil {
					if err := f.SetCellValue(SHEET, nameCell, name); err != nil {
						return err
					}
				}
				for _, entry := range entries {
					// HACK: format date
					day, err := strconv.Atoi(entry.Date.Format("2006-01-02")[8:])
					if err != nil {
						return err
					}

					fmt.Printf("day: %d, hours: %d\n", day, entry.Hours)

					// C{7+i*2}-AG{7+i*2}
					hourCell, err := excelize.CoordinatesToCellName(3+day-1, 7+int(i)*2)
					if err != nil {
						return err
					}

					hourValue := 0

					if hourStr, err := f.GetCellValue(SHEET, hourCell); err == nil {
						if hourStr != "" {
							hourValue, err = strconv.Atoi(hourStr)
							if err != nil {
								return err
							}
						}
					}

					if err = f.SetCellValue(SHEET, hourCell, entry.Hours+uint(hourValue)); err != nil {
						return err
					}
				}
				i++
			}

			return nil
		},
	}
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
