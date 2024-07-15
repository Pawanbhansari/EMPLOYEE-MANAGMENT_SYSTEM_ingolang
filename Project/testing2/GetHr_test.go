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

func init() {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic("failed to open a stub database connection")
	}
	repo = functions.NewRepo(db)
}

func TestGetHR(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid HR records", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMPLOYEE_ID", "EMPLOYEE_NAME", "EMPLOYEE_EMAIL", "EMPLOYEE_PHONE", "EMPLOYEE_ADDRESS", "EMPLOYEE_DOB", "DEPT_ID", "MANAGER_ID", "HR_ID"}).
			AddRow(1, "John HR", "john.hr@example.com", "1234567890", "123 HR St", "1990-01-01", 1, 2, 101).
			AddRow(2, "Jane HR", "jane.hr@example.com", "0987654321", "456 HR St", "1992-02-02", 2, 3, 102)

		mock.ExpectQuery("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/hrs", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs", repo.GetHR)
		router.ServeHTTP(w, req)

		expected := `[{"EmpID":1,"Name":"John HR","Email":"john.hr@example.com","Phone":"1234567890","Address":"123 HR St","DOB":"1990-01-01","DeptID":1,"ManagerID":2,"HR_ID":101},{"EmpID":2,"Name":"Jane HR","Email":"jane.hr@example.com","Phone":"0987654321","Address":"456 HR St","DOB":"1992-02-02","DeptID":2,"ManagerID":3,"HR_ID":102}]`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("no HR records", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMPLOYEE_ID", "EMPLOYEE_NAME", "EMPLOYEE_EMAIL", "EMPLOYEE_PHONE", "EMPLOYEE_ADDRESS", "EMPLOYEE_DOB", "DEPT_ID", "MANAGER_ID", "HR_ID"})

		mock.ExpectQuery("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/hrs", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs", repo.GetHR)
		router.ServeHTTP(w, req)

		expected := `[]`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID").
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/hrs", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs", repo.GetHR)
		router.ServeHTTP(w, req)

		expected := `{"error":"driver: bad connection"}`
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
