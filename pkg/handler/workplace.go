package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/taxio/errors"
)

func GetWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workplace, err := repo.GetWorkplace(c, id)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workplace)

}

func GetWorkplaces(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	officeID, err := strconv.ParseInt(c.Param("office_id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workplaces, err := repo.GetWorkplaces(c, officeID)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workplaces)
}

func PostWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	var input rdb.CreateWorkplaceParams
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workplace, err := repo.CreateWorkplace(c, input)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workplace)
}

func DeleteWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	if err := repo.SoftDeleteWorkplace(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.Status(http.StatusNoContent)
}
