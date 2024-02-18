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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	c.Request, _ = http.NewRequest("POST", ui.EmployeePath, body)
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

func TestDeleteEmployee(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateEmployee(t, c, dbConn, nil)

	c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.EmployeePath, created.ID), nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, test.CheckDeletedEmployee(t, c, dbConn, created.ID))

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteEmployee(c, created.ID))
	})
}
