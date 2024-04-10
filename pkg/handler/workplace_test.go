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

func TestGetWorkplaces(t *testing.T) {
	router := ui.SetupRouter()
	ctx := context.Background()
	dbConn := infra.ConnectDB(ctx)

	wp1 := test.CreateWorkplace(t, ctx, dbConn, nil)
	wp2 := test.CreateWorkplace(t, ctx, dbConn, nil)

	tests := map[string]struct {
		Ours   *rdb.Workplace
		Others *rdb.Workplace
	}{
		"user1": {
			Ours:   wp1,
			Others: wp2,
		},
		"user2": {
			Ours:   wp2,
			Others: wp1,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.OfficeID = tt.Ours.OfficeID
			})

			var err error
			c.Request, err = http.NewRequest("GET", ui.WorkplacePath, nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			assert.Equal(t, http.StatusOK, w.Code)
			var res []rdb.Workplace
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.NotEmpty(t, res)
			assert.Contains(t, res, *tt.Ours)
			assert.NotContains(t, res, *tt.Others)
		})
	}
}

func TestGetWorkplace(t *testing.T) {
	router := ui.SetupRouter()
	ctx := context.Background()
	dbConn := infra.ConnectDB(ctx)

	wp1 := test.CreateWorkplace(t, ctx, dbConn, nil)
	wp2 := test.CreateWorkplace(t, ctx, dbConn, nil)

	tests := map[string]struct {
		Ours    *rdb.Workplace
		Target  *rdb.Workplace
		WantErr bool
	}{
		"user1-wp1": {
			Ours:    wp1,
			Target:  wp1,
			WantErr: false,
		},
		"user1-wp2": {
			Ours:    wp1,
			Target:  wp2,
			WantErr: true,
		},
		"user2-wp1": {
			Ours:    wp2,
			Target:  wp1,
			WantErr: true,
		},
		"user2-wp2": {
			Ours:    wp2,
			Target:  wp2,
			WantErr: false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			_, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.OfficeID = tt.Ours.OfficeID
			})

			var err error
			c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s%d/", ui.WorkplacePath, tt.Target.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				var res rdb.Workplace
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				assert.Equal(t, *tt.Ours, res)
			}
		})
	}
}

func TestPostWorkplace(t *testing.T) {
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
			user, token, _ := test.CreateUserWithToken(t, c, dbConn, func(v *rdb.User) {
				v.OfficeID = office.ID
				v.Role = tt.Role
				if tt.Role == rdb.UserTypeManager || tt.Role == rdb.UserTypeEmployee {
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
			})

			var p = rdb.CreateWorkplaceParams{
				Name:     faker.Username(),
				WorkType: rdb.WorkTypeHours,
			}
			if tt.OtherOffice {
				p.OfficeID = test.CreateOffice(t, c, dbConn, nil).ID
			} else {
				p.OfficeID = user.OfficeID
			}
			b, err := json.Marshal(p)
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			c.Request, err = http.NewRequest("POST", ui.WorkplacePath, body)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusCreated, w.Code)
				var res rdb.Workplace
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				assert.Equal(t, p.Name, res.Name)
				assert.Equal(t, p.OfficeID, res.OfficeID)
				assert.Equal(t, p.WorkType, res.WorkType)
				assert.Empty(t, res.DeletedAt)

				t.Cleanup(func() {
					require.NoError(t, rdb.New(dbConn).TestDeleteWorkplace(c, res.ID))
				})
			}
		})
	}
}

func TestDeleteWorkplace(t *testing.T) {
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

			var err error
			c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s%d/", ui.WorkplacePath, workplace.ID), nil)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusNoContent, w.Code)
				assert.NotEmpty(t, test.GetDeletedAtWorkplace(t, c, dbConn, workplace.ID))
			}
		})
	}
}
