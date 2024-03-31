package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/util"
)

func PostLogin(c *gin.Context) {
	dbConn := infra.ConnectDB(c)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	if err = util.CompareHashAndPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid password",
		})
		return
	}

	token, err := util.GenerateToken(uint64(user.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	// TODO: return NULL if !user.EmployeeID.Valid
	workplaceID := 0
	if user.EmployeeID.Valid {
		employee, err := repo.GetEmployeeById(c, user.EmployeeID.Int64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Employee not found",
			})
			return
		}
		workplaceID = int(employee.WorkplaceID)
	}

	domain := os.Getenv("DOMAIN")
	c.SetCookie("token", token, 3600, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{
		"office_id":    user.OfficeID,
		"user_id":      user.ID,
		"employee_id":  user.EmployeeID,
		"workplace_id": workplaceID,
		"name":         user.Name,
		"role":         user.Role,
	})
}
