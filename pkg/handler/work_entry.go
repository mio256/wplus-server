package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/taxio/errors"
)

type PostWorkEntryParams struct {
	EmployeeID  int64  `json:"employee_id"`
	WorkplaceID int64  `json:"workplace_id"`
	Date        string `json:"date"`
	Hours       int    `json:"hours"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Attendance  bool   `json:"attendance"`
	Comment     string `json:"comment"`
}

func GetWorkEntriesByOffice(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin",
		})
		return
	}

	workEntries, err := repo.GetWorkEntriesByOffice(c, int64(user.OfficeID))
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workEntries)
}

func GetWorkEntriesByWorkplace(c *gin.Context) {
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
		c.Error(errors.Wrap(err))
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
		c.Error(errors.Wrap(err))
		return
	}
	if workplace.OfficeID != int64(user.OfficeID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "your office is different",
		})
		return
	}

	workEntries, err := repo.GetWorkEntriesByWorkplace(c, workplaceID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workEntries)
}

func GetWorkEntries(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	employeeID, err := strconv.ParseInt(c.Param("employee_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	employee, err := repo.GetEmployee(c, employeeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	if user.Role == "admin" {
		employeeOfficeID, err := repo.GetEmployeeOffice(c, employeeID)
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employeeOfficeID != int64(user.OfficeID) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your office is different",
			})
			return
		}
	} else if user.Role == "manager" {
		if user.EmployeeID == 0 {
			c.Error(errors.New("user is not manager: employee_id is not set"))
			return
		}
		me, err := repo.GetEmployee(c, int64(user.EmployeeID))
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employee.WorkplaceID != me.WorkplaceID {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your workplace is different",
			})
			return
		}
	} else if user.Role == "employee" {
		if user.EmployeeID == 0 {
			c.Error(errors.New("user is not employee: employee_id is not set"))
			return
		}
		me, err := repo.GetEmployee(c, int64(user.EmployeeID))
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employeeID != me.ID {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your employee is different",
			})
			return
		}
	}

	workEntries, err := repo.GetWorkEntriesByEmployee(c, employeeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workEntries)
}

func PostWorkEntry(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	var input PostWorkEntryParams
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	employee, err := repo.GetEmployee(c, input.EmployeeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	if employee.WorkplaceID != input.WorkplaceID {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	if user.Role == "admin" {
		employeeOfficeID, err := repo.GetEmployeeOffice(c, employee.ID)
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employeeOfficeID != int64(user.OfficeID) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your office is different",
			})
			return
		}
	} else if user.Role == "manager" {
		if user.EmployeeID == 0 {
			c.Error(errors.New("user is not manager: employee_id is not set"))
			return
		}
		me, err := repo.GetEmployee(c, int64(user.EmployeeID))
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employee.WorkplaceID != me.WorkplaceID {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your workplace is different",
			})
			return
		}
	} else if user.Role == "employee" {
		if user.EmployeeID == 0 {
			c.Error(errors.New("user is not employee: employee_id is not set"))
			return
		}
		me, err := repo.GetEmployee(c, int64(user.EmployeeID))
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employee.ID != me.ID {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your employee is different",
			})
			return
		}
	}

	var p rdb.CreateWorkEntryParams

	p.EmployeeID = input.EmployeeID
	p.WorkplaceID = input.WorkplaceID
	wp, err := repo.GetWorkplace(c, p.WorkplaceID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	date, err := time.Parse("2006-01-02T15:04:05.000+09:00", input.Date)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	p.Date = pgtype.Date{
		Time:  date,
		Valid: true,
	}

	if input.Attendance && wp.WorkType == rdb.WorkTypeAttendance {
		p.Attendance = pgtype.Bool{
			Bool:  true,
			Valid: true,
		}
	} else if input.Hours > 0 && wp.WorkType == rdb.WorkTypeHours {
		p.Hours = pgtype.Int2{
			Int16: int16(input.Hours),
			Valid: true,
		}
	} else if input.StartTime != "" && input.EndTime != "" && wp.WorkType == rdb.WorkTypeTime {
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
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

	user := c.MustGet("user").(*util.UserClaims)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workEntry, err := repo.GetWorkEntry(c, id)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	employee, err := repo.GetEmployee(c, workEntry.EmployeeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	if user.Role == "admin" {
		employeeOfficeID, err := repo.GetEmployeeOffice(c, employee.ID)
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		if employeeOfficeID != int64(user.OfficeID) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your office is different",
			})
			return
		}
	} else if user.Role == "manager" {
		if user.EmployeeID == 0 {
			c.Error(errors.New("user is not manager: employee_id is not set"))
			return
		}
		if employee.WorkplaceID != int64(user.WorkplaceID) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your workplace is different",
			})
			return
		}
	} else if user.Role == "employee" {
		if user.EmployeeID == 0 {
			c.Error(errors.New("user is not employee: employee_id is not set"))
			return
		}
		if employee.ID != int64(user.EmployeeID) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "your employee is different",
			})
			return
		}
	}

	if err := repo.SoftDeleteWorkEntry(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	c.Status(http.StatusNoContent)
}
