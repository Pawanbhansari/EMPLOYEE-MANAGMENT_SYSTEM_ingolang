package testing2

import (
	"Project/Employee"
	"github.com/go-sql-driver/mysql"

	//"Project/main/a"
	"Project/main/functions"
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

func TestCreateEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize mock database and repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := functions.NewRepo(db)

	// Clear mock expectations after the test
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	t.Run("successful creation", func(t *testing.T) {
		// Mock SQL query for successful creation
		mock.ExpectExec("INSERT INTO EMPLOYEE").
			WithArgs("John Doe", "john@example.com", int64(1234567890), "123 Main St", "1990-01-01", 1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Define a newEmployee object with fixed values and dynamic IDs
		newEmployee := Employee.Employee{
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     1234567890,
			Address:   "123 Main St",
			DOB:       "1990-01-01",
			DeptID:    func(i int) *int { return &i }(1),
			ManagerID: func(i int) *int { return &i }(2),
		}

		// Marshal the newEmployee object to JSON
		jsonValue, _ := json.Marshal(newEmployee)

		// Create a new HTTP request with the JSON body
		req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Initialize a new gin router and define the POST endpoint for creating employees
		router := gin.Default()
		router.POST("/employees", repo.CreateEmployee)
		router.ServeHTTP(w, req)

		// Assert that the HTTP response status code is 201 (Created)
		assert.Equal(t, http.StatusCreated, w.Code)
		// Assert that the response body is not empty
		assert.NotEmpty(t, w.Body)

		// Initialize a map to store the parsed JSON response
		var response map[string]interface{}
		// Unmarshal the response body into the response map
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Assert specific fields of the response against expected values
		assert.Equal(t, "John Doe", response["name"])
		assert.Equal(t, "john@example.com", response["email"])
		assert.Equal(t, "123 Main St", response["address"])
		assert.Equal(t, "1990-01-01", response["dob"])
		assert.Equal(t, 1234567890.0, response["phone"])

		// Assert that emp_id is present and non-zero
		empID, ok := response["emp_id"].(float64)
		assert.True(t, ok)
		assert.NotZero(t, empID)
	})

	t.Run("invalid input", func(t *testing.T) {
		// Define a newEmployee object with missing required fields
		newEmployee := Employee.Employee{
			// Missing Name, Email, Phone, Address, DOB, DeptID, ManagerID fields intentionally
		}

		// Marshal the newEmployee object to JSON
		jsonValue, _ := json.Marshal(newEmployee)

		// Create a new HTTP request with the JSON body
		req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Initialize a new gin router and define the POST endpoint for creating employees
		router := gin.Default()
		router.POST("/employees", repo.CreateEmployee)
		router.ServeHTTP(w, req)

		// Assert that the HTTP response status code is 400 (Bad Request)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		// Assert the error message in the response body
		expected := `{"error":"Invalid input: Name, Email, Phone, Address, DOB, DeptID, ManagerID are required fields"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	// Add other test cases as needed (database error, duplicate email, etc.)

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO EMPLOYEE").
			WithArgs("John Doe", "john@example.com", int64(1234567890), "123 Main St", "1990-01-01", 1, 2).
			WillReturnError(sql.ErrConnDone)

		newEmployee := Employee.Employee{
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     1234567890,
			Address:   "123 Main St",
			DOB:       "1990-01-01",
			DeptID:    func(i int) *int { return &i }(1),
			ManagerID: func(i int) *int { return &i }(2),
		}

		jsonValue, _ := json.Marshal(newEmployee)
		req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/employees", repo.CreateEmployee)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"sql: connection is already closed"}`
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("database connection is nil", func(t *testing.T) {
		repoNilDB := functions.NewRepo(nil)

		newEmployee := Employee.Employee{
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     1234567890,
			Address:   "123 Main St",
			DOB:       "1990-01-01",
			DeptID:    func(i int) *int { return &i }(1),
			ManagerID: func(i int) *int { return &i }(2),
		}

		jsonValue, _ := json.Marshal(newEmployee)
		req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/employees", repoNilDB.CreateEmployee)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"Database connection is nil"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestCreateDepartment2(t *testing.T) {
	// Set up the mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repo instance with the mock db
	repo := functions.NewRepo(db)

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	t.Run("successful creation", func(t *testing.T) {
		department := map[string]string{
			"name": "New Department",
		}
		jsonValue, _ := json.Marshal(department)

		mock.ExpectExec("INSERT INTO department").
			WithArgs(department["name"]).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		expected := `{"message":"Department created successfully","id":1}`
		assert.JSONEq(t, expected, w.Body.String())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("duplicate department name", func(t *testing.T) {
		department := map[string]string{
			"name": "Existing Department",
		}
		jsonValue, _ := json.Marshal(department)

		mock.ExpectExec("INSERT INTO department").
			WithArgs(department["name"]).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'Existing Department' for key 'name'"})

		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		expected := `{"error":"Department name already exists"}`
		assert.JSONEq(t, expected, w.Body.String())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		department := map[string]string{
			"name": "Error Department",
		}
		jsonValue, _ := json.Marshal(department)

		mock.ExpectExec("INSERT INTO department").
			WithArgs(department["name"]).
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/departments", repo.CreateDepartment)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"error":"sql: connection is already closed"}`
		assert.JSONEq(t, expected, w.Body.String())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}
func TestCreateHR2(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize mock database and repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := functions.NewRepo(db)

	t.Run("successful creation", func(t *testing.T) {
		// Mock begin transaction
		mock.ExpectBegin()

		// Mock employee insertion
		mock.ExpectExec("INSERT INTO employee").
			WithArgs("John Doe", "john@example.com", int64(1234567890), "123 Main St", "1990-01-01", 1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock hr insertion
		mock.ExpectExec("INSERT INTO hr").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock commit
		mock.ExpectCommit()

		deptID := 1
		managerID := 2
		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Doe",
				Email:     "john@example.com",
				Phone:     1234567890,
				Address:   "123 Main St",
				DOB:       "1990-01-01",
				DeptID:    &deptID,
				ManagerID: &managerID,
			},
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hr", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hr", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Employee.HR
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 1, response.EmpID)
		assert.Equal(t, 1, response.HR_ID)
		assert.Equal(t, "John Doe", response.Name)
		assert.Equal(t, "john@example.com", response.Email)
		assert.Equal(t, int64(1234567890), response.Phone)
		assert.Equal(t, "123 Main St", response.Address)
		assert.Equal(t, "1990-01-01", response.DOB)
		assert.Equal(t, 1, *response.DeptID)
		assert.Equal(t, 2, *response.ManagerID)
	})

	t.Run("invalid input", func(t *testing.T) {
		newHR := Employee.HR{
			// Missing required fields
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hr", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hr", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	t.Run("transaction begin error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		deptID := 1
		managerID := 2
		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Doe",
				Email:     "john@example.com",
				Phone:     1234567890,
				Address:   "123 Main St",
				DOB:       "1990-01-01",
				DeptID:    &deptID,
				ManagerID: &managerID,
			},
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hr", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hr", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to start transaction")
	})

	t.Run("employee insertion error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO employee").
			WithArgs("John Doe", "john@example.com", int64(1234567890), "123 Main St", "1990-01-01", 1, 2).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		deptID := 1
		managerID := 2
		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Doe",
				Email:     "john@example.com",
				Phone:     1234567890,
				Address:   "123 Main St",
				DOB:       "1990-01-01",
				DeptID:    &deptID,
				ManagerID: &managerID,
			},
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hr", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hr", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "sql: connection is already closed")
	})

	t.Run("hr insertion error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO employee").
			WithArgs("John Doe", "john@example.com", int64(1234567890), "123 Main St", "1990-01-01", 1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO hr").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		deptID := 1
		managerID := 2
		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Doe",
				Email:     "john@example.com",
				Phone:     1234567890,
				Address:   "123 Main St",
				DOB:       "1990-01-01",
				DeptID:    &deptID,
				ManagerID: &managerID,
			},
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hr", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hr", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "sql: connection is already closed")
	})

	t.Run("commit error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO employee").
			WithArgs("John Doe", "john@example.com", int64(1234567890), "123 Main St", "1990-01-01", 1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO hr").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		deptID := 1
		managerID := 2
		newHR := Employee.HR{
			Employee: Employee.Employee{
				Name:      "John Doe",
				Email:     "john@example.com",
				Phone:     1234567890,
				Address:   "123 Main St",
				DOB:       "1990-01-01",
				DeptID:    &deptID,
				ManagerID: &managerID,
			},
		}

		jsonValue, _ := json.Marshal(newHR)
		req, _ := http.NewRequest(http.MethodPost, "/hr", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/hr", repo.CreateHR)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to commit transaction")
	})

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
