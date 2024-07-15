package testing2

import (
	"Project/Employee"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHR(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO employee").
			WithArgs(sqlmock.AnyArg(), "John HR", "john.hr@example.com", 1234567890, "123 HR St", "1990-01-01", 1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO hr").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John HR",
				Email:     "john.hr@example.com",
				Phone:     1234567890,
				Address:   "123 HR St",
				DOB:       "1990-01-01",
				DeptID:    func(i int) *int { return &i }(1),
				ManagerID: func(i int) *int { return &i }(2),
			},
			HR_ID: 1,
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hrs", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hrs", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "id")
	})

	t.Run("invalid input", func(t *testing.T) {
		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name: "John HR",
				// Missing required fields
			},
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hrs", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hrs", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expected := `{"error":"Invalid input"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO employee").
			WithArgs(sqlmock.AnyArg(), "John HR", "john.hr@example.com", "1234567890", "123 HR St", "1990-01-01", 1, 2).
			WillReturnError(sql.ErrConnDone)

		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John HR",
				Email:     "john.hr@example.com",
				Phone:     1234567890,
				Address:   "123 HR St",
				DOB:       "1990-01-01",
				DeptID:    func(i int) *int { return &i }(1),
				ManagerID: func(i int) *int { return &i }(2),
			},
			HR_ID: 1,
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hrs", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hrs", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"driver: bad connection"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestDeleteHR(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectExec("DELETE e, h FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?").
			WithArgs(101).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest(http.MethodDelete, "/hrs/101", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE("/hrs/:hrId", repo.DeleteHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expected := `{"message":"HR deleted successfully"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("hr not found", func(t *testing.T) {
		mock.ExpectExec("DELETE e, h FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		req, _ := http.NewRequest(http.MethodDelete, "/hrs/999", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE("/hrs/:hrId", repo.DeleteHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		expected := `{"error":"HR not found"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("invalid hr id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/hrs/invalid", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE("/hrs/:hrId", repo.DeleteHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expected := `{"error":"Invalid HR ID"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("DELETE e, h FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?").
			WithArgs(101).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodDelete, "/hrs/101", nil)
		w := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE("/hrs/:hrId", repo.DeleteHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"driver: bad connection"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestGetHRByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid HR record", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"EMPLOYEE_ID", "EMPLOYEE_NAME", "EMPLOYEE_EMAIL", "EMPLOYEE_PHONE", "EMPLOYEE_ADDRESS", "EMPLOYEE_DOB", "DEPT_ID", "MANAGER_ID", "HR_ID"}).
			AddRow(1, "John HR", "john.hr@example.com", 1234567890, "123 HR St", "1990-01-01", 1, 2, 101)

		mock.ExpectQuery("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?").
			WithArgs(101).
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/hrs/101", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs/:hrId", repo.GetHRByID)
		router.ServeHTTP(w, req)

		expected := `{"EmpID":1,"Name":"John HR","Email":"john.hr@example.com","Phone":1234567890,"Address":"123 HR St","DOB":"1990-01-01","DeptID":1,"ManagerID":2,"HR_ID":101}`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("HR not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		req, _ := http.NewRequest(http.MethodGet, "/hrs/999", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs/:hrId", repo.GetHRByID)
		router.ServeHTTP(w, req)

		expected := `{"error":"HR not found"}`
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("invalid HR ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/hrs/invalid", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs/:hrId", repo.GetHRByID)
		router.ServeHTTP(w, req)

		expected := `{"error":"Invalid HR ID"}`
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?").
			WithArgs(101).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/hrs/101", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/hrs/:hrId", repo.GetHRByID)
		router.ServeHTTP(w, req)

		expected := `{"error":"driver: bad connection"}`
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestUpdateHR(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectExec("UPDATE employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID").
			WithArgs("John Updated", "john.updated@example.com", int64(9876543210), "456 Updated St", "1991-02-02", sqlmock.AnyArg(), sqlmock.AnyArg(), 101).
			WillReturnResult(sqlmock.NewResult(0, 1))

		updatedHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Updated",
				Email:     "john.updated@example.com",
				Phone:     9876543210,
				Address:   "456 Updated St",
				DOB:       "1991-02-02",
				DeptID:    func(i int) *int { return &i }(2),
				ManagerID: func(i int) *int { return &i }(3),
			},
			HR_ID: 101,
		}

		jsonValue, _ := json.Marshal(updatedHR)
		req, _ := http.NewRequest(http.MethodPut, "/hrs/101", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.PUT("/hrs/:hrId", repo.UpdateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expected := `{"message":"HR updated successfully"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("hr not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID").
			WithArgs("John Updated", "john.updated@example.com", int64(9876543210), "456 Updated St", "1991-02-02", sqlmock.AnyArg(), sqlmock.AnyArg(), 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		updatedHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Updated",
				Email:     "john.updated@example.com",
				Phone:     9876543210,
				Address:   "456 Updated St",
				DOB:       "1991-02-02",
				DeptID:    func(i int) *int { return &i }(2),
				ManagerID: func(i int) *int { return &i }(3),
			},
			HR_ID: 999,
		}

		jsonValue, _ := json.Marshal(updatedHR)
		req, _ := http.NewRequest(http.MethodPut, "/hrs/999", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.PUT("/hrs/:hrId", repo.UpdateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		expected := `{"error":"HR not found"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("invalid input", func(t *testing.T) {
		updatedHR := Employee.HR{
			Employee: Employee.Employee{
				Name: "John Updated",
				// Missing required fields
			},
			HR_ID: 101,
		}

		jsonValue, _ := json.Marshal(updatedHR)
		req, _ := http.NewRequest(http.MethodPut, "/hrs/101", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.PUT("/hrs/:hrId", repo.UpdateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expected := `{"error":"Invalid input"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID").
			WithArgs("John Updated", "john.updated@example.com", int64(9876543210), "456 Updated St", "1991-02-02", sqlmock.AnyArg(), sqlmock.AnyArg(), 101).
			WillReturnError(sql.ErrConnDone)

		updatedHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Updated",
				Email:     "john.updated@example.com",
				Phone:     9876543210,
				Address:   "456 Updated St",
				DOB:       "1991-02-02",
				DeptID:    func(i int) *int { return &i }(2),
				ManagerID: func(i int) *int { return &i }(3),
			},
			HR_ID: 101,
		}

		jsonValue, _ := json.Marshal(updatedHR)
		req, _ := http.NewRequest(http.MethodPut, "/hrs/101", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.PUT("/hrs/:hrId", repo.UpdateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"driver: bad connection"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
