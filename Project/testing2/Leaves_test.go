package testing2

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateLeaveType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid leave type ID and payload", func(t *testing.T) {
		payload := `{"LeaveTypeName":"Personal Leave"}`

		mock.ExpectExec("UPDATE leave_types SET LEAVE_TYPE_NAME = ? WHERE LEAVE_TYPE_ID = ?").
			WithArgs("Personal Leave", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest(http.MethodPut, "/leave-type/1", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := gin.Default()
		router.PUT("/leave-type/:id", repo.UpdateLeaveType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"LeavetypeId":1,"LeaveTypeName":"Personal Leave"}`, w.Body.String())
	})

	t.Run("invalid leave type ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/leave-type/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.PUT("/leave-type/:id", repo.UpdateLeaveType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid leave type ID"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		payload := `{"LeaveTypeName":"Personal Leave"}`

		mock.ExpectExec("UPDATE leave_types SET LEAVE_TYPE_NAME = ? WHERE LEAVE_TYPE_ID = ?").
			WithArgs("Personal Leave", 2).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodPut, "/leave-type/2", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := gin.Default()
		router.PUT("/leave-type/:id", repo.UpdateLeaveType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
func TestDeleteLeaveType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid leave type ID", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM leave_types WHERE LEAVE_TYPE_ID = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest(http.MethodDelete, "/leave-type/1", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/leave-type/:id", repo.DeleteLeaveType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Leave type deleted"}`, w.Body.String())
	})

	t.Run("invalid leave type ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/leave-type/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/leave-type/:id", repo.DeleteLeaveType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid leave type ID"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM leave_types WHERE LEAVE_TYPE_ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodDelete, "/leave-type/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/leave-type/:id", repo.DeleteLeaveType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
func TestGetLeavesTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get all leaves", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMP_ID", "START_DATE", "END_DATE", "LEAVE_TYPE_ID", "APPROVAL_STATUS", "APPROVAL_BY"}).
			AddRow(1, "2024-07-15", "2024-07-20", 1, true, 2).
			AddRow(2, "2024-08-01", "2024-08-05", 2, false, nil)

		mock.ExpectQuery("SELECT EMP_ID, START_DATE, END_DATE, LEAVE_TYPE_ID, APPROVAL_STATUS, APPROVAL_BY FROM leaves").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/leaves", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leaves", repo.GetLeaves)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"EmpId":1,"StartDate":"2024-07-15","EndDate":"2024-07-20","LeaveType_id":1,"Approval_status":true,"ApprovedBy":2},{"EmpId":2,"StartDate":"2024-08-01","EndDate":"2024-08-05","LeaveType_id":2,"Approval_status":false}]`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EMP_ID, START_DATE, END_DATE, LEAVE_TYPE_ID, APPROVAL_STATUS, APPROVAL_BY FROM leaves").
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/leaves", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leaves", repo.GetLeaves)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
func TestGetLeaveTypeByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid leave type ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"LEAVE_TYPE_ID", "LEAVE_TYPE_NAME"}).
			AddRow(1, "Vacation")

		mock.ExpectQuery("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types WHERE LEAVE_TYPE_ID = ?").
			WithArgs(1).
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/leave-type/1", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leave-type/:id", repo.GetLeaveTypeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"LeavetypeId":1,"LeaveTypeName":"Vacation"}`, w.Body.String())
	})

	t.Run("invalid leave type ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/leave-type/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leave-type/:id", repo.GetLeaveTypeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid leave type ID"}`, w.Body.String())
	})

	t.Run("leave type not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types WHERE LEAVE_TYPE_ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		req, _ := http.NewRequest(http.MethodGet, "/leave-type/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leave-type/:id", repo.GetLeaveTypeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Leave type not found"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types WHERE LEAVE_TYPE_ID = ?").
			WithArgs(3).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/leave-type/3", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leave-type/:id", repo.GetLeaveTypeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
