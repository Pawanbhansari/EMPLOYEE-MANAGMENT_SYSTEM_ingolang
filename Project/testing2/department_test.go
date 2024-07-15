package testing2

import (
	"Project/Employee"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateDepartment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO department").
			WithArgs("Engineering").
			WillReturnResult(sqlmock.NewResult(1, 1))

		newDept := Employee.Department{
			Dept_ID:   1, // Ensure Dept_ID is initialized
			Dept_Name: "Engineering",
		}

		jsonValue, _ := json.Marshal(newDept)
		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		expected := `{"dept_id":1,"Dept_Name":"Engineering"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("invalid input", func(t *testing.T) {
		newDept := Employee.Department{
			// Missing Dept_Name
		}

		jsonValue, _ := json.Marshal(newDept)
		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expected := `{"error":"Invalid input: Department name is required"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO department").
			WithArgs("Finance").
			WillReturnError(sql.ErrConnDone)

		newDept := Employee.Department{
			Dept_ID:   1, // Ensure Dept_ID is initialized
			Dept_Name: "Finance",
		}

		jsonValue, _ := json.Marshal(newDept)
		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"Failed to create department: driver: bad connection"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("duplicate department name", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO department").
			WithArgs("Human Resources").
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'Human Resources' for key 'dept_name'"})

		newDept := Employee.Department{
			Dept_ID:   1, // Ensure Dept_ID is initialized
			Dept_Name: "Human Resources",
		}

		jsonValue, _ := json.Marshal(newDept)
		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		expected := `{"error":"Department name already exists"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestGetDepartment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"DEPT_ID", "DEPT_NAME"}).
			AddRow(1, "Human Resources").
			AddRow(2, "Engineering")

		mock.ExpectQuery("SELECT DEPT_ID, DEPT_NAME FROM department").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/departments", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.GET("/departments", repo.GetDepartments)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expected := `[{"dept_id":1,"Dept_Name":"Human Resources"},{"dept_id":2,"Dept_Name":"Engineering"}]`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("no departments found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"DEPT_ID", "DEPT_NAME"})

		mock.ExpectQuery("SELECT DEPT_ID, DEPT_NAME FROM department").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/departments", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.GET("/departments", repo.GetDepartments)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expected := `[]`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DEPT_ID, DEPT_NAME FROM department").
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/departments", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.GET("/departments", repo.GetDepartments)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"driver: bad connection"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
