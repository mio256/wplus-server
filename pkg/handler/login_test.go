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
		"manager": {
			Role:           rdb.UserTypeManager,
			WantEmployeeID: true,
		},
		"employee": {
			Role:           rdb.UserTypeEmployee,
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
				if tt.WantEmployeeID {
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

			assert.Equal(t, http.StatusOK, w.Code)

			cookie := w.Header().Get("Set-Cookie")
			re := regexp.MustCompile(`token=([^;]+)`)
			matches := re.FindStringSubmatch(cookie)
			require.Equal(t, 2, len(matches))
			token := matches[1]

			t.Log(token)

			// NOTE: time.Now().Unix() を使っているため、同時刻に実行すると token が同じになる
			userClaims := util.UserClaims{
				UserID:   uint64(user.ID),
				OfficeID: uint64(user.OfficeID),
				Name:     user.Name,
				Role:     string(user.Role),
			}
			if tt.WantEmployeeID {
				userClaims.WorkplaceID = uint64(employee.WorkplaceID)
				userClaims.EmployeeID = uint64(employee.ID)
			}
			token2, err := util.GenerateToken(userClaims)
			require.NoError(t, err)
			// NOTE: tokenが一致しているので、すべてのパラメータが同一である
			assert.Equal(t, token, token2)
		})
	}
}
