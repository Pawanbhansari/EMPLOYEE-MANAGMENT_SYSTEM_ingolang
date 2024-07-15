package testing2

import (
	"Project/main/functions"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock the database connection
var db *sql.DB
var mock sqlmock.Sqlmock
var repo *functions.Repo

func init() {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic("failed to open a stub database connection")
	}
	repo = functions.NewRepo(db)
}

func TestGetDepartmentByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid department ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"DEPT_ID", "DEPT_NAME"}).
			AddRow(1, "HR")

		mock.ExpectQuery("SELECT DEPT_ID, DEPT_NAME FROM department WHERE DEPT_ID = ?").
			WithArgs(1).
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/department/1", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/department/:id", repo.GetDepartmentByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"DEPT_ID":1,"DEPT_NAME":"HR"}`, w.Body.String())
	})

	t.Run("invalid department ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/department/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/department/:id", repo.GetDepartmentByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid department ID"}`, w.Body.String())
	})

	t.Run("department not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT DEPT_ID, DEPT_NAME FROM department WHERE DEPT_ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		req, _ := http.NewRequest(http.MethodGet, "/department/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/department/:id", repo.GetDepartmentByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Department not found"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DEPT_ID, DEPT_NAME FROM department WHERE DEPT_ID = ?").
			WithArgs(3).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/department/3", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/department/:id", repo.GetDepartmentByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
