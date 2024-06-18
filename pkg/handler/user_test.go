package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/test"
	"github.com/mio256/wplus-server/pkg/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostUserAndEmployee(t *testing.T) {
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

			var p = struct {
				Name        string `json:"name"`
				WorkplaceID int64  `json:"workplace_id"`
				Role        string `json:"role"`
				Password    string `json:"password"`
			}{
				Name:        faker.Username(),
				WorkplaceID: workplace.ID,
				Role:        "employee",
				Password:    faker.Password(),
			}
			var err error
			b, err := json.Marshal(p)
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			c.Request, err = http.NewRequest("POST", ui.UserPath, body)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			if tt.WantErr {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusCreated, w.Code)
				var res rdb.User
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				assert.Equal(t, res.Name, p.Name)
				assert.Equal(t, res.OfficeID, office.ID)
				assert.Equal(t, res.Role, rdb.UserType(p.Role))
				t.Logf("EmployeeID: %v", res.EmployeeID)
				assert.True(t, res.EmployeeID.Valid)

				t.Cleanup(func() {
					require.NoError(t, rdb.New(dbConn).TestDeleteUser(c, res.ID))
					require.NoError(t, rdb.New(dbConn).TestDeleteEmployee(c, res.EmployeeID.Int64))
				})
			}
		})
	}
}
