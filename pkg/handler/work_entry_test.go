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

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
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

	tests := map[string]struct {
		Hours      int
		StartTime  string
		EndTime    string
		Attendance bool
	}{
		"hours": {
			Hours: rand.Intn(23) + 1,
		},
		"time": {
			StartTime: "1970-01-01T08:00:00.000000Z",
			EndTime:   "1970-01-01T17:00:00.000000Z",
		},
		"attendance": {
			Attendance: true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			dbConn := infra.ConnectDB(c)

			wp := test.CreateWorkplace(t, c, dbConn, nil)
			e := test.CreateEmployee(t, c, dbConn, nil)

			var p = struct {
				EmployeeID  int64  `json:"employee_id"`
				WorkplaceID int64  `json:"workplace_id"`
				Date        string `json:"date"`
				Hours       int    `json:"hours"`
				StartTime   string `json:"start_time"`
				EndTime     string `json:"end_time"`
				Attendance  bool   `json:"attendance"`
				Comment     string `json:"comment"`
			}{
				EmployeeID:  e.ID,
				WorkplaceID: wp.ID,
				Date:        "2006-01-02T00:00:00.000000Z",
				Hours:       tt.Hours,
				StartTime:   tt.StartTime,
				EndTime:     tt.EndTime,
				Attendance:  tt.Attendance,
				Comment:     "test",
			}
			b, err := json.Marshal(p)
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			user, _ := test.CreateUser(t, c, dbConn, nil)
			token, err := util.GenerateToken(uint64(user.ID))
			require.NoError(t, err)
			c.Request, err = http.NewRequest("POST", ui.WorkEntryPath, body)
			require.NoError(t, err)
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			router.ServeHTTP(w, c.Request)

			utc, err := time.LoadLocation("UTC")
			require.NoError(t, err)

			require.Equal(t, 200, w.Code)
			var res rdb.WorkEntry
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			require.Equal(t, p.EmployeeID, res.EmployeeID)
			require.Equal(t, p.WorkplaceID, res.WorkplaceID)
			require.Equal(t, p.Date, res.Date.Time.Format("2006-01-02T15:04:05.000000Z"))
			if res.Hours.Valid {
				require.Equal(t, p.Hours, int(res.Hours.Int16))
			} else if res.Attendance.Valid {
				require.Equal(t, p.Attendance, res.Attendance.Bool)
			} else {
				require.Equal(t, p.StartTime, time.UnixMicro(res.StartTime.Microseconds).In(utc).Format("2006-01-02T15:04:05.000000Z"))
				require.Equal(t, p.EndTime, time.UnixMicro(res.EndTime.Microseconds).In(utc).Format("2006-01-02T15:04:05.000000Z"))
			}
			if res.Comment.Valid {
				require.Equal(t, p.Comment, res.Comment.String)
			}

			t.Cleanup(func() {
				require.NoError(t, rdb.New(dbConn).TestDeleteWorkEntry(c, res.ID))
			})
		})
	}
}

func TestDeleteWorkEntry(t *testing.T) {
	router := ui.SetupRouter()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	dbConn := infra.ConnectDB(c)

	created := test.CreateWorkEntries(t, c, dbConn, nil)

	user, _ := test.CreateUser(t, c, dbConn, nil)
	token, err := util.GenerateToken(uint64(user.ID))
	require.NoError(t, err)
	c.Request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ui.WorkEntryPath, created.ID), nil)
	require.NoError(t, err)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, c.Request)

	require.Equal(t, 204, w.Code)
	require.NotEmpty(t, test.CheckDeletedWorkEntry(t, c, dbConn, created.ID))
}
