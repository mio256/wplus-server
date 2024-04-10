package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/test"
	"github.com/mio256/wplus-server/pkg/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEmployeesByOffice(t *testing.T) {
	router := ui.SetupRouter()
	ctx := context.Background()
	dbConn := infra.ConnectDB(ctx)

	e1, wp1 := test.CreateEmployeeWithWorkplace(t, ctx, dbConn, nil)
	e2, wp2 := test.CreateEmployeeWithWorkplace(t, ctx, dbConn, nil)

	tests := map[string]struct {
		OursEmployee    *rdb.Employee
		OursWorkplace   *rdb.Workplace
		OthersEmployee  *rdb.Employee
		OthersWorkplace *rdb.Workplace
	}{
		"user1": {
			OursEmployee:    e1,
			OursWorkplace:   wp1,
			OthersEmployee:  e2,
			OthersWorkplace: wp2,
		},
		"user2": {
			OursEmployee:    e2,
			OursWorkplace:   wp2,
			OthersEmployee:  e1,
			OthersWorkplace: wp1,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.OfficeID = tt.OursWorkplace.OfficeID
			})

			var err error
			c.Request, err = http.NewRequest("GET", ui.EmployeePath, nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			assert.Equal(t, http.StatusOK, w.Code)
			var res []rdb.Employee
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.NotEmpty(t, res)
			assert.Contains(t, res, *tt.OursEmployee)
			assert.NotContains(t, res, *tt.OthersEmployee)
		})
	}
}

func TestGetEmployees(t *testing.T) {
	router := ui.SetupRouter()
	ctx := context.Background()
	dbConn := infra.ConnectDB(ctx)

	office := test.CreateOffice(t, ctx, dbConn, nil)
	wp1 := test.CreateWorkplace(t, ctx, dbConn, func(v *rdb.Workplace) {
		v.OfficeID = office.ID
	})
	wp2 := test.CreateWorkplace(t, ctx, dbConn, func(v *rdb.Workplace) {
		v.OfficeID = office.ID
	})
	e1 := test.CreateEmployee(t, ctx, dbConn, func(v *rdb.Employee) {
		v.WorkplaceID = wp1.ID
	})
	e2 := test.CreateEmployee(t, ctx, dbConn, func(v *rdb.Employee) {
		v.WorkplaceID = wp2.ID
	})

	tests := map[string]struct {
		OursEmployee    *rdb.Employee
		OursWorkplace   *rdb.Workplace
		OthersEmployee  *rdb.Employee
		OthersWorkplace *rdb.Workplace
	}{
		"user1": {
			OursEmployee:    e1,
			OursWorkplace:   wp1,
			OthersEmployee:  e2,
			OthersWorkplace: wp2,
		},
		"user2": {
			OursEmployee:    e2,
			OursWorkplace:   wp2,
			OthersEmployee:  e1,
			OthersWorkplace: wp1,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.OfficeID = tt.OursWorkplace.OfficeID
			})

			var err error
			c.Request, err = http.NewRequest("GET", fmt.Sprintf("%sworkplace/%d/", ui.EmployeePath, tt.OursWorkplace.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			assert.Equal(t, http.StatusOK, w.Code)
			var res []rdb.Employee
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.NotEmpty(t, res)
			assert.Contains(t, res, *tt.OursEmployee)
			assert.NotContains(t, res, *tt.OthersEmployee)
		})
	}
}

func TestGetEmployee(t *testing.T) {
	router := ui.SetupRouter()
	ctx := context.Background()
	dbConn := infra.ConnectDB(ctx)

	e1, wp1 := test.CreateEmployeeWithWorkplace(t, ctx, dbConn, nil)
	e2, wp2 := test.CreateEmployeeWithWorkplace(t, ctx, dbConn, nil)

	tests := map[string]struct {
		OfficeID uint64
		Employee *rdb.Employee
	}{
		"ours": {
			OfficeID: uint64(wp1.OfficeID),
			Employee: e1,
		},
		"others": {
			OfficeID: uint64(wp2.OfficeID),
			Employee: e2,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.OfficeID = int64(tt.OfficeID)
			})

			var err error
			c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s%d/", ui.EmployeePath, tt.Employee.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			assert.Equal(t, http.StatusOK, w.Code)
			var res rdb.Employee
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.NotEmpty(t, res)
			assert.Equal(t, res, *tt.Employee)
		})
	}
}

func TestPostEmployee(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role        rdb.UserType
		OtherOffice bool
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
			})

			var p = rdb.CreateEmployeeParams{
				Name:        faker.Username(),
				WorkplaceID: workplace.ID,
			}
			var err error
			b, err := json.Marshal(p)
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			c.Request, err = http.NewRequest("POST", ui.EmployeePath, body)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				var res rdb.Employee
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				assert.Equal(t, res.Name, p.Name)
				assert.Equal(t, res.WorkplaceID, p.WorkplaceID)

				t.Cleanup(func() {
					require.NoError(t, rdb.New(dbConn).TestDeleteEmployee(c, res.ID))
				})
			}
		})
	}
}

func TestChangeEmployeeWorkplace(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role            rdb.UserType
		OtherOffice     bool
		OthersWorkplace bool
		WantErr         bool
	}{
		"admin": {
			Role:    rdb.UserTypeAdmin,
			WantErr: false,
		},
		"admin-others-workplace": {
			Role:            rdb.UserTypeAdmin,
			OthersWorkplace: true,
			WantErr:         true,
		},
		"admin-other-office": {
			Role:        rdb.UserTypeAdmin,
			OtherOffice: true,
			WantErr:     true,
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
			user, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				if tt.OtherOffice {
					v.OfficeID = test.CreateOffice(t, c, dbConn, nil).ID
				} else {
					v.OfficeID = office.ID
				}
				v.Role = tt.Role
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
			})

			newWp := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
				if !tt.OthersWorkplace {
					v.OfficeID = user.OfficeID
				}
			})

			var p = struct {
				WorkplaceID int64 `json:"workplace_id"`
			}{
				WorkplaceID: newWp.ID,
			}
			var err error
			b, err := json.Marshal(p)
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			c.Request, err = http.NewRequest("PUT", fmt.Sprintf("%s%d/", ui.EmployeePath, employee.ID), body)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
				t.Log(w.Body.String())
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				var res rdb.Employee
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				assert.NotEmpty(t, res)
				assert.Equal(t, res.WorkplaceID, newWp.ID)
				assert.Equal(t, res.ID, employee.ID)
				assert.Equal(t, res.Name, employee.Name)

				t.Cleanup(func() {
					require.NoError(t, rdb.New(dbConn).TestDeleteEmployee(c, employee.ID))
				})
			}
		})
	}
}

func TestDeleteEmployee(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role        rdb.UserType
		OtherOffice bool
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
			})

			entry := test.CreateWorkEntries(t, c, dbConn, func(v *rdb.WorkEntry) {
				v.EmployeeID = employee.ID
				v.WorkplaceID = workplace.ID
			})

			var err error
			c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s%d/", ui.EmployeePath, employee.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusNoContent, w.Code)
				assert.NotEmpty(t, test.GetDeletedAtEmployee(t, c, dbConn, employee.ID))
				assert.NotEmpty(t, test.GetDeletedAtWorkEntry(t, c, dbConn, entry.ID))
			}
		})
	}
}
