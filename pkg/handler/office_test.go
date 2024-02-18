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
	c.Request, _ = http.NewRequest("POST", ui.OfficePath, body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, w.Code, 200)
	var res rdb.Office
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, res.Name, p.Name)

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

	c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.OfficePath, created.ID), nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, 204, w.Code)
	assert.NotEmpty(t, test.CheckDeletedOffice(t, c, dbConn, created.ID))

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteOffice(c, created.ID))
	})
}