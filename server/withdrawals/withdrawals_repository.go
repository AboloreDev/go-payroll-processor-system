package withdrawals

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"abah-go/projects/payment/employees"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

//  A single withdrawal result per employee
type WithdrawalResult struct {
	EmployeeId             int `json:"employee_id"`
	EmployeeName           string `json:"employee_name"`
	AccountNumber          string	`json:"account_number"`
	WithdrawnAmount        float64	`json:"withdrawn_amount"`
	BalanceAfterWithdrawal float64	`json:"balance_after_withdrawal"`
	Status                 string	`json:"status"`
	ErrorMessage           string	`json:"error_message,omitempty"`
	DurationMS             int64	`json:"duration_ms"`
}

// Summary of All withdrawal results
type WithdrawalSummary struct {
	WithdrwalRunId int64 `json:"withdrawal_run_id"`
	TotalWithdrawn float64	`json:"total_withdrawn"`
	SuccessCount   int `json:"success_count"`
	FailCount      int	`json:"fail_count"`
	Results        []WithdrawalResult `json:"results"`
	RanAt          time.Time `json:"ran_at"`
}

// A withdrawal run struct used to get history
type WithdrawalRun struct {
	Id int `json:"id"`
	TotalWithdrawn float64	`json:"total_withdrawn"`
	SuccessCount   int `json:"success_count"`
	FailCount      int	`json:"fail_count"`
	RanAt          time.Time `json:"ran_at"`
}

type WithdrawalRepository interface {
	SaveRun(totalWithdrawn float64, successCount, failCount int) (int64, error)
	SaveLog(runId int64, result WithdrawalResult) error
	GetHistory() ([]WithdrawalRun, error)
}

type SQLWithdrawalRepository struct {
	db *sql.DB
}

func NewSQLWithdrawalRepository(db *sql.DB) *SQLWithdrawalRepository {
	return &SQLWithdrawalRepository{db: db}
}

func (r *SQLWithdrawalRepository) SaveRun(totalWithdrawn float64, successCount, failCount int) (int64, error) {
	// Prepare statement
	queryStatement := `INSERT INTO withdrawal_runs (total_withdrawn, success_count, fail_count) VALUES (?, ?, ?)`

	result, err := r.db.Exec(queryStatement, totalWithdrawn, successCount, failCount)
	if err != nil {
		return 0, fmt.Errorf("Failed to run payroll run %w", err)
	}

	runId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get payroll run id: %w", err)
	}
	return runId, nil
}

func (r *SQLWithdrawalRepository) SaveLog(runId int64, result WithdrawalResult) error {
	// Prepare statement
	queryStatement := `INSERT INTO withdrawal_logs (run_id, employee_id, account_number, amount_withdrawn, balance_after, status, error_message , duration_ms)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(queryStatement, runId,
		result.EmployeeId,
		result.AccountNumber,
		result.WithdrawnAmount,
		result.BalanceAfterWithdrawal,
		result.Status,
		result.ErrorMessage,
		result.DurationMS,
	)

	if err != nil {
		return fmt.Errorf("could not save withdrawal log for employee %d: %w", result.EmployeeId, err)
	}

	return nil
}

// generateWithdrawalAmount returns a random amount
// between 10% and 50% of the employee's current balance
// This simulates realistic withdrawal behaviour
func GenerateWithdrawalAmount(balance float64) float64 {
	if balance <= 0 {
		return 0
	}
	// Pick a random percentage between 10% and 50%
	percentage := 0.20 + rand.Float64()*0.30
	amount := balance * percentage
	// Round to 2 decimal places
	return float64(int(amount*100)) / 100
}

// Process withdrawal
// ConcurrentLogic
func ProcessWithdrawal(employeesList []employees.Employee, accountRepo bankaccount.AccountRepository, withdrawalRepo WithdrawalRepository, maxConcurrent int) (*WithdrawalSummary, error) {
	// define the result channel
	resultChannel := make(chan WithdrawalResult, len(employeesList))
	limiter := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, employee := range employeesList {
		wg.Add(1)

		go func(emp employees.Employee) {
			defer wg.Done()

			limiter <- struct{}{}
			defer func() { <-limiter }()
			startTime := time.Now()

			// Fetch the latest balance for this employee
			// We fetch fresh from DB because payroll may have
			// just updated it — we need the real current number
			account, err := accountRepo.GetAccountsById(emp.ID)
			if err != nil {
				resultChannel <- WithdrawalResult{
					EmployeeId:    emp.ID,
					EmployeeName:  emp.Name,
					AccountNumber: emp.AccountNumber,
					Status:        "failed",
					ErrorMessage:  err.Error(),
					DurationMS:    time.Since(startTime).Milliseconds(),
				}
				return
			}

			// No account found for this employee
			if account == nil {
				resultChannel <- WithdrawalResult{
					EmployeeId:    emp.ID,
					EmployeeName:  emp.Name,
					AccountNumber: emp.AccountNumber,
					Status:        "failed",
					ErrorMessage:  "no bank account found",
					DurationMS:    time.Since(startTime).Milliseconds(),
				}
				return
			}

			// Genrate withdrawal amount
			amount := GenerateWithdrawalAmount(account.Balance)
			if amount <= 0 {
				resultChannel <- WithdrawalResult{
					EmployeeId:    emp.ID,
					EmployeeName:  emp.Name,
					AccountNumber: emp.AccountNumber,
					Status:        "failed",
					ErrorMessage:  "balance is zero, nothing to withdraw",
					DurationMS:    time.Since(startTime).Milliseconds(),
				}
				return
			}

			// Deduct the withdrawal from their balance
			// We pass a negative amount to UpdateBalance
			// because UpdateBalance does balance + amount
			// so balance + (-amount) = balance - amount
			err = accountRepo.UpdateBalance(emp.ID, -amount)
			if err != nil {
				resultChannel <- WithdrawalResult{
					EmployeeId:      emp.ID,
					EmployeeName:    emp.Name,
					AccountNumber:   emp.AccountNumber,
					WithdrawnAmount: amount,
					Status:          "failed",
					ErrorMessage:    err.Error(),
					DurationMS:      time.Since(startTime).Milliseconds(),
				}
				return
			}

			// Calculate balance after withdrawal
			balanceAfter := account.Balance - amount

			resultChannel <- WithdrawalResult{
				EmployeeId:             emp.ID,
				EmployeeName:           emp.Name,
				AccountNumber:          emp.AccountNumber,
				WithdrawnAmount:        amount,
				BalanceAfterWithdrawal: balanceAfter,
				Status:                 "success",
				ErrorMessage:           "",
				DurationMS:             time.Since(startTime).Milliseconds(),
			}

		}(employee)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var results []WithdrawalResult
	var totalWithdrawn float64
	var successCount, failCount int

	for result := range resultChannel {
		results = append(results, result)
		if result.Status == "success" {
			successCount++
			totalWithdrawn += result.WithdrawnAmount
		} else {
			failCount++
		}
	}

	// Save the run summary first
	runID, err := withdrawalRepo.SaveRun(totalWithdrawn, successCount, failCount)
	if err != nil {
		return nil, fmt.Errorf("could not save withdrawal run: %w", err)
	}

	// Save one log per employee
	for _, result := range results {
		if err := withdrawalRepo.SaveLog(runID, result); err != nil {
			return nil, err
		}
	}

	return &WithdrawalSummary{
		WithdrwalRunId: runID,
		TotalWithdrawn: totalWithdrawn,
		SuccessCount:   successCount,
		FailCount:      failCount,
		Results:        results,
		RanAt:          time.Now(),
	}, nil

}


// Get withdrawal history 
func (r *SQLWithdrawalRepository) GetHistory() ([]WithdrawalRun, error) {
	// Query statemnet
	queryStatement := `
		SELECT id, total_withdrawn, success_count, fail_count, ran_at
		FROM withdrawal_runs
		ORDER BY ran_at DESC
	`

	rows, err := r.db.Query(queryStatement)
	if err != nil {
		return nil, fmt.Errorf("could not fetch withdrawal history: %w", err)
	}
	defer rows.Close()

	var runs []WithdrawalRun
	for rows.Next() {
		var run WithdrawalRun
		err := rows.Scan(
			&run.Id,
			&run.TotalWithdrawn,
			&run.SuccessCount,
			&run.FailCount,
			&run.RanAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan withdrawal run: %w", err)
		}
		runs = append(runs, run)
	}
	return runs, nil
}