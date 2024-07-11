package Employee

import (
	"errors"
	"fmt"
	"time"
)

type Employee struct {
	EmpID     int    `json:"emp_id" db:"EMPLOYEE_ID"`
	Name      string `json:"name" db:"EMPLOYEE_NAME"`
	Email     string `json:"email" db:"EMPLOYEE_EMAIL"`
	Phone     int64  `json:"phone" db:"EMPLOYEE_PHONE"`
	Address   string `json:"address" db:"EMPLOYEE_ADDRESS"`
	DOB       string `json:"dob" db:"EMPLOYEE_DOB"`
	DeptID    *int   `json:"dept_id" db:"DEPT_ID"`
	ManagerID *int   `json:"manager_id" db:"MANAGER_ID"`
}

func (e Employee) String() string {
	empStr := fmt.Sprintf("Employee[ID: %d, Name: %s, Email: %s, Phone: %s, Address: %s, DOB: %s,DepartmentID: %d, ManagerID: %d]",
		e.EmpID, e.Name, e.Email, e.Phone, e.Address, e.DOB, *e.DeptID, *e.ManagerID)
	return empStr
}

type Leave struct {
	EmpId           int       `json:"emp_id" db:"EMP_ID"`
	StartDate       time.Time `json:"start_date" db:"START_DATE"`
	EndDate         time.Time `json:"end_date" db:"END_DATE"`
	LeaveType_id    int       `json:"leave_type_id" db:"LEAVE_TYPE_ID"`
	Approval_status *bool     `json:"approval_Status" db:"APPROVAL_STATUS"`
	ApprovedBy      *int      `json:"approved_by,omitempty" db:"APPROVAL_BY"`
}

type Department struct {
	Dept_ID   int    `json:"dept_id" db:"DEPT_ID"`
	Dept_Name string `json:"Dept_Name" db:"DEPT_NAME"`
}
type NationalHolidays struct {
	HolidayID   int64     `json:"holiday_id" db:"HOLIDAY_ID"`
	Holidayname string    `json:"Holiday name" db:"HOLIDAY_NAME"`
	StartDate   time.Time `json:"start_date" db:"FROM_DATE"`
	EndDate     time.Time `json:"end_date" db:"TO_DATE"`
}
type LeaveType struct {
	LeavetypeId   int    `json:"leave_type_id"  db:"LEAVE_TYPE_ID"`
	LeaveTypeName string `json:"LeaveType name" db:"LEAVE_TYPE_NAME"`
}

var nationalHolidays = []NationalHolidays{
	newNationalHoliday("New Year’s Day", "2024-01-01", "2024-01-01"),
	newNationalHoliday("Eid Al Fitr holiday", "2024-04-08", "2024-04-12"),
	newNationalHoliday("Eid al Adha holiday", "2024-06-15", "2024-06-18"),
	newNationalHoliday("Islamic New Year", "2024-07-07", "2024-07-07"),
	newNationalHoliday("Prophet Muhammad’s (PBUH) birthday", "2024-09-15", "2024-09-15"),
	newNationalHoliday("National Day", "2024-12-02", "2024-12-03"),
}

func newNationalHoliday(name, startDateStr, endDateStr string) NationalHolidays {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		fmt.Printf("Error parsing start date for %s: %v\n", name, err)
		startDate = time.Time{}
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		fmt.Printf("Error parsing end date for %s: %v\n", name, err)
		endDate = time.Time{}
	}
	return NationalHolidays{
		Holidayname: name,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}
func GetNationalHolidays() []NationalHolidays {
	return nationalHolidays
}

type HR struct {
	Employee
	HR_ID int `json:"hr_id" db:"HR_ID"`
}

//	func (hr *HR) ApproveLeave(leave *Leave) error {
//		if leave.Approval_status == true {
//			return errors.New(fmt.Sprintf("\n Leave has already been approved for : %d days by %v id: %d", leave.NoOfDays(), hr.Name, hr.EmpID))
//		}
//		leave.Approval_status = true
//		leave.ApprovedBy = &hr.EmpID
//		fmt.Printf("\n Leave approved for : %d days\n", leave.NoOfDays())
//
//		return nil
//	}
//
//	func (hr *HR) RejectLeave(leave *Leave) error {
//		if leave.Approval_status != true {
//			return errors.New("\n Leave has already been rejected")
//		}
//		leave.Approval_status = false
//		leave.ApprovedBy = nil
//		return nil
//	}
func (hr *HR) ApproveLeave(leave *Leave) error {
	if leave.Approval_status != nil && *leave.Approval_status == true {
		return errors.New(fmt.Sprintf("\n Leave has already been approved for : %d days by %v id: %d", leave.NoOfDays(), hr.Name, hr.EmpID))
	}
	*leave.Approval_status = true
	leave.ApprovedBy = &hr.EmpID
	fmt.Printf("\n Leave approved for : %d days\n", leave.NoOfDays())

	return nil
}
func (hr *HR) RejectLeave(leave *Leave) error {
	if leave.Approval_status != nil && *leave.Approval_status == false {
		return errors.New("\n Leave has already been rejected")
	}
	*leave.Approval_status = false
	leave.ApprovedBy = nil
	return nil
}

func isNationalHoliday(date time.Time) bool {
	for _, nationalHoliday := range nationalHolidays {
		if !date.Before(nationalHoliday.StartDate) && !date.After(nationalHoliday.EndDate) {
			return true
		}
	}
	return false
}
func isWeekend(date time.Time) bool {
	day := date.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func (leave *Leave) NoOfDays() int {
	totalLeaves := 0
	for date := leave.StartDate; !date.After(leave.EndDate); date = date.AddDate(0, 0, 1) {
		if !isNationalHoliday(date) && !isWeekend(date) {
			totalLeaves++
		}
	}
	if totalLeaves == 0 {
		fmt.Println("All Leaves requested for are national holidays")
	}
	//totalLeaves := int(leave.EndDate.Sub(leave.StartDate).Hours()/24) + 1
	return totalLeaves
}

//func SetupRouter() *gin.Engine {
//	router := gin.Default()
//
//	// Define all your routes here
//
//	return router
//}
