package testing2

import (
	"Project/main/functions"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"
)

func init() {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic("failed to open a stub database connection")
	}
	repo = functions.NewRepo(db)
}

func TestGetNationalHolidayByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid holiday ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"HOLIDAY_ID", "HOLIDAY_NAME", "FROM_DATE", "TO_DATE"}).
			AddRow(1, "New Year", "2024-01-01", "2024-01-01")

		mock.ExpectQuery("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays WHERE HOLIDAY_ID = ?").
			WithArgs(1).
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/holiday/1", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holiday/:id", repo.GetNationalHolidayByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"HolidayID":1,"Holidayname":"New Year","StartDate":"2024-01-01","EndDate":"2024-01-01"}`, w.Body.String())
	})

	t.Run("invalid holiday ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/holiday/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holiday/:id", repo.GetNationalHolidayByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid holiday ID"}`, w.Body.String())
	})

	t.Run("holiday not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays WHERE HOLIDAY_ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		req, _ := http.NewRequest(http.MethodGet, "/holiday/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holiday/:id", repo.GetNationalHolidayByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"National holiday not found"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays WHERE HOLIDAY_ID = ?").
			WithArgs(3).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/holiday/3", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holiday/:id", repo.GetNationalHolidayByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
func TestUpdateNationalHoliday(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid holiday ID and payload", func(t *testing.T) {
		// Mock request body
		payload := `{"Holidayname":"New Year","StartDate":"2024-01-01","EndDate":"2024-01-02"}`

		mock.ExpectExec("UPDATE national_holidays SET HOLIDAY_NAME = ?, FROM_DATE = ?, TO_DATE = ? WHERE HOLIDAY_ID = ?").
			WithArgs("New Year", "2024-01-01", "2024-01-02", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest(http.MethodPut, "/national-holiday/1", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := gin.Default()
		router.PUT("/national-holiday/:id", repo.UpdateNationalHoliday)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"HolidayID":1,"Holidayname":"New Year","StartDate":"2024-01-01","EndDate":"2024-01-02"}`, w.Body.String())
	})

	t.Run("invalid holiday ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/national-holiday/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.PUT("/national-holiday/:id", repo.UpdateNationalHoliday)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid holiday ID"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectExec("UPDATE national_holidays SET HOLIDAY_NAME = ?, FROM_DATE = ?, TO_DATE = ? WHERE HOLIDAY_ID = ?").
			WithArgs("New Year", "2024-01-01", "2024-01-02", 2).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodPut, "/national-holiday/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.PUT("/national-holiday/:id", repo.UpdateNationalHoliday)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
func TestDeleteNationalHoliday(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid holiday ID", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM national_holidays WHERE HOLIDAY_ID = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest(http.MethodDelete, "/national-holiday/1", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/national-holiday/:id", repo.DeleteNationalHoliday)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"National holiday deleted"}`, w.Body.String())
	})

	t.Run("invalid holiday ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/national-holiday/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/national-holiday/:id", repo.DeleteNationalHoliday)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid holiday ID"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM national_holidays WHERE HOLIDAY_ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodDelete, "/national-holiday/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/national-holiday/:id", repo.DeleteNationalHoliday)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
func TestGetLeaveTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get all leave types", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"LEAVE_TYPE_ID", "LEAVE_TYPE_NAME"}).
			AddRow(1, "Vacation").
			AddRow(2, "Sick Leave")

		mock.ExpectQuery("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/leave-types", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leave-types", repo.GetLeaveTypes)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"LeavetypeId":1,"LeaveTypeName":"Vacation"},{"LeavetypeId":2,"LeaveTypeName":"Sick Leave"}]`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types").
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/leave-types", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/leave-types", repo.GetLeaveTypes)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
