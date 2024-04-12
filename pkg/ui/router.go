package ui

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/handler"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/util"
)

const LoginPath = "/login/"
const OfficePath = "/offices/"
const WorkplacePath = "/workplaces/"
const EmployeePath = "/employees/"
const WorkEntryPath = "/work_entries/"

func DBContext() gin.HandlerFunc {
	ctx := context.Background()
	dbConn := infra.ConnectDB(ctx)
	return func(c *gin.Context) {
		c.Set("db", dbConn)
	}
}

func UserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, err := util.GetUserClaims(c)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		c.Set("user", userClaims)
	}

}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(DBContext())

	// ping
	r.GET("/ping/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// db-ping
	r.GET("/db-ping/", func(ctx *gin.Context) {
		infra.CheckConnectDB(ctx)
		ctx.JSON(200, gin.H{
			"message": "db-pong",
		})
	})

	// login
	r.POST(LoginPath, handler.PostLogin)

	// private
	p := r.Group("")
	p.Use(util.AuthMiddleware)
	p.Use(UserContext())
	// office
	p.GET(OfficePath, handler.GetOffice)
	// workplace
	p.GET(WorkplacePath, handler.GetWorkplaces)
	p.GET(WorkplacePath+":id/", handler.GetWorkplace)
	p.POST(WorkplacePath, handler.PostWorkplace)
	p.DELETE(WorkplacePath+":id/", handler.DeleteWorkplace)
	// employee
	p.GET(EmployeePath, handler.GetEmployeesByOffice)
	p.GET(EmployeePath+"workplace/:workplace_id/", handler.GetEmployees)
	p.GET(EmployeePath+":id/", handler.GetEmployee)
	p.POST(EmployeePath, handler.PostEmployee)
	p.PUT(EmployeePath+":id/", handler.ChangeEmployeeWorkplace)
	p.DELETE(EmployeePath+":id/", handler.DeleteEmployee)
	// work_entry
	p.GET(WorkEntryPath, handler.GetWorkEntriesByOffice)
	p.GET(WorkEntryPath+"workplace/:workplace_id/", handler.GetWorkEntriesByWorkplace)
	p.GET(WorkEntryPath+"employee/:employee_id/", handler.GetWorkEntries)
	p.POST(WorkEntryPath, handler.PostWorkEntry)
	p.DELETE(WorkEntryPath+":id/", handler.DeleteWorkEntry)

	return r
}
