package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/taxio/errors"
)

func GetEmployeesByOffice(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	employees, err := repo.GetEmployeesByOffice(c, int64(user.OfficeID))
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employees)
}

func GetEmployees(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	workplaceID, err := strconv.ParseInt(c.Param("workplace_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
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

	employees, err := repo.GetEmployees(c, workplaceID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employees)
}

func GetEmployee(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	employeeOfficeID, err := repo.GetEmployeeOffice(c, id)
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

	employee, err := repo.GetEmployee(c, id)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employee)

}

func PostEmployee(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin",
		})
		return
	}

	var input rdb.CreateEmployeeParams
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workplace, err := repo.GetWorkplace(c, input.WorkplaceID)
	if workplace.OfficeID != int64(user.OfficeID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "your office is different",
		})
		return
	}

	employee, err := repo.CreateEmployee(c, input)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employee)
}

func ChangeEmployeeWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin",
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	officeID, err := repo.GetEmployeeOffice(c, id)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	if officeID != int64(user.OfficeID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "this employee is not in your office",
		})
		return
	}

	var input struct {
		WorkplaceID int64 `json:"workplace_id"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workplace, err := repo.GetWorkplace(c, input.WorkplaceID)
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

	if err := repo.UpdateEmployeeWorkplace(c, rdb.UpdateEmployeeWorkplaceParams{
		ID:          id,
		WorkplaceID: input.WorkplaceID,
	}); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	employee, err := repo.GetEmployee(c, id)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employee)
}

func DeleteEmployee(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin",
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	employeeOfficeID, err := repo.GetEmployeeOffice(c, id)
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

	if err := repo.SoftDeleteEmployee(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	if err := repo.SoftDeleteWorkEntriesByEmployee(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.Status(http.StatusNoContent)
}
