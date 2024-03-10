package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/taxio/errors"
)

func GetEmployees(c *gin.Context) {
	dbConn := infra.ConnectDB(c)
	repo := rdb.New(dbConn)

	workplaceID, err := strconv.ParseInt(c.Param("workplace_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
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
	dbConn := infra.ConnectDB(c)
	repo := rdb.New(dbConn)

	workplaceID, err := strconv.ParseInt(c.Param("workplace_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	employee, err := repo.GetEmployee(c, rdb.GetEmployeeParams{
		WorkplaceID: workplaceID,
		ID:          id,
	})
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employee)

}

func PostEmployee(c *gin.Context) {
	dbConn := infra.ConnectDB(c)
	repo := rdb.New(dbConn)

	var input rdb.CreateEmployeeParams
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	employee, err := repo.CreateEmployee(c, input)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, employee)
}

func DeleteEmployee(c *gin.Context) {
	dbConn := infra.ConnectDB(c)
	repo := rdb.New(dbConn)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	log.Print(id)

	if err := repo.SoftDeleteEmployee(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	log.Print("success")

	c.Status(http.StatusNoContent)
}
