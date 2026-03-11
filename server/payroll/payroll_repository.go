package payroll

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"abah-go/projects/payment/employees"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// calculate Salary
func CalculateSalary(employeeType string, annualSalary, transportAllowance, feedingAllowance, hourlyRate, hoursWorked, taxRate float64) (float64, error) {
	// Using switch statement 
	// Switch the cases between different employee types
	switch employeeType {
	case "fulltime" : 
		// Add allowances to annual salary first
		// then divide by 12 for monthly
		// then deduct tax
		annualGross := annualSalary + transportAllowance + feedingAllowance
		monthlyGross := annualGross / 12
		monthlyNet := monthlyGross * (1 - taxRate)
		return monthlyNet, nil
	
	case "remote" :
		// Remote employees are paid by hours worked
		gross := hourlyRate * hoursWorked
		net := gross * (1 - taxRate)
		return net, nil

	case "hybrid":
		// Hybrid is like fulltime but no allowances
		monthlyGross := annualSalary / 12
		monthlyNet := monthlyGross * (1 - taxRate)
		return monthlyNet, nil
	default:
		// If we get an unknown type something is
		// wrong with the data — return an error
		return 0, fmt.Errorf("unknown employee type: %s", employeeType)
	}
}


// PayrollResult is the result of processing
// one single employee during a payroll run
type PayrollResult struct {
	EmployeeID    int `json:"employee_id"`
	EmployeeName  string `json:"employee_name"`
	AccountNumber string `json:"account_number"`
	AmountPaid    float64 `json:"amount_paid"`
	Status        string `json:"status"`
	ErrorMessage  string `json:"error-message,omitempty"`
	DurationMS    int64 `json:"duration_ms"`
}

// PayrollSummary is the overall result of
// an entire payroll run — what gets saved
// to payroll_runs and returned to the handler
type PayrollSummary struct {
	RunID        int64 `json:"run_id"`
	TotalPaid    float64 `json:"total_paid"`
	SuccessCount int	`json:"success_count"`
	FailCount    int	`json:"fail_count"`
	Results      []PayrollResult `json:"results"`
	RanAt        time.Time `json:"ran_at"`
}

// Payroll Run History
type PayrollRun struct {
	ID int `json:"id"`
	TotalPaid float64 `json:"total_paid"`
	SuccessCount int `json:"success_count"`
	FailCount int `json:"fail_count"`
	RanAt time.Time 	`json:"ran_at"`
}

// Payroll repository interface
type PayrollRepository interface {
	SaveRun(totalPaid float64, successCount, failCount int) (int64, error)
	SaveLog(runId int64, result PayrollResult) error
	GetHistory() ([]PayrollRun, error)
}

type SQLPayrollRepository struct {
	db *sql.DB
}

func NewSQLPayrollRepository(db *sql.DB) *SQLPayrollRepository {
	return  &SQLPayrollRepository{db: db}
}

// SaveRun inserts one row into payroll_runs
// and returns the new run id
func (r *SQLPayrollRepository) SaveRun(totalPaid float64, successCount, failCount int) (int64, error) {
	// Prepare the query statement
	queryStatement := `INSERT INTO payroll_runs (total_paid, success_count, fail_count) VALUES (?, ?, ?)`

	// execute on insert
	result, err := r.db.Exec(queryStatement, totalPaid, successCount, failCount)
	if err != nil {
		return 0, fmt.Errorf("Unable to run payroll %w", err)
	}

	runId, err := result.LastInsertId()
	if err != nil {
		return  0, fmt.Errorf("could not get payroll run id: %w", err)
	}

	return  runId, nil
}

// SaveLog inserts one row into payroll_logs
// for one employee in a payroll run
func(r *SQLPayrollRepository) SaveLog(runId int64, result PayrollResult) error {
	// Prepare the stament
	queryStatement := `INSERT INTO payroll_logs (
			run_id,
			employee_id,
			account_number,
			amount_paid,
			status,
			error_message,
			duration_ms
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(queryStatement, runId,
		result.EmployeeID,
		result.AccountNumber,
		result.AmountPaid,
		result.Status,
		result.ErrorMessage,
		result.DurationMS,
	)
	if err != nil {
		return fmt.Errorf("could not save payroll log for employee %d: %w", result.EmployeeID, err)
	}

	return nil
}


// ProcessPayroll is your concurrent logic
// adapted to read from DB and write results back
// maxConcurrent controls how many goroutines run at once
func ProcessPayroll(employeeList []employees.Employee, accountRepo bankaccount.AccountRepository, payrollRepository PayrollRepository, maxConcurrent int ) (*PayrollSummary, error) {
	resultChannel := make(chan PayrollResult, len(employeeList))
	limiter := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, employee := range employeeList {
		wg.Add(1)
		go func(emp employees.Employee){
			defer wg.Done()

			// Acquire a slot in the limiter
			// blocks if maxConcurrent goroutines are already running
			limiter <- struct{}{}
			defer func(){<- limiter}()
			startTime := time.Now()

			// Calculate this employee's salary
			salary, err := CalculateSalary(
				emp.EmployeeType,
				emp.AnnualSalary,
				emp.TransportAllowance,
				emp.FeedingAllowance,
				emp.HourlyRate,
				emp.HoursWorked,
				emp.TaxRate,
			)
			if err != nil {
				resultChannel <- PayrollResult{
					EmployeeID:    emp.ID,
					EmployeeName:  emp.Name,
					AccountNumber: emp.AccountNumber,
					AmountPaid:    0,
					Status:        "failed",
					ErrorMessage:  err.Error(),
					DurationMS:    time.Since(startTime).Milliseconds(),
				}
				return
			}

			// Deposit salary into their account in the DB
			err = accountRepo.UpdateBalance(emp.ID, salary)
			if err != nil {
				resultChannel <- PayrollResult{
					EmployeeID:    emp.ID,
					EmployeeName:  emp.Name,
					AccountNumber: emp.AccountNumber,
					AmountPaid:    salary,
					Status:        "failed",
					ErrorMessage:  err.Error(),
					DurationMS:    time.Since(startTime).Milliseconds(),
				}
				return
			}

			// Success
			resultChannel <- PayrollResult{
				EmployeeID:    emp.ID,
				EmployeeName:  emp.Name,
				AccountNumber: emp.AccountNumber,
				AmountPaid:    salary,
				Status:        "success",
				ErrorMessage:  "",
				DurationMS:    time.Since(startTime).Milliseconds(),
			}
		}(employee)
	}

	// Close the channel once all goroutines finish
	// Same pattern as your original code
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	
	// Collect all results from the channel
	var results []PayrollResult
	var totalPaid float64
	var successCount, failCount int

	for result := range resultChannel {
		results = append(results, result)
		if result.Status == "success" {
			successCount++
			totalPaid += result.AmountPaid
		} else {
			failCount++
		}
	}

	// Save the run summary to payroll_runs table
	runID, err := payrollRepository.SaveRun(totalPaid, successCount, failCount)
	if err != nil {
		return nil, fmt.Errorf("could not save payroll run: %w", err)
	}

	// Save one log row per employee to payroll_logs table
	for _, result := range results {
		if err := payrollRepository.SaveLog(runID, result); err != nil {
			return nil, err
		}
	}

	return &PayrollSummary{
		RunID:        runID,
		TotalPaid:    totalPaid,
		SuccessCount: successCount,
		FailCount:    failCount,
		Results:      results,
		RanAt:        time.Now(),
	}, nil

}

func (r *SQLPayrollRepository) GetHistory() ([]PayrollRun, error) {
	// Prepare statemet
	queryStatement := `SELECT id, total_paid, success_count, fail_count, ran_at FROM payroll_runs ORDER BY ran_at DESC`

	rows, err := r.db.Query(queryStatement)
	if err != nil {
		return nil, fmt.Errorf("could not fetch payroll history: %w", err)
	}
	defer rows.Close()

	var runs []PayrollRun
	for rows.Next() {
		var run PayrollRun 
		err := rows.Scan(
			&run.ID,
			&run.TotalPaid,
			&run.SuccessCount,
			&run.FailCount,
			&run.RanAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan payroll run: %w", err)
		}

		runs = append(runs, run)
	}

	return  runs, nil
}