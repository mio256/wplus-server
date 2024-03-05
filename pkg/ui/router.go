package ui

import (
	"github.com/gin-gonic/gin"
	"github.com/mio256/wplus-server/pkg/handler"
)

const OfficePath = "/offices"
const WorkplacePath = "/workplaces"
const EmployeePath = "/employees"
const WorkEntryPath = "/work_entries"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// office
	r.GET(OfficePath, handler.GetOffices)
	r.POST(OfficePath, handler.PostOffice)
	r.DELETE(OfficePath+"/:id", handler.DeleteOffice)
	// workplace
	r.GET(WorkplacePath+"/:office_id", handler.GetWorkplaces)
	r.POST(WorkplacePath, handler.PostWorkplace)
	r.DELETE(WorkplacePath+"/:id", handler.DeleteWorkplace)
	// employee
	// r.GET(EmployeePath+"/:workplace_id", handler.GetEmployees)
	r.POST(EmployeePath, handler.PostEmployee)
	r.DELETE(EmployeePath+"/:id", handler.DeleteEmployee)
	// work_entry
	// r.GET(WorkEntryPath+"/:employee_id", handler.GetWorkEntries)
	r.POST(WorkEntryPath, handler.PostWorkEntry)
	r.DELETE(WorkEntryPath+"/:id", handler.DeleteWorkEntry)

	return r
}
