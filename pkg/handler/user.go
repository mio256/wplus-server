package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/taxio/errors"
	"net/http"
)

func PostUserAndEmployee(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	user := c.MustGet("user").(*util.UserClaims)

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you are not admin",
		})
		return
	}

	var input struct {
		Name        string `json:"name"`
		WorkplaceID int64  `json:"workplace_id"`
		Role        string `json:"role"`
		Password    string `json:"password"`
	}
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

	employee, err := repo.CreateEmployee(c, rdb.CreateEmployeeParams{
		Name:        input.Name,
		WorkplaceID: input.WorkplaceID,
	})
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	created, err := repo.CreateUser(c, rdb.CreateUserParams{
		OfficeID: int64(user.OfficeID),
		Name:     input.Name,
		Password: input.Password,
		Role:     rdb.UserType(input.Role),
		EmployeeID: pgtype.Int8{
			Int64: employee.ID,
			Valid: true,
		},
	})
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	c.IndentedJSON(http.StatusCreated, created)
}
