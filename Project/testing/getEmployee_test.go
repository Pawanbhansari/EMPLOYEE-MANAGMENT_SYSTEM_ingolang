package testing

import (
	"Project/Employee"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Define all your routes here

	return router
}

func TestGetEmployees(t *testing.T) {
	// Initialize router and set up test HTTP server
	router := SetupRouter()

	// Create a GET request to /employees
	req, err := http.NewRequest("GET", "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve HTTP into the response recorder
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "status code should be OK")

	// Decode the response body into []Employee.Employee
	var employees []Employee.Employee
	if err := json.Unmarshal(rr.Body.Bytes(), &employees); err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	// Assert that at least one employee is returned (assuming test data is seeded)
	assert.NotEmpty(t, employees, "expected employees to be returned")
}
