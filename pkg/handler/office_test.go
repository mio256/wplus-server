package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/mio256/wplus-server/pkg/infra"
	"github.com/mio256/wplus-server/pkg/infra/rdb"
	"github.com/mio256/wplus-server/pkg/test"
	"github.com/mio256/wplus-server/pkg/ui"
	"github.com/mio256/wplus-server/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOffices(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateOffice(t, c, dbConn, nil)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", ui.OfficePath, nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res []rdb.Office
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res)
	assert.Contains(t, res, *created)
}

func TestPostOffice(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	var p struct {
		Name string `json:"name"`
	}
	p.Name = faker.Username()
	b, err := json.Marshal(p)
	require.NoError(t, err)
	body := bytes.NewBuffer(b)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("POST", ui.OfficePath, body)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res rdb.Office
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, p.Name, res.Name)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteOffice(c, res.ID))
	})
}

func TestDeleteOffice(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateOffice(t, c, dbConn, nil)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.OfficePath, created.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, test.CheckDeletedOffice(t, c, dbConn, created.ID))
}
