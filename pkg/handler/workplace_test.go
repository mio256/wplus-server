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
	c.Request, _ = http.NewRequest("POST", ui.WorkplacePath, body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 200, w.Code)
	var res rdb.Workplace
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, p.Name, res.Name)
	assert.Equal(t, p.OfficeID, res.OfficeID)
	assert.Equal(t, p.WorkType, res.WorkType)
	assert.Empty(t, res.DeletedAt)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteOffice(c, res.ID))
		require.NoError(t, rdb.New(dbConn).TestDeleteWorkplace(c, res.ID))
	})
}

func TestDeleteWorkplace(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateWorkplace(t, c, dbConn, nil)

	c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.WorkplacePath, created.ID), nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, test.CheckDeletedWorkplace(t, c, dbConn, created.ID))

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteWorkplace(c, created.ID))
	})
}
