package testing2

import (
	"Project/main/functions"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestDeleteEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid employee ID", func(t *testing.T) {
		empID := 1
		mock.ExpectExec("DELETE FROM EMPLOYEE WHERE EMPLOYEE_ID = ?").
			WithArgs(empID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req, _ := http.NewRequest(http.MethodDelete, "/employee/"+strconv.Itoa(empID), nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/employee/:id", repo.DeleteEmployee)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Employee deleted"}`, w.Body.String())
	})

	t.Run("invalid employee ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/employee/abc", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/employee/:id", repo.DeleteEmployee)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid employee ID"}`, w.Body.String())
	})

	t.Run("employee not found", func(t *testing.T) {
		empID := 2
		mock.ExpectExec("DELETE FROM EMPLOYEE WHERE EMPLOYEE_ID = ?").
			WithArgs(empID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		req, _ := http.NewRequest(http.MethodDelete, "/employee/"+strconv.Itoa(empID), nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/employee/:id", repo.DeleteEmployee)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Employee deleted"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		empID := 3
		mock.ExpectExec("DELETE FROM EMPLOYEE WHERE EMPLOYEE_ID = ?").
			WithArgs(empID).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodDelete, "/employee/"+strconv.Itoa(empID), nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.DELETE("/employee/:id", repo.DeleteEmployee)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "driver: bad connection"}`, w.Body.String())
	})
}
