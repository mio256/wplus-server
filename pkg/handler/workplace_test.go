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

func TestGetWorkplace(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateWorkplace(t, c, dbConn, nil)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s/%d", ui.WorkplacePath, created.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res rdb.Workplace
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, created.Name, res.Name)
	assert.Equal(t, created.OfficeID, res.OfficeID)
	assert.Equal(t, created.WorkType, res.WorkType)

}

func TestGetWorkplaces(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	o := test.CreateOffice(t, c, dbConn, nil)
	wp := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
		v.OfficeID = o.ID
	})

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s/office/%d", ui.WorkplacePath, o.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res []rdb.Workplace
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res)
	assert.Equal(t, wp.Name, res[0].Name)
	assert.Equal(t, wp.OfficeID, res[0].OfficeID)
	assert.Equal(t, wp.WorkType, res[0].WorkType)
}

func TestPostWorkplace(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	o := test.CreateOffice(t, c, dbConn, nil)

	var p = rdb.CreateWorkplaceParams{
		Name:     faker.Username(),
		OfficeID: o.ID,
		WorkType: rdb.WorkTypeHours,
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	body := bytes.NewBuffer(b)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("POST", ui.WorkplacePath, body)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
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

func TestDeleteWorkplace(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateWorkplace(t, c, dbConn, nil)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.WorkplacePath, created.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, test.CheckDeletedWorkplace(t, c, dbConn, created.ID))
}
