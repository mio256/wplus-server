package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestGetWorkEntries(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateWorkEntries(t, c, dbConn, nil)

	token, err := util.GenerateToken(uint64(test.CreateUser(t, c, dbConn, nil).ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("GET", fmt.Sprintf("%s/%d", ui.WorkEntryPath, created.EmployeeID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	require.Equal(t, 200, w.Code)
	var res []rdb.WorkEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	require.NotEmpty(t, res)
	assert.Equal(t, created.EmployeeID, res[0].EmployeeID)
	assert.Equal(t, created.WorkplaceID, res[0].WorkplaceID)
	assert.Equal(t, created.Date, res[0].Date)
	assert.Equal(t, created.Hours, res[0].Hours)
	assert.Equal(t, created.StartTime, res[0].StartTime)
	assert.Equal(t, created.EndTime, res[0].EndTime)
	assert.Equal(t, created.Attendance, res[0].Attendance)
	assert.Equal(t, created.Comment, res[0].Comment)
}

func TestPostWorkEntry(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	wp := test.CreateWorkplace(t, c, dbConn, nil)
	e := test.CreateEmployee(t, c, dbConn, nil)

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var p = rdb.CreateWorkEntryParams{
		EmployeeID:  e.ID,
		WorkplaceID: wp.ID,
		Date:        pgtype.Date{Time: date, Valid: true},
		Hours:       pgtype.Int2{Int16: int16(rand.Int31n(24)), Valid: true},
		StartTime:   pgtype.Time{},
		EndTime:     pgtype.Time{},
		Attendance:  pgtype.Bool{},
		Comment:     pgtype.Text{String: "test", Valid: true},
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	body := bytes.NewBuffer(b)

	token, err := util.GenerateToken(uint64(test.CreateUser(t, c, dbConn, nil).ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("POST", ui.WorkEntryPath, body)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	require.Equal(t, 200, w.Code)
	var res rdb.WorkEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	require.Equal(t, p.EmployeeID, res.EmployeeID)
	require.Equal(t, p.WorkplaceID, res.WorkplaceID)
	require.Equal(t, p.Date, res.Date)
	require.Equal(t, p.Hours, res.Hours)
	require.Equal(t, p.StartTime, res.StartTime)
	require.Equal(t, p.EndTime, res.EndTime)
	require.Equal(t, p.Attendance, res.Attendance)
	require.Equal(t, p.Comment, res.Comment)

	t.Cleanup(func() {
		require.NoError(t, rdb.New(dbConn).TestDeleteWorkEntry(c, res.ID))
	})
}

func TestDeleteWorkEntry(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateWorkEntries(t, c, dbConn, nil)

	token, err := util.GenerateToken(uint64(test.CreateUser(t, c, dbConn, nil).ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.WorkEntryPath, created.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	require.Equal(t, 204, w.Code)
	require.NotEmpty(t, test.CheckDeletedWorkEntry(t, c, dbConn, created.ID))
}
