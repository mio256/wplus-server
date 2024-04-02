package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/taxio/errors"
)

func GetWorkEntries(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	employeeID, err := strconv.ParseInt(c.Param("employee_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workEntries, err := repo.GetWorkEntriesByEmployee(c, employeeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workEntries)
}

func GetWorkEntriesByOffice(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	officeID, err := strconv.ParseInt(c.Param("office_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workEntries, err := repo.GetWorkEntriesByOffice(c, officeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workEntries)
}

func GetWorkEntriesByWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	workplaceID, err := strconv.ParseInt(c.Param("workplace_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workEntries, err := repo.GetWorkEntriesByWorkplace(c, workplaceID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workEntries)

}

func PostWorkEntry(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	var input struct {
		EmployeeID  int64  `json:"employee_id"`
		WorkplaceID int64  `json:"workplace_id"`
		Date        string `json:"date"`
		Hours       int    `json:"hours"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		Attendance  bool   `json:"attendance"`
		Comment     string `json:"comment"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	var p rdb.CreateWorkEntryParams

	p.EmployeeID = input.EmployeeID
	p.WorkplaceID = input.WorkplaceID
	date, err := time.Parse("2006-01-02T15:04:05.000Z", input.Date)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	p.Date = pgtype.Date{
		Time:  date,
		Valid: true,
	}

	if input.Attendance {
		p.Attendance = pgtype.Bool{
			Bool:  true,
			Valid: true,
		}
	} else if input.Hours > 0 {
		p.Hours = pgtype.Int2{
			Int16: int16(input.Hours),
			Valid: true,
		}
	} else {
		startTime, err := time.Parse("2006-01-02T15:04:05.000Z", input.StartTime)
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		p.StartTime = pgtype.Time{
			Microseconds: startTime.UnixMicro(),
			Valid:        true,
		}
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		endTime, err := time.Parse("2006-01-02T15:04:05.000Z", input.EndTime)
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		p.EndTime = pgtype.Time{
			Microseconds: endTime.UnixMicro(),
			Valid:        true,
		}
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
	}

	if p.Comment.Valid {
		p.Comment = pgtype.Text{String: input.Comment, Valid: true}
	}

	workEntry, err := repo.CreateWorkEntry(c, p)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	c.IndentedJSON(http.StatusOK, workEntry)
}

func DeleteWorkEntry(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	if err := repo.SoftDeleteWorkEntry(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	c.Status(http.StatusNoContent)
}
