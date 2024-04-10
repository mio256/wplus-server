package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetOffice(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateOffice(t, c, dbConn, nil)

	user, _ := test.CreateUser(t, c, dbConn, func(v *rdb.User) {
		v.OfficeID = created.ID
	})
	token, err := util.GenerateToken(util.UserClaims{
		UserID:   uint64(user.ID),
		OfficeID: uint64(user.OfficeID),
	})
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", ui.OfficePath, nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res rdb.Office
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, *created, res)
}
