package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
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
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	user, plain := test.CreateUser(t, c, dbConn, nil)

	require.NotEqual(t, user.Password, plain)
	// GeneratePasswordHash はランダムな salt を使うため、hash は毎回異なる
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

	b, err := json.Marshal(p)
	require.NoError(t, err)
	body := bytes.NewBuffer(b)
	c.Request, _ = http.NewRequest("POST", ui.LoginPath, body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res struct {
		OfficeID uint64 `json:"office_id"`
		UserID   uint64 `json:"user_id"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, p.OfficeID, res.OfficeID)
	assert.Equal(t, p.UserID, res.UserID)
	assert.Equal(t, user.Name, res.Name)
	assert.Equal(t, user.Role, rdb.UserType(res.Role))

	cookie := w.Header().Get("Set-Cookie")
	re := regexp.MustCompile(`token=([^;]+)`)
	matches := re.FindStringSubmatch(cookie)
	require.Equal(t, 2, len(matches))
	token := matches[1]

	// time.Now().Unix() を使っているため、同時刻に実行すると token が同じになる
	token2, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	assert.Equal(t, token, token2)
}
