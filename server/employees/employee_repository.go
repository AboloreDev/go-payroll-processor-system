package employees

import (
	"database/sql"
	"fmt"
	"time"
)

// Definr the structure of the employee id and employee
type EmployeeID struct {
	ID   int
	Type string
}

type Employee struct {
	ID                 int `json:"id"`
	Name               string `json:"name"`
	EmployeeType       string `json:"employee_type"`
	AnnualSalary       float64 `json:"annual_salary"`
	TransportAllowance float64	`json:"transport_allowance"`
	FeedingAllowance   float64	`json:"feeding_allowance"`
	HourlyRate         float64 `json:"hourly_rate"`
	HoursWorked        float64 `json:"hours_worked"`
	TaxRate            float64	`json:"tax_rate"`
	AccountNumber      string `json:"account_number"`
	AccountOwner       string	`json:"account_owner"`
	Balance            float64 `json:"balance"`
	CreatedAt          time.Time `json:"created_at"`
}

// Create an interface that determnes what
// operations we can perform on the db
type EmployeeRepository interface {
	GetAll() ([]Employee, error)
	GetById(id int) (*Employee, error) 
	GetCount() (int, error)
}

// define the db struct
type SQLEmployeeRepository struct {
	db *sql.DB
}

func NewSQLEmployeeRepository(db *sql.DB) *SQLEmployeeRepository {
	return &SQLEmployeeRepository{db: db}
}

// GetAll() fetches all employees joined with their
// bank account and employee type
// We need the JOIN because salary calculation needs
// bank account info and type info — not just the employee row
func (r *SQLEmployeeRepository) GetAll() ([]Employee, error) {
	//Prepare query statements
	queryStatement := `
		SELECT
			e.id,
			e.name,
			et.type,
			e.annual_salary,
			e.transport_allowance,
			e.feeding_allowance,
			e.hourly_rate,
			e.hours_worked,
			e.tax_rate,
			ba.account_number,
			ba.account_owner,
			ba.balance,
			e.created_at
		FROM employees e
		JOIN employee_types et ON e.employee_type_id = et.id
		JOIN bank_accounts ba  ON ba.employee_id = e.id
		ORDER BY e.id ASC
	`
	rows, err := r.db.Query(queryStatement)
	if err != nil {
		return nil, fmt.Errorf("could not query employees: %w", err)
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.EmployeeType,
			&e.AnnualSalary,
			&e.TransportAllowance,
			&e.FeedingAllowance,
			&e.HourlyRate,
			&e.HoursWorked,
			&e.TaxRate,
			&e.AccountNumber,
			&e.AccountOwner,
			&e.Balance,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan employee row: %w", err)
		}
		employees = append(employees, e)
	}
	// rows.Err() catches any error that happened
	// while iterating — not just the initial query
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}
	return employees, nil
}

// GetId() get an employee by their id
// Retrus one employee
func (r *SQLEmployeeRepository) GetById(id int) (*Employee, error) {
	// Prepare the query statement
	queryStatement := `
		SELECT
			e.id,
			e.name,
			et.type,
			e.annual_salary,
			e.transport_allowance,
			e.feeding_allowance,
			e.hourly_rate,
			e.hours_worked,
			e.tax_rate,
			ba.account_number,
			ba.account_owner,
			ba.balance,
			e.created_at
		FROM employees e
		JOIN employee_types et ON e.employee_type_id = et.id
		JOIN bank_accounts ba  ON ba.employee_id = e.id
		WHERE e.id = ?
	`

	rows := r.db.QueryRow(queryStatement, id)

	var e Employee
	err := rows.Scan(
		&e.ID,
		&e.Name,
		&e.EmployeeType,
		&e.AnnualSalary,
		&e.TransportAllowance,
		&e.FeedingAllowance,
		&e.HourlyRate,
		&e.HoursWorked,
		&e.TaxRate,
		&e.AccountNumber,
		&e.AccountOwner,
		&e.Balance,
		&e.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not scan employee row: %w", err)
	}

	return &e, nil
}

// GetCount returns the total number of employees
// used by the dashboard to show total headcount
func (r *SQLEmployeeRepository) GetCount() (int, error) {
	var count int
	// prepare db statement
	queryStatement := `SELECT COUNT(*) FROM employees`

	rows := r.db.QueryRow(queryStatement)

	err := rows.Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("could not count employees: %w", err)
	}

	return count, nil
}