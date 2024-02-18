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
	r.POST(OfficePath, handler.PostOffice)
	r.DELETE(OfficePath+"/:id", handler.DeleteOffice)
	// workplace
	r.POST(WorkplacePath, handler.PostWorkplace)
	r.DELETE(WorkplacePath+"/:id", handler.DeleteWorkplace)
	// employee
	r.POST(EmployeePath, handler.PostEmployee)
	r.DELETE(EmployeePath+"/:id", handler.DeleteEmployee)
	// work_entry
	r.POST(WorkEntryPath, handler.PostWorkEntry)
	r.DELETE(WorkEntryPath+"/:id", handler.DeleteWorkEntry)

	return r
}
