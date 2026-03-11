package main

import (
	"database/sql"
	"fmt"
	"math/rand"
)

var firstNames = []string{
	"James", "Sarah", "David", "Maria", "Michael",
	"Jennifer", "Carlos", "Aisha", "Hassan", "Emily",
	"Daniel", "Fatima", "Samuel", "Grace", "Robert",
	"Amara", "Kevin", "Zainab", "Patrick", "Olivia",
	"Victor", "Chioma", "Emmanuel", "Linda", "George",
	"Ngozi", "Benjamin", "Adaeze", "Christopher", "Blessing",
}

var lastNames = []string{
	"Martinez", "Chen", "Kim", "Rodriguez", "Williams",
	"O'Brien", "Patel", "Ahmed", "Thompson", "Santos",
	"Johnson", "Okafor", "Smith", "Adeyemi", "Brown",
	"Nwosu", "Wilson", "Ibrahim", "Taylor", "Eze",
	"Anderson", "Musa", "Jackson", "Obi", "Harris",
	"Bello", "Martin", "Chukwu", "Garcia", "Emeka",
}

var employeeTypes = []int{1, 2, 3}

var fulltimeSalaries = []struct {
    annual    float64
    transport float64
    feeding   float64
}{
    {85000, 1500, 3000},
    {72000, 1200, 2500},
    {95000, 1800, 3500},
    {60000, 1000, 2000},
    {110000, 2000, 4000},
    {78000, 1400, 2800},
    {55000, 1000, 1800},
    {120000, 2500, 4500},
    {68000, 1200, 2400},
    {90000, 1600, 3200},
}

var remoteSalaries = []struct {
    hourlyRate  float64
    hoursWorked float64
}{
    {45, 180},
    {60, 160},
    {35, 200},
    {75, 140},
    {50, 175},
    {40, 190},
    {80, 150},
    {55, 165},
    {30, 200},
    {65, 155},
}

var hybridSalaries = []float64{
    72000, 58000, 88000, 64000, 96000,
    70000, 52000, 82000, 76000, 66000,
}

// Employee holds the generated data before inserting into DB
type Employee struct {
	Name               string
	EmployeeTypeID     int
	AnnualSalary       float64
	TransportAllowance float64
	FeedingAllowance   float64
	HourlyRate         float64
	HoursWorked        float64
	TaxRate            float64
}

type SeedRepository interface {
	Setup() error
	Clear() error
	InsertEmployee(e Employee) (int64, error)
	InsertBankAccount(employeeID int64, accountNumber string, name string) error
}

type SQLSeedRepository struct {
	db *sql.DB
}

// NewSQLSeedRepository create new SQLSeedRepository type
func NewSeedRepository(db *sql.DB) *SQLSeedRepository {
	return &SQLSeedRepository{db: db}
}

func (r *SQLSeedRepository) Setup() error {
	_, err := r.db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("could not enable foreign keys: %w", err)
	}

	if err := r.CreateEmployeeTypesTable(); err != nil {
		return err
	}
	if err := r.CreateEmployeesTable(); err != nil {
		return err
	}
	if err := r.CreateBankAccountsTable(); err != nil {
		return err
	}
	if err := r.CreatePayrollRunsTable(); err != nil {
		return err
	}
	if err := r.CreatePayrollLogsTable(); err != nil {
		return err
	}
	if err := r.CreateWithdrawalRunsTable(); err != nil {
		return err
	}
	if err := r.CreateWithdrawalLogsTable(); err != nil {
		return err
	}

	// Seed employee types
	for _, t := range []string{"fulltime", "remote", "hybrid"} {
		_, err := r.db.Exec("INSERT OR IGNORE INTO employee_types (type) VALUES (?)", t)
		if err != nil {
			return fmt.Errorf("could not insert employee type %s: %w", t, err)
		}
	}

	return nil
}

func (r *SQLSeedRepository) CreateEmployeeTypesTable() error {
	_, err := r.db.Exec(employeeTypesSchema)
	if err != nil {
		return fmt.Errorf("could not create employee_types table: %w", err)
	}
	return nil
}

func (r *SQLSeedRepository) CreateEmployeesTable() error {
	_, err := r.db.Exec(employeesSchema)
	if err != nil {
		return fmt.Errorf("could not create employees table: %w", err)
	}
	return nil
}

func (r *SQLSeedRepository) CreateBankAccountsTable() error {
	_, err := r.db.Exec(bankAccountsSchema)
	if err != nil {
		return fmt.Errorf("could not create bank_accounts table: %w", err)
	}
	return nil
}

func (r *SQLSeedRepository) CreatePayrollRunsTable() error {
	_, err := r.db.Exec(payrollRunsSchema)
	if err != nil {
		return fmt.Errorf("could not create payroll_runs table: %w", err)
	}
	return nil
}

func (r *SQLSeedRepository) CreatePayrollLogsTable() error {
	_, err := r.db.Exec(payrollLogsSchema)
	if err != nil {
		return fmt.Errorf("could not create payroll_logs table: %w", err)
	}
	return nil
}

func (r *SQLSeedRepository) CreateWithdrawalRunsTable() error {
	_, err := r.db.Exec(withdrawalRunsSchema)
	if err != nil {
		return fmt.Errorf("could not create withdrawal_runs table: %w", err)
	}
	return nil
}

func (r *SQLSeedRepository) CreateWithdrawalLogsTable() error {
	_, err := r.db.Exec(withdrawalLogsSchema)
	if err != nil {
		return fmt.Errorf("could not create withdrawal_logs table: %w", err)
	}
	return nil
}

// Clear() wipes employees and bank accounts for a fresh seed
func (r *SQLSeedRepository) Clear() error {

	// Delete in order — children first, parents last
	// because child tables reference parent tables via foreign keys

	tables := []string{
		"payroll_logs",    // references payroll_runs + employees
		"withdrawal_logs", // references withdrawal_runs + employees
		"payroll_runs",    // references nothing but has children
		"withdrawal_runs", // references nothing but has children
		"bank_accounts",   // references employees
		"employees",       // parent — deleted last
	}

	for _, table := range tables {
		_, err := r.db.Exec("DELETE FROM " + table)
		if err != nil {
			return fmt.Errorf("could not clear %s: %w", table, err)
		}
	}

	// Reset autoincrement so IDs always start at 1
	sequences := []string{"employees", "bank_accounts", "payroll_runs", "payroll_logs", "withdrawal_runs", "withdrawal_logs"}
	for _, seq := range sequences {
		_, err := r.db.Exec("DELETE FROM sqlite_sequence WHERE name=?", seq)
		if err != nil {
			return fmt.Errorf("could not reset sequence for %s: %w", seq, err)
		}
	}

	fmt.Println("✓ cleared all tables and reset sequences")
	return nil
}

// InsertEmployee inserts one employee and returns their new id
func (r *SQLSeedRepository) InsertEmployee(e Employee) (int64, error) {
	dbStatement := `
		INSERT INTO employees (
			name,
			employee_type_id,
			annual_salary,
			transport_allowance,
			feeding_allowance,
			hourly_rate,
			hours_worked,
			tax_rate
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(dbStatement,
		e.Name,
		e.EmployeeTypeID,
		e.AnnualSalary,
		e.TransportAllowance,
		e.FeedingAllowance,
		e.HourlyRate,
		e.HoursWorked,
		e.TaxRate,
	)
	if err != nil {
		return 0, fmt.Errorf("could not insert employee %s: %w", e.Name, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last insert id: %w", err)
	}

	return id, nil
}

// InsertBankAccount inserts a bank account linked to an employee
func (r *SQLSeedRepository) InsertBankAccount(employeeID int64, accountNumber string, name string) error {
	query := `
		INSERT INTO bank_accounts (
			employee_id,
			account_number,
			account_owner,
			balance
		) VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, employeeID, accountNumber, name, 0.0)
	if err != nil {
		return fmt.Errorf("could not insert bank account for employee %d: %w", employeeID, err)
	}

	return nil
}

func RandomEmployeeGenerator(i int) Employee {
	// Pick a random name from our lists
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	name := firstName + " " + lastName

	// Pick a random employee type (1, 2, or 3)
	employeeTypeID := employeeTypes[rand.Intn(len(employeeTypes))]

	e := Employee{
		Name:           name,
		EmployeeTypeID: employeeTypeID,
		TaxRate:        0.10,
	}

	switch employeeTypeID {
    case 1: // fulltime
        s := fulltimeSalaries[i%len(fulltimeSalaries)]
        e.AnnualSalary = s.annual
        e.TransportAllowance = s.transport
        e.FeedingAllowance = s.feeding

    case 2: // remote
        s := remoteSalaries[i%len(remoteSalaries)]
        e.HourlyRate = s.hourlyRate
        e.HoursWorked = s.hoursWorked

    case 3: // hybrid
        e.AnnualSalary = hybridSalaries[i%len(hybridSalaries)]
    }

	return e
}

