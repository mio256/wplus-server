package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/taxio/errors"
)

func PostLogin(c *gin.Context) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	var input struct {
		OfficeID uint64 `json:"office_id"`
		UserID   uint64 `json:"user_id"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	user, err := repo.GetUser(c, rdb.GetUserParams{
		OfficeID: int64(input.OfficeID),
		ID:       int64(input.UserID),
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	if err = util.CompareHashAndPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var workplaceID int64
	if user.EmployeeID.Valid {
		employee, err := repo.GetEmployee(c, user.EmployeeID.Int64)
		if err != nil {
			c.Error(errors.Wrap(err))
			return
		}
		workplaceID = employee.WorkplaceID
	}

	token, err := util.GenerateToken(util.UserClaims{
		UserID:      uint64(user.ID),
		OfficeID:    uint64(user.OfficeID),
		WorkplaceID: uint64(workplaceID),
		EmployeeID:  uint64(user.EmployeeID.Int64),
		Name:        user.Name,
		Role:        string(user.Role),
	})
	if err != nil {
		c.Error(errors.Wrap(err))
		return
	}

	domain := os.Getenv("DOMAIN")
	c.SetCookie("token", token, 3600, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
