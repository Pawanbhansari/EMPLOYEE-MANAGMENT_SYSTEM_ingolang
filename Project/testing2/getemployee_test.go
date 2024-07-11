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

func TestGetEmployees(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid employees", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMPLOYEE_ID", "EMPLOYEE_NAME", "EMPLOYEE_EMAIL", "EMPLOYEE_PHONE", "EMPLOYEE_ADDRESS", "EMPLOYEE_DOB", "DEPT_ID", "MANAGER_ID"}).
			AddRow(1, "John Doe", "john.doe@example.com", "1234567890", "123 Main St", "1990-01-01", 1, 2).
			AddRow(2, "Jane Doe", "jane.doe@example.com", "0987654321", "456 Elm St", "1992-02-02", 2, 3)

		mock.ExpectQuery("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/employees", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employees", repo.GetEmployees)
		router.ServeHTTP(w, req)

		expected := `[{"EmpID":1,"Name":"John Doe","Email":"john.doe@example.com","Phone":"1234567890","Address":"123 Main St","DOB":"1990-01-01","DeptID":1,"ManagerID":2},{"EmpID":2,"Name":"Jane Doe","Email":"jane.doe@example.com","Phone":"0987654321","Address":"456 Elm St","DOB":"1992-02-02","DeptID":2,"ManagerID":3}]`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("no employees", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMPLOYEE_ID", "EMPLOYEE_NAME", "EMPLOYEE_EMAIL", "EMPLOYEE_PHONE", "EMPLOYEE_ADDRESS", "EMPLOYEE_DOB", "DEPT_ID", "MANAGER_ID"})

		mock.ExpectQuery("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/employees", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employees", repo.GetEmployees)
		router.ServeHTTP(w, req)

		expected := `[]`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee").
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/employees", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/employees", repo.GetEmployees)
		router.ServeHTTP(w, req)

		expected := `{"error":"driver: bad connection"}`
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
