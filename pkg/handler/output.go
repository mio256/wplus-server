package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/taxio/errors"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"time"
)

const TEMPLATE = "resource/template.xlsx"
const SHEET = "Sheet1"

func GetOutputByWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	if user.Role != "admin" && user.Role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin or manager",
		})
		return
	}

	workplaceID, err := strconv.ParseInt(c.Param("workplace_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.Wrap(err))
		return
	}
	if user.Role == "manager" {
		if workplaceID != int64(user.WorkplaceID) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your workplace is different",
			})
			return
		}
	}

	workplace, err := repo.GetWorkplace(c, workplaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Wrap(err))
		return
	}
	if workplace.OfficeID != int64(user.OfficeID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "your office is different",
		})
		return
	}

	var input struct {
		Year  int `json:"year"`
		Month int `json:"month"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.Wrap(err))
		return
	}

	f, err := excelize.OpenFile(TEMPLATE)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	if err := f.SetCellValue(SHEET, "C4", input.Year); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Wrap(err))
		return
	}
	if err := f.SetCellValue(SHEET, "E4", input.Month); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	entries, err := repo.OutputWorkEntriesByWorkplaceAndDate(c, rdb.OutputWorkEntriesByWorkplaceAndDateParams{
		ID: workplace.ID,
		MinDate: pgtype.Date{
			Time:  time.Date(input.Year, time.Month(input.Month), 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		MaxDate: pgtype.Date{
			Time:  time.Date(input.Year, time.Month(input.Month+1), 0, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Wrap(err))
		return
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

	if err = f.SetCellValue(SHEET, "E4", input.Month); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	i := 0
	for name, entries := range employeeEntries {
		fmt.Printf("%s:\n", name)
		if nameCell, err := excelize.CoordinatesToCellName(2, 7+int(i)*2); err == nil {
			if err := f.SetCellValue(SHEET, nameCell, name); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		for _, entry := range entries {
			// HACK: format date
			day, err := strconv.Atoi(entry.Date.Format("2006-01-02")[8:])
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			fmt.Printf("day: %d, hours: %d\n", day, entry.Hours)

			// C{7+i*2}-AG{7+i*2}
			hourCell, err := excelize.CoordinatesToCellName(3+day-1, 7+int(i)*2)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			hourValue := 0

			if hourStr, err := f.GetCellValue(SHEET, hourCell); err == nil {
				if hourStr != "" {
					hourValue, err = strconv.Atoi(hourStr)
					if err != nil {
						c.JSON(http.StatusInternalServerError, err.Error())
						return
					}
				}
			}

			if err = f.SetCellValue(SHEET, hourCell, entry.Hours+uint(hourValue)); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		i++
	}

	var b bytes.Buffer
	if f.WorkBook != nil && f.WorkBook.CalcPr != nil {
		f.WorkBook.CalcPr.FullCalcOnLoad = true
	}
	if err := f.Write(&b); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := f.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	name := fmt.Sprintf("%s-%d-%d_%s.xlsx", workplace.Name, input.Year, input.Month, time.Now().Format("20060102150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Data(http.StatusOK, "application/octet-stream", b.Bytes())
}
