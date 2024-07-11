package functions

import (
	"Project/Employee"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) GetEmployeeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	emp, err := queryEmployeeByID(r.db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func queryEmployeeByID(db *sql.DB, id int) (Employee.Employee, error) {
	var emp Employee.Employee
	err := db.QueryRow("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee WHERE EMPLOYEE_ID = ?", id).Scan(&emp.EmpID, &emp.Name, &emp.Email, &emp.Phone, &emp.Address, &emp.DOB, &emp.DeptID, &emp.ManagerID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Employee not found")
		}
		return emp, err
	}
	return emp, nil
}
func (r *Repo) GetDepartmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var dept Employee.Department
	err = r.db.QueryRow("SELECT DEPT_ID, DEPT_NAME FROM department WHERE DEPT_ID = ?", id).Scan(&dept.Dept_ID, &dept.Dept_Name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dept)
}
func (r *Repo) GetEmployees(c *gin.Context) {
	rows, err := r.db.Query("SELECT EMPLOYEE_ID, EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID FROM employee")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var employees []Employee.Employee
	for rows.Next() {
		var emp Employee.Employee
		err := rows.Scan(&emp.EmpID, &emp.Name, &emp.Email, &emp.Phone, &emp.Address, &emp.DOB, &emp.DeptID, &emp.ManagerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}

	c.JSON(http.StatusOK, employees)
}
func (r *Repo) CreateEmployee(c *gin.Context) {
	var newEmployee Employee.Employee
	if err := c.ShouldBindJSON(&newEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming newEmployee doesn't have EmpID set (or set to 0)
	newEmployee.EmpID = rand.Intn(1000) + 1

	// Ensure db is not nil
	if r.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is nil"})
		return
	}

	// Execute the SQL insert statement
	_, err := r.db.Exec("INSERT INTO EMPLOYEE (EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newEmployee.Name, newEmployee.Email, newEmployee.Phone, newEmployee.Address, newEmployee.DOB, newEmployee.DeptID, newEmployee.ManagerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newEmployee)
}
func (r *Repo) UpdateEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var updatedEmployee Employee.Employee
	if err := c.ShouldBindJSON(&updatedEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedEmployee.EmpID = id

	_, err = r.db.Exec("UPDATE EMPLOYEE SET EMPLOYEE_NAME = ?, EMPLOYEE_EMAIL = ?, EMPLOYEE_PHONE = ?, EMPLOYEE_ADDRESS = ?, EMPLOYEE_DOB = ?, DEPT_ID = ?, MANAGER_ID = ? WHERE EMPLOYEE_ID = ?",
		updatedEmployee.Name, updatedEmployee.Email, updatedEmployee.Phone, updatedEmployee.Address, updatedEmployee.DOB, updatedEmployee.DeptID, updatedEmployee.ManagerID, updatedEmployee.EmpID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}
func (r *Repo) DeleteEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	_, err = r.db.Exec("DELETE FROM EMPLOYEE WHERE EMPLOYEE_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted"})
}
func (r *Repo) GetDepartments(c *gin.Context) {
	rows, err := r.db.Query("SELECT DEPT_ID, DEPT_NAME FROM department")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var departments []Employee.Department
	for rows.Next() {
		var dept Employee.Department
		err := rows.Scan(&dept.Dept_ID, &dept.Dept_Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		departments = append(departments, dept)
	}

	c.JSON(http.StatusOK, departments)
}

func (r *Repo) CreateDepartment(c *gin.Context) {
	var newDept Employee.Department
	if err := c.ShouldBindJSON(&newDept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := r.db.Exec("INSERT INTO department (DEPT_NAME) VALUES (?)", newDept.Dept_Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	newDept.Dept_ID = int(id)

	c.JSON(http.StatusCreated, newDept)
}
func (r *Repo) UpdateDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var updatedDept Employee.Department
	if err := c.ShouldBindJSON(&updatedDept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedDept.Dept_ID = id

	_, err = r.db.Exec("UPDATE department SET DEPT_NAME = ? WHERE DEPT_ID = ?", updatedDept.Dept_Name, updatedDept.Dept_ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedDept)
}
func (r *Repo) DeleteDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	_, err = r.db.Exec("DELETE FROM department WHERE DEPT_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted"})
}
func (r *Repo) GetNationalHolidays(c *gin.Context) {
	rows, err := r.db.Query("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var holidays []Employee.NationalHolidays
	for rows.Next() {
		var holiday Employee.NationalHolidays
		err := rows.Scan(&holiday.HolidayID, &holiday.Holidayname, &holiday.StartDate, &holiday.EndDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		holidays = append(holidays, holiday)
	}

	c.JSON(http.StatusOK, holidays)
}

func (r *Repo) GetNationalHolidayByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid holiday ID"})
		return
	}

	var holiday Employee.NationalHolidays
	err = r.db.QueryRow("SELECT HOLIDAY_ID, HOLIDAY_NAME, FROM_DATE, TO_DATE FROM national_holidays WHERE HOLIDAY_ID = ?", id).Scan(&holiday.HolidayID, &holiday.Holidayname, &holiday.StartDate, &holiday.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "National holiday not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, holiday)
}

func (r *Repo) CreateNationalHoliday(c *gin.Context) {
	var newHoliday Employee.NationalHolidays
	if err := c.ShouldBindJSON(&newHoliday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := r.db.Exec("INSERT INTO national_holidays (HOLIDAY_NAME, FROM_DATE, TO_DATE) VALUES (?, ?, ?)", newHoliday.Holidayname, newHoliday.StartDate, newHoliday.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	newHoliday.HolidayID = id

	c.JSON(http.StatusCreated, newHoliday)
}

func (r *Repo) UpdateNationalHoliday(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid holiday ID"})
		return
	}

	var updatedHoliday Employee.NationalHolidays
	if err := c.ShouldBindJSON(&updatedHoliday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedHoliday.HolidayID = id

	_, err = r.db.Exec("UPDATE national_holidays SET HOLIDAY_NAME = ?, FROM_DATE = ?, TO_DATE = ? WHERE HOLIDAY_ID = ?", updatedHoliday.Holidayname, updatedHoliday.StartDate, updatedHoliday.EndDate, updatedHoliday.HolidayID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedHoliday)
}

func (r *Repo) DeleteNationalHoliday(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid holiday ID"})
		return
	}

	_, err = r.db.Exec("DELETE FROM national_holidays WHERE HOLIDAY_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "National holiday deleted"})
}
func (r *Repo) GetLeaveTypes(c *gin.Context) {
	rows, err := r.db.Query("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var leaveTypes []Employee.LeaveType
	for rows.Next() {
		var leaveType Employee.LeaveType
		err := rows.Scan(&leaveType.LeavetypeId, &leaveType.LeaveTypeName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		leaveTypes = append(leaveTypes, leaveType)
	}

	c.JSON(http.StatusOK, leaveTypes)
}

func (r *Repo) GetLeaveTypeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave type ID"})
		return
	}

	var leaveType Employee.LeaveType
	err = r.db.QueryRow("SELECT LEAVE_TYPE_ID, LEAVE_TYPE_NAME FROM leave_types WHERE LEAVE_TYPE_ID = ?", id).Scan(&leaveType.LeavetypeId, &leaveType.LeaveTypeName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave type not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, leaveType)
}

func (r *Repo) CreateLeaveType(c *gin.Context) {
	var newLeaveType Employee.LeaveType
	if err := c.ShouldBindJSON(&newLeaveType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := r.db.Exec("INSERT INTO leave_types (LEAVE_TYPE_NAME) VALUES (?)", newLeaveType.LeaveTypeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	newLeaveType.LeavetypeId = int(id)

	c.JSON(http.StatusCreated, newLeaveType)
}

func (r *Repo) UpdateLeaveType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave type ID"})
		return
	}

	var updatedLeaveType Employee.LeaveType
	if err := c.ShouldBindJSON(&updatedLeaveType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedLeaveType.LeavetypeId = id

	_, err = r.db.Exec("UPDATE leave_types SET LEAVE_TYPE_NAME = ? WHERE LEAVE_TYPE_ID = ?", updatedLeaveType.LeaveTypeName, updatedLeaveType.LeavetypeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedLeaveType)
}

func (r *Repo) DeleteLeaveType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave type ID"})
		return
	}

	_, err = r.db.Exec("DELETE FROM leave_types WHERE LEAVE_TYPE_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave type deleted"})
}
func (r *Repo) GetLeaves(c *gin.Context) {
	rows, err := r.db.Query("SELECT EMP_ID, START_DATE, END_DATE, LEAVE_TYPE_ID, APPROVAL_STATUS, APPROVAL_BY FROM leaves")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var leaves []Employee.Leave
	for rows.Next() {
		var leave Employee.Leave
		var approvalStatus sql.NullBool
		var approvedBy sql.NullInt64
		err := rows.Scan(&leave.EmpId, &leave.StartDate, &leave.EndDate, &leave.LeaveType_id, &approvalStatus, &approvedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if approvalStatus.Valid {
			leave.Approval_status = &approvalStatus.Bool
		}
		if approvedBy.Valid {
			approvedByInt := int(approvedBy.Int64)
			leave.ApprovedBy = &approvedByInt
		}
		leaves = append(leaves, leave)
	}

	c.JSON(http.StatusOK, leaves)
}

func (r *Repo) GetLeaveByEmpID(c *gin.Context) {
	empID, err := strconv.Atoi(c.Param("empId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	rows, err := r.db.Query("SELECT EMP_ID, START_DATE, END_DATE, LEAVE_TYPE_ID, APPROVAL_STATUS, APPROVAL_BY FROM leaves WHERE EMP_ID = ?", empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var leaves []Employee.Leave
	for rows.Next() {
		var leave Employee.Leave
		var approvalStatus sql.NullBool
		var approvedBy sql.NullInt64
		err := rows.Scan(&leave.EmpId, &leave.StartDate, &leave.EndDate, &leave.LeaveType_id, &approvalStatus, &approvedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if approvalStatus.Valid {
			leave.Approval_status = &approvalStatus.Bool
		}
		if approvedBy.Valid {
			approvedByInt := int(approvedBy.Int64)
			leave.ApprovedBy = &approvedByInt
		}
		leaves = append(leaves, leave)
	}

	if len(leaves) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No leaves found for this employee"})
		return
	}

	c.JSON(http.StatusOK, leaves)
}

func (r *Repo) CreateLeave(c *gin.Context) {
	var newLeave Employee.Leave
	if err := c.ShouldBindJSON(&newLeave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := r.db.Exec("INSERT INTO leaves (EMP_ID, START_DATE, END_DATE, LEAVE_TYPE_ID, APPROVAL_STATUS) VALUES (?, ?, ?, ?, ?)",
		newLeave.EmpId, newLeave.StartDate, newLeave.EndDate, newLeave.LeaveType_id, newLeave.Approval_status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newLeave)
}
func (r *Repo) GetHR(c *gin.Context) {
	rows, err := r.db.Query("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var hrs []Employee.HR
	for rows.Next() {
		var hr Employee.HR
		err := rows.Scan(&hr.EmpID, &hr.Name, &hr.Email, &hr.Phone, &hr.Address, &hr.DOB, &hr.DeptID, &hr.ManagerID, &hr.HR_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hrs = append(hrs, hr)
	}

	c.JSON(http.StatusOK, hrs)
}

// getHRByID retrieves a specific HR by their HR_ID
func (r *Repo) GetHRByID(c *gin.Context) {
	hrID, err := strconv.Atoi(c.Param("hrId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid HR ID"})
		return
	}

	var hr Employee.HR
	err = r.db.QueryRow("SELECT e.EMPLOYEE_ID, e.EMPLOYEE_NAME, e.EMPLOYEE_EMAIL, e.EMPLOYEE_PHONE, e.EMPLOYEE_ADDRESS, e.EMPLOYEE_DOB, e.DEPT_ID, e.MANAGER_ID, h.HR_ID FROM employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID WHERE h.HR_ID = ?", hrID).
		Scan(&hr.EmpID, &hr.Name, &hr.Email, &hr.Phone, &hr.Address, &hr.DOB, &hr.DeptID, &hr.ManagerID, &hr.HR_ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "HR not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, hr)
}

// createHR creates a new HR personnel
func (r *Repo) CreateHR(c *gin.Context) {
	var newHR Employee.HR
	if err := c.ShouldBindJSON(&newHR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := r.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Insert into employee table
	result, err := tx.Exec("INSERT INTO employee (EMPLOYEE_NAME, EMPLOYEE_EMAIL, EMPLOYEE_PHONE, EMPLOYEE_ADDRESS, EMPLOYEE_DOB, DEPT_ID, MANAGER_ID) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newHR.Name, newHR.Email, newHR.Phone, newHR.Address, newHR.DOB, newHR.DeptID, newHR.ManagerID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	empID, _ := result.LastInsertId()
	newHR.EmpID = int(empID)

	// Insert into hr table
	result, err = tx.Exec("INSERT INTO hr (EMPLOYEE_ID) VALUES (?)", empID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hrID, _ := result.LastInsertId()
	newHR.HR_ID = int(hrID)

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, newHR)
}

// updateHR updates an existing HR personnel
func (r *Repo) UpdateHR(c *gin.Context) {
	hrID, err := strconv.Atoi(c.Param("hrId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid HR ID"})
		return
	}

	var updatedHR Employee.HR
	if err := c.ShouldBindJSON(&updatedHR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = r.db.Exec("UPDATE employee e JOIN hr h ON e.EMPLOYEE_ID = h.EMPLOYEE_ID SET e.EMPLOYEE_NAME = ?, e.EMPLOYEE_EMAIL = ?, e.EMPLOYEE_PHONE = ?, e.EMPLOYEE_ADDRESS = ?, e.EMPLOYEE_DOB = ?, e.DEPT_ID = ?, e.MANAGER_ID = ? WHERE h.HR_ID = ?",
		updatedHR.Name, updatedHR.Email, updatedHR.Phone, updatedHR.Address, updatedHR.DOB, updatedHR.DeptID, updatedHR.ManagerID, hrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedHR.HR_ID = hrID
	c.JSON(http.StatusOK, updatedHR)
}

// deleteHR deletes an HR personnel
func (r *Repo) DeleteHR(c *gin.Context) {
	hrID, err := strconv.Atoi(c.Param("hrId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid HR ID"})
		return
	}

	tx, err := r.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Delete from hr table
	result, err := tx.Exec("DELETE FROM hr WHERE HR_ID = ?", hrID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "HR not found"})
		return
	}

	// Delete from employee table
	_, err = tx.Exec("DELETE FROM employee WHERE EMPLOYEE_ID = (SELECT EMPLOYEE_ID FROM hr WHERE HR_ID = ?)", hrID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "HR deleted successfully"})
}
