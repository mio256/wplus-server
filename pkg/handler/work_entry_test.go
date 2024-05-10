package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/handler"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/test"
	"github.com/mio256/wplus-server/pkg/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWorkEntriesByOffice(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role    rdb.UserType
		WantErr bool
	}{
		"admin": {
			Role:    rdb.UserTypeAdmin,
			WantErr: false,
		},
		"manager": {
			Role:    rdb.UserTypeManager,
			WantErr: true,
		},
		"employee": {
			Role:    rdb.UserTypeEmployee,
			WantErr: true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			office := test.CreateOffice(t, c, dbConn, nil)
			workplace := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
				v.OfficeID = office.ID
			})
			employee := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
				v.WorkplaceID = workplace.ID
			})
			created := test.CreateWorkEntries(t, c, dbConn, func(v *rdb.WorkEntry) {
				v.EmployeeID = employee.ID
				v.WorkplaceID = workplace.ID
			})
			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.Role = tt.Role
				v.OfficeID = office.ID
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
			})

			var err error
			c.Request, err = http.NewRequest("GET", ui.WorkEntryPath, nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				require.Equal(t, http.StatusOK, w.Code)
				var res []rdb.GetWorkEntriesByOfficeRow
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				require.NotEmpty(t, res)
				var q rdb.GetWorkEntriesByOfficeRow
				q.ID = created.ID
				q.EmployeeID = created.EmployeeID
				q.WorkplaceID = created.WorkplaceID
				q.Date = created.Date
				q.Hours = created.Hours
				q.StartTime = created.StartTime
				q.EndTime = created.EndTime
				q.Attendance = created.Attendance
				q.Comment = created.Comment
				q.EmployeeName = employee.Name
				q.WorkplaceName = workplace.Name
				q.UpdatedAt = created.UpdatedAt
				q.CreatedAt = created.CreatedAt
				require.Contains(t, res, q)
			}
		})
	}
}

func TestGetWorkEntriesByWorkplace(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role        rdb.UserType
		OtherOffice bool
		OtherWp     bool
		WantErr     bool
	}{
		"admin": {
			Role:    rdb.UserTypeAdmin,
			WantErr: false,
		},
		"admin-other-office": {
			Role:        rdb.UserTypeAdmin,
			OtherOffice: true,
			WantErr:     true,
		},
		"manager": {
			Role:    rdb.UserTypeManager,
			WantErr: false,
		},
		"manager-other-wp": {
			Role:    rdb.UserTypeManager,
			OtherWp: true,
			WantErr: true,
		},
		"employee": {
			Role:    rdb.UserTypeEmployee,
			WantErr: true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			office := test.CreateOffice(t, c, dbConn, nil)
			workplace := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
				v.OfficeID = office.ID
			})
			employee := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
				v.WorkplaceID = workplace.ID
			})
			created := test.CreateWorkEntries(t, c, dbConn, func(v *rdb.WorkEntry) {
				v.EmployeeID = employee.ID
				v.WorkplaceID = workplace.ID
			})
			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				if tt.OtherOffice {
					v.OfficeID = test.CreateOffice(t, c, dbConn, nil).ID
				} else {
					v.OfficeID = office.ID
				}
				v.Role = tt.Role
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
				if tt.OtherWp {
					v.EmployeeID = pgtype.Int8{Int64: test.CreateEmployee(t, c, dbConn, nil).ID, Valid: true}
				}
			})

			var err error
			c.Request, err = http.NewRequest("GET", fmt.Sprintf("%sworkplace/%d/", ui.WorkEntryPath, workplace.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				require.Equal(t, http.StatusOK, w.Code)
				var res []rdb.GetWorkEntriesByEmployeeRow
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				require.NotEmpty(t, res)
				require.Equal(t, created.EmployeeID, res[0].EmployeeID)
				require.Equal(t, created.WorkplaceID, res[0].WorkplaceID)
				require.Equal(t, created.Date.Time.Format("2006-01-02T15:04:05.000Z"), res[0].Date.Time.Format("2006-01-02T15:04:05.000Z"))
				if created.Hours.Valid {
					require.Equal(t, int(created.Hours.Int16), int(res[0].Hours.Int16))
				} else if created.Attendance.Valid {
					require.Equal(t, created.Attendance.Bool, res[0].Attendance.Bool)
				} else {
					require.Equal(t, time.UnixMicro(created.StartTime.Microseconds).Format("2006-01-02T15:04:05.000Z"), time.UnixMicro(res[0].StartTime.Microseconds).Format("2006-01-02T15:04:05.000Z"))
					require.Equal(t, time.UnixMicro(created.EndTime.Microseconds).Format("2006-01-02T15:04:05.000Z"), time.UnixMicro(res[0].EndTime.Microseconds).Format("2006-01-02T15:04:05.000Z"))
				}
				if created.Comment.Valid {
					require.Equal(t, created.Comment.String, res[0].Comment.String)
				}
			}
		})
	}
}

func TestGetWorkEntries(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role        rdb.UserType
		OtherOffice bool
		OtherWp     bool
		WantErr     bool
	}{
		"admin": {
			Role:    rdb.UserTypeAdmin,
			WantErr: false,
		},
		"admin-other-office": {
			Role:        rdb.UserTypeAdmin,
			OtherOffice: true,
			WantErr:     true,
		},
		"manager": {
			Role:    rdb.UserTypeManager,
			WantErr: false,
		},
		"manager-other-wp": {
			Role:    rdb.UserTypeManager,
			OtherWp: true,
			WantErr: true,
		},
		"employee": {
			Role:    rdb.UserTypeEmployee,
			WantErr: false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			office := test.CreateOffice(t, c, dbConn, nil)
			workplace := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
				v.OfficeID = office.ID
			})
			employee := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
				v.WorkplaceID = workplace.ID
			})
			created := test.CreateWorkEntries(t, c, dbConn, func(v *rdb.WorkEntry) {
				v.EmployeeID = employee.ID
				v.WorkplaceID = workplace.ID
			})
			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				if tt.OtherOffice {
					v.OfficeID = test.CreateOffice(t, c, dbConn, nil).ID
				} else {
					v.OfficeID = office.ID
				}
				v.Role = tt.Role
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
				if tt.OtherWp {
					v.EmployeeID = pgtype.Int8{Int64: test.CreateEmployee(t, c, dbConn, nil).ID, Valid: true}
				}
			})

			var err error
			c.Request, err = http.NewRequest("GET", fmt.Sprintf("%semployee/%d/", ui.WorkEntryPath, employee.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				require.Equal(t, http.StatusOK, w.Code)
				var res []rdb.GetWorkEntriesByEmployeeRow
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				require.NotEmpty(t, res)
				require.Equal(t, created.EmployeeID, res[0].EmployeeID)
				require.Equal(t, created.WorkplaceID, res[0].WorkplaceID)
				require.Equal(t, created.Date.Time.Format("2006-01-02T15:04:05.000Z"), res[0].Date.Time.Format("2006-01-02T15:04:05.000Z"))
				if created.Hours.Valid {
					require.Equal(t, int(created.Hours.Int16), int(res[0].Hours.Int16))
				} else if created.Attendance.Valid {
					require.Equal(t, created.Attendance.Bool, res[0].Attendance.Bool)
				} else {
					require.Equal(t, time.UnixMicro(created.StartTime.Microseconds).Format("2006-01-02T15:04:05.000Z"), time.UnixMicro(res[0].StartTime.Microseconds).Format("2006-01-02T15:04:05.000Z"))
					require.Equal(t, time.UnixMicro(created.EndTime.Microseconds).Format("2006-01-02T15:04:05.000Z"), time.UnixMicro(res[0].EndTime.Microseconds).Format("2006-01-02T15:04:05.000Z"))
				}
				if created.Comment.Valid {
					require.Equal(t, created.Comment.String, res[0].Comment.String)
				}
			}
		})
	}
}

func TestPostWorkEntry(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		WorkType    rdb.WorkType
		Hours       int
		StartTime   string
		EndTime     string
		Attendance  bool
		Role        rdb.UserType
		OtherOffice bool
		OtherWp     bool
		WantErr     bool
	}{
		"admin": {
			WorkType: rdb.WorkTypeHours,
			Hours:    rand.Intn(23) + 1,
			Role:     rdb.UserTypeAdmin,
			WantErr:  false,
		},
		"admin-other-office": {
			WorkType:    rdb.WorkTypeHours,
			Hours:       rand.Intn(23) + 1,
			Role:        rdb.UserTypeAdmin,
			OtherOffice: true,
			WantErr:     true,
		},
		"manager": {
			WorkType: rdb.WorkTypeHours,
			Hours:    rand.Intn(23) + 1,
			Role:     rdb.UserTypeManager,
			WantErr:  false,
		},
		"manager-other-wp": {
			WorkType: rdb.WorkTypeHours,
			Hours:    rand.Intn(23) + 1,
			Role:     rdb.UserTypeManager,
			OtherWp:  true,
			WantErr:  true,
		},
		"employee": {
			WorkType: rdb.WorkTypeHours,
			Hours:    rand.Intn(23) + 1,
			Role:     rdb.UserTypeEmployee,
			WantErr:  false,
		},
		"hours": {
			WorkType: rdb.WorkTypeHours,
			Hours:    rand.Intn(23) + 1,
			Role:     rdb.UserTypeAdmin,
			WantErr:  false,
		},
		"time": {
			WorkType:  rdb.WorkTypeTime,
			StartTime: "1970-01-01T08:00:00.000Z",
			EndTime:   "1970-01-01T17:00:00.000Z",
			Role:      rdb.UserTypeAdmin,
			WantErr:   false,
		},
		"attendance": {
			WorkType:   rdb.WorkTypeAttendance,
			Attendance: true,
			Role:       rdb.UserTypeAdmin,
			WantErr:    false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			o := test.CreateOffice(t, c, dbConn, nil)
			wp := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
				v.WorkType = tt.WorkType
				v.OfficeID = o.ID
			})
			e := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
				v.WorkplaceID = wp.ID
			})
			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				if tt.OtherOffice {
					v.OfficeID = test.CreateOffice(t, c, dbConn, nil).ID
				} else {
					v.OfficeID = o.ID
				}
				v.Role = tt.Role
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: e.ID, Valid: true}
				}
				if tt.OtherWp {
					v.EmployeeID = pgtype.Int8{Int64: test.CreateEmployee(t, c, dbConn, nil).ID, Valid: true}
				}
			})

			p := handler.PostWorkEntryParams{
				EmployeeID:  e.ID,
				WorkplaceID: wp.ID,
				Date:        "2006-01-02T00:00:00.000+09:00",
				Hours:       tt.Hours,
				StartTime:   tt.StartTime,
				EndTime:     tt.EndTime,
				Attendance:  tt.Attendance,
				Comment:     "test",
			}
			b, err := json.Marshal(p)
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			require.NoError(t, err)
			c.Request, err = http.NewRequest("POST", ui.WorkEntryPath, body)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			utc, err := time.LoadLocation("UTC")
			require.NoError(t, err)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				require.Equal(t, http.StatusOK, w.Code)
				var res rdb.WorkEntry
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				require.Equal(t, p.EmployeeID, res.EmployeeID)
				require.Equal(t, p.WorkplaceID, res.WorkplaceID)
				require.Equal(t, p.Date, res.Date.Time.Format("2006-01-02T15:04:05.000+09:00"))
				if res.Hours.Valid {
					require.Equal(t, p.Hours, int(res.Hours.Int16))
				} else if res.Attendance.Valid {
					require.Equal(t, p.Attendance, res.Attendance.Bool)
				} else {
					require.Equal(t, p.StartTime, time.UnixMicro(res.StartTime.Microseconds).In(utc).Format("2006-01-02T15:04:05.000Z"))
					require.Equal(t, p.EndTime, time.UnixMicro(res.EndTime.Microseconds).In(utc).Format("2006-01-02T15:04:05.000Z"))
				}
				if res.Comment.Valid {
					require.Equal(t, p.Comment, res.Comment.String)
				}

				t.Cleanup(func() {
					require.NoError(t, rdb.New(dbConn).TestDeleteWorkEntry(c, res.ID))
				})
			}
		})
	}
}

func TestDeleteWorkEntry(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role        rdb.UserType
		OtherOffice bool
		OtherWp     bool
		WantErr     bool
	}{
		"admin": {
			Role:    rdb.UserTypeAdmin,
			WantErr: false,
		},
		"admin-other-office": {
			Role:        rdb.UserTypeAdmin,
			OtherOffice: true,
			WantErr:     true,
		},
		"manager": {
			Role:    rdb.UserTypeManager,
			WantErr: false,
		},
		"manager-other-wp": {
			Role:    rdb.UserTypeManager,
			OtherWp: true,
			WantErr: true,
		},
		"employee": {
			Role:    rdb.UserTypeEmployee,
			WantErr: false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			office := test.CreateOffice(t, c, dbConn, nil)
			workplace := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
				v.OfficeID = office.ID
			})
			employee := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
				v.WorkplaceID = workplace.ID
			})
			created := test.CreateWorkEntries(t, c, dbConn, func(v *rdb.WorkEntry) {
				v.EmployeeID = employee.ID
				v.WorkplaceID = workplace.ID
			})
			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				if tt.OtherOffice {
					v.OfficeID = test.CreateOffice(t, c, dbConn, nil).ID
				} else {
					v.OfficeID = office.ID
				}
				v.Role = tt.Role
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
				if tt.OtherWp {
					v.EmployeeID = pgtype.Int8{Int64: test.CreateEmployee(t, c, dbConn, nil).ID, Valid: true}
				}
			})

			var err error
			c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s%d/", ui.WorkEntryPath, created.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				require.Equal(t, http.StatusNoContent, w.Code)
				require.NotEmpty(t, test.GetDeletedAtWorkEntry(t, c, dbConn, created.ID))
			}
		})
	}
}
