package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/taxio/errors"
)

func GetWorkplaces(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	workplaces, err := repo.GetWorkplaces(c, int64(user.OfficeID))
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusOK, workplaces)
}

func GetWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

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

	if workplace.OfficeID != int64(user.OfficeID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "your office is different",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, workplace)
}

func PostWorkplace(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin",
		})
		return
	}

	var input rdb.CreateWorkplaceParams
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	if input.OfficeID != int64(user.OfficeID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "your office is different",
		})
		return
	}

	workplace, err := repo.CreateWorkplace(c, input)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusCreated, workplace)
}

// DeleteWorkplace NOTE: DeleteWorkplaceによって削除されるWorkplaceに属するEmployeeをChangeEmployeeWorkplaceを使用して移動させるようにフロントエンドで促す
func DeleteWorkplace(c *gin.Context) {
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

	workplace, err := repo.GetWorkplace(c, id)
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

	if err := repo.SoftDeleteWorkplace(c, id); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.Status(http.StatusNoContent)
}
