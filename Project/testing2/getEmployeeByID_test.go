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
//var db *sql.DB
//var mock sqlmock.Sqlmock
//var repo *functions.Repo

func init() {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic("failed to open a stub database connection")
	}
	repo = functions.NewRepo(db)
}

func TestGetEmployeeByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid employee ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMPLOYEE_ID", "EMPLOYEE_NAME", "EMPLOYEE_EMAIL", "EMPLOYEE_PHONE", "EMPLOYEE_ADDRESS", "EMPLOYEE_DOB", "DEPT_ID", "MANAGER_ID"}).
			AddRow(1, "John Doe", "john.doe@example.com", "1234567890", "123 Main St", "1990-01-01", 1, 2)

		mock.ExpectQuery("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee WHERE EMPLOYEE_ID = ?").
			WithArgs(1).
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/employee/1", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employee/:id", repo.GetEmployeeByID)
		router.ServeHTTP(w, req)

		expected := `{"EmpID":1,"Name":"John Doe","Email":"john.doe@example.com","Phone":"1234567890","Address":"123 Main St","DOB":"1990-01-01","DeptID":1,"ManagerID":2}`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("invalid employee ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/employee/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employee/:id", repo.GetEmployeeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid employee ID"}`, w.Body.String())
	})

	t.Run("employee not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee WHERE EMPLOYEE_ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		req, _ := http.NewRequest(http.MethodGet, "/employee/2", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employee/:id", repo.GetEmployeeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Employee not found"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee WHERE EMPLOYEE_ID = ?").
			WithArgs(3).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/employee/3", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employee/:id", repo.GetEmployeeByID)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
