package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/test"
	"github.com/mio256/wplus-server/pkg/ui"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	router := ui.SetupRouter()

	tests := map[string]struct {
		Role           rdb.UserType
		WantEmployeeID bool
	}{
		"admin": {
			Role:           rdb.UserTypeAdmin,
			WantEmployeeID: false,
		},
		"employee": {
			Role:           rdb.UserTypeEmployee,
			WantEmployeeID: true,
		},
		"manager": {
			Role:           rdb.UserTypeManager,
			WantEmployeeID: true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			employee := test.CreateEmployee(t, c, dbConn, nil)
			user, plain := test.CreateUser(t, c, dbConn, func(v *rdb.User) {
				v.Role = tt.Role
				if tt.Role != rdb.UserTypeAdmin {
					require.True(t, tt.WantEmployeeID)
					v.EmployeeID = pgtype.Int8{Int64: employee.ID, Valid: true}
				}
			})

			require.NotEqual(t, user.Password, plain)
			// NOTE: GeneratePasswordHash はランダムな salt を使うため、hash は毎回異なる
			// hash, err := util.GeneratePasswordHash(plain)
			// require.NoError(t, err)
			// require.Equal(t, user.Password, hash)

			var p = struct {
				OfficeID uint64 `json:"office_id"`
				UserID   uint64 `json:"user_id"`
				Password string `json:"password"`
			}{
				OfficeID: uint64(user.OfficeID),
				UserID:   uint64(user.ID),
				Password: plain,
			}

			t.Log(p)

			b, err := json.Marshal(p)
			require.NoError(t, err)
			t.Log(b)
			body := bytes.NewBuffer(b)
			t.Log(body)
			c.Request, _ = http.NewRequest("POST", ui.LoginPath, body)
			router.ServeHTTP(w, c.Request)

			assert.Equal(t, 200, w.Code)
			var res struct {
				OfficeID   uint64 `json:"office_id"`
				UserID     uint64 `json:"user_id"`
				Name       string `json:"name"`
				Role       string `json:"role"`
				EmployeeID uint64 `json:"employee_id"`
			}
			t.Log(w.Body.Bytes())
			t.Log(json.Unmarshal(w.Body.Bytes(), &res))
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Equal(t, p.OfficeID, res.OfficeID)
			assert.Equal(t, p.UserID, res.UserID)
			assert.Equal(t, user.Name, res.Name)
			assert.Equal(t, user.Role, rdb.UserType(res.Role))
			if tt.WantEmployeeID {
				assert.Equal(t, employee.ID, int64(res.EmployeeID))
			} else {
				assert.Empty(t, res.EmployeeID)
			}

			cookie := w.Header().Get("Set-Cookie")
			re := regexp.MustCompile(`token=([^;]+)`)
			matches := re.FindStringSubmatch(cookie)
			require.Equal(t, 2, len(matches))
			token := matches[1]

			// NOTE: time.Now().Unix() を使っているため、同時刻に実行すると token が同じになる
			token2, err := util.GenerateToken(uint64(user.ID))
			require.NoError(t, err)
			assert.Equal(t, token, token2)
		})
	}
}
