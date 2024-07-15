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

//package testing2
//
//import (
//"Project/main/functions"
//"database/sql"
//"github.com/DATA-DOG/go-sqlmock"
//"github.com/gin-gonic/gin"
//"github.com/stretchr/testify/assert"
//"net/http"
//"net/http/httptest"
//"testing"
//"time"
//)

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

func TestGetNationalHolidays(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid holidays list", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"HOLIDAY_ID", "HOLIDAY_NAME", "FROM_DATE", "TO_DATE"}).
			AddRow(1, "New Year", "2024-01-01", "2024-01-01").
			AddRow(2, "Independence Day", "2024-07-04", "2024-07-04")

		mock.ExpectQuery("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/holidays", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holidays", repo.GetNationalHolidays)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"HolidayID":1,"Holidayname":"New Year","StartDate":"2024-01-01","EndDate":"2024-01-01"},{"HolidayID":2,"Holidayname":"Independence Day","StartDate":"2024-07-04","EndDate":"2024-07-04"}]`, w.Body.String())
	})

	t.Run("no holidays found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{})

		mock.ExpectQuery("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays").
			WillReturnRows(rows)

		req, _ := http.NewRequest(http.MethodGet, "/holidays", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holidays", repo.GetNationalHolidays)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[]`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		mock.ExpectQuery("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays").
			WillReturnError(sql.ErrConnDone)

		req, _ := http.NewRequest(http.MethodGet, "/holidays", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/holidays", repo.GetNationalHolidays)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"driver: bad connection"}`, w.Body.String())
	})
}
