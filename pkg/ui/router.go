package ui

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/handler"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/util"
)

const LoginPath = "/login"
const OfficePath = "/offices"
const WorkplacePath = "/workplaces"
const EmployeePath = "/employees"
const WorkEntryPath = "/work_entries"

func DBContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn := infra.ConnectDB(c)
		c.Set("db", dbConn)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(DBContext())

	// ping
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// db-ping
	r.GET("/db-ping", func(ctx *gin.Context) {
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
	// office
	p.GET(OfficePath, handler.GetOffices)
	p.POST(OfficePath, handler.PostOffice)
	p.DELETE(OfficePath+"/:id", handler.DeleteOffice)
	// workplace
	p.GET(WorkplacePath+"/:office_id", handler.GetWorkplaces)
	p.POST(WorkplacePath, handler.PostWorkplace)
	p.DELETE(WorkplacePath+"/:id", handler.DeleteWorkplace)
	// employee
	p.GET(EmployeePath+"/:workplace_id", handler.GetEmployees)
	p.GET(EmployeePath+"/:workplace_id/:id", handler.GetEmployee)
	p.POST(EmployeePath, handler.PostEmployee)
	p.PUT(EmployeePath+"/:id", handler.ChangeEmployeeWorkplace)
	p.DELETE(EmployeePath+"/:id", handler.DeleteEmployee)
	// work_entry
	p.GET(WorkEntryPath+"/employee/:employee_id", handler.GetWorkEntries)
	p.GET(WorkEntryPath+"/office/:office_id", handler.GetWorkEntriesByOffice)
	p.GET(WorkEntryPath+"/workplace/:office_id", handler.GetWorkEntriesByWorkplace)
	p.POST(WorkEntryPath, handler.PostWorkEntry)
	p.DELETE(WorkEntryPath+"/:id", handler.DeleteWorkEntry)

	return r
}
