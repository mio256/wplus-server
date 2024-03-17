package ui

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/handler"
	"github.com/mio256/wplus-server/pkg/util"
)

const LoginPath = "/login"
const OfficePath = "/offices"
const WorkplacePath = "/workplaces"
const EmployeePath = "/employees"
const WorkEntryPath = "/work_entries"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	// login
	r.POST(LoginPath, handler.PostLogin)

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
	p.DELETE(EmployeePath+"/:id", handler.DeleteEmployee)
	// work_entry
	p.GET(WorkEntryPath+"/:employee_id", handler.GetWorkEntries)
	p.POST(WorkEntryPath, handler.PostWorkEntry)
	p.DELETE(WorkEntryPath+"/:id", handler.DeleteWorkEntry)

	return r
}
