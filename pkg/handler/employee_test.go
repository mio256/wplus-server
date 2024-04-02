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

func TestGetEmployeesByOffice(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	office := test.CreateOffice(t, c, dbConn, nil)
	wp := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
		v.OfficeID = office.ID
	})
	e := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
		v.Name = faker.Username()
		v.WorkplaceID = wp.ID
	})

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s/office/%d", ui.EmployeePath, office.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res []rdb.Employee
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res)
	assert.Equal(t, e.Name, res[0].Name)
	assert.Equal(t, e.WorkplaceID, res[0].WorkplaceID)
}

func TestGetEmployees(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	wp := test.CreateWorkplace(t, c, dbConn, nil)
	e := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
		v.Name = faker.Username()
		v.WorkplaceID = wp.ID
	})

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s/workplace/%d", ui.EmployeePath, wp.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res []rdb.Employee
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res)
	assert.Equal(t, e.Name, res[0].Name)
	assert.Equal(t, e.WorkplaceID, res[0].WorkplaceID)
}

func TestGetEmployee(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	wp := test.CreateWorkplace(t, c, dbConn, nil)
	e := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
		v.Name = faker.Username()
		v.WorkplaceID = wp.ID
	})

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s/%d", ui.EmployeePath, e.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res rdb.Employee
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res)
	assert.Equal(t, e.Name, res.Name)
	assert.Equal(t, e.WorkplaceID, res.WorkplaceID)
}

func TestPostEmployee(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	wp := test.CreateWorkplace(t, c, dbConn, nil)

	var p = rdb.CreateEmployeeParams{
		Name:        faker.Username(),
		WorkplaceID: wp.ID,
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	body := bytes.NewBuffer(b)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("POST", ui.EmployeePath, body)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, w.Code, 200)
	var res rdb.Employee
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, res.Name, p.Name)
	assert.Equal(t, res.WorkplaceID, p.WorkplaceID)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteEmployee(c, res.ID))
	})
}

func TestChangeEmployeeWorkplace(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	wp1 := test.CreateWorkplace(t, c, dbConn, nil)
	wp2 := test.CreateWorkplace(t, c, dbConn, func(v *rdb.Workplace) {
		v.OfficeID = wp1.OfficeID
	})
	e := test.CreateEmployee(t, c, dbConn, func(v *rdb.Employee) {
		v.Name = faker.Username()
		v.WorkplaceID = wp1.ID
	})

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"workplace_id": %d}`, wp2.ID)))
	c.Request, err = http.NewRequest("PUT", fmt.Sprintf("%s/%d", ui.EmployeePath, e.ID), body)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res rdb.Employee
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res)
	assert.Equal(t, e.Name, res.Name)
	assert.NotEqual(t, wp1.ID, res.WorkplaceID)
	assert.Equal(t, wp2.ID, res.WorkplaceID)
}

func TestDeleteEmployee(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateEmployee(t, c, dbConn, nil)
	entry := test.CreateWorkEntries(t, c, dbConn, func(v *rdb.WorkEntry) {
		v.EmployeeID = created.ID
		v.WorkplaceID = created.WorkplaceID
	})

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.EmployeePath, created.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, test.CheckDeletedEmployee(t, c, dbConn, created.ID))
	assert.NotEmpty(t, test.CheckDeletedWorkEntry(t, c, dbConn, entry.ID))
}
