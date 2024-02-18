package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/taxio/errors"
)

func PostWorkEntry(c *gin.Context) {
	dbConn := infra.ConnectDB(c)
	repo := rdb.New(dbConn)

	var input rdb.CreateWorkEntryParams
	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	workEntry, err := repo.CreateWorkEntry(c, input)
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}
	c.IndentedJSON(http.StatusOK, workEntry)
}

func DeleteWorkEntry(c *gin.Context) {
	dbConn := infra.ConnectDB(c)
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
