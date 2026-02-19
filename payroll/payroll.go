package payroll

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"fmt"
	"sync"
	"time"
)

// Structs
type FulltimeEmployee struct {
	Name                    string
	AnnualSalary            float64
	TransportationAllowance float64
	FeedingAllowance        float64
	TaxDeductions           string
	AccountNumber           string
	BankAccount             *bankaccount.BankAccount
	
}

type RemoteEmployee struct {
	Name          string
	HoursWorked   float64
	HourlyRate    float64
	TaxDeductions string
	AccountNumber string
	BankAccount   *bankaccount.BankAccount
}

type HybridEmployee struct {
	Name          string
	AnnualSalary  float64
	TaxDeductions string
	AccountNumber string
	BankAccount   *bankaccount.BankAccount
}

type PaymentResult struct {
	EmployeeName  any
	Error         error
	AmountPaid    float64
	Duration      time.Duration
	AccountNumber string
}

// Interface
type Payment interface {
	CalculateMonthlySalary() float64
	GetBankAccount() *bankaccount.BankAccount
	fmt.Stringer
	
}

// Functions
func (fulltime FulltimeEmployee) CalculateMonthlySalary() float64 {
	annualGrossSalary := fulltime.AnnualSalary + fulltime.FeedingAllowance + fulltime.TransportationAllowance
	monthlyGrossSalary := annualGrossSalary / 12
	monthlyNetSalary := monthlyGrossSalary * 0.90

	return monthlyNetSalary
}

func (remote RemoteEmployee) CalculateMonthlySalary() float64 {
	grossSalary := remote.HourlyRate * remote.HoursWorked
	netSalary := grossSalary / 0.10

	return netSalary
}

func (hybrid HybridEmployee) CalculateMonthlySalary() float64 {
	monthlyGrossSalary := hybrid.AnnualSalary / 12
	monthlyNetSalary := monthlyGrossSalary * 0.90

	return monthlyNetSalary
}

func (e FulltimeEmployee) GetBankAccount() *bankaccount.BankAccount { return e.BankAccount }
func (e RemoteEmployee) GetBankAccount() *bankaccount.BankAccount   { return e.BankAccount }
func (e HybridEmployee) GetBankAccount() *bankaccount.BankAccount   { return e.BankAccount }


// Stringer Interface for printing
func (e FulltimeEmployee) String() string {
	return fmt.Sprintf(
		"Full-Time Employee: %s\n"+
			"  Annual Salary: $%.2f\n"+
			"  Feeding Allowance: $%.2f\n"+
			"  Transportation Allowance: $%.2f\n"+
			"  Tax Deduction: %s\n"+
			"  Bank Account: %s",
		e.Name, e.AnnualSalary, e.FeedingAllowance,
		e.TransportationAllowance, e.TaxDeductions, e.BankAccount.AccountNumber,
	)
}

func (e HybridEmployee) String() string {
	return fmt.Sprintf(
		"Hybrid Employee: %s\n"+
			"  Annual Salary: $%.2f\n"+
			"  Tax Deduction: %s\n"+
			"  Bank Account: %s",
		e.Name, e.AnnualSalary, e.TaxDeductions, e.BankAccount.AccountNumber,
	)
}

func (e RemoteEmployee) String() string {
	return fmt.Sprintf(
		"Remote Employee: %s\n"+
			"  Hours Worked: %.2f\n"+
			"  Hourly Rate: $%.2f/hr\n"+
			"  Tax Deduction: %s\n"+
			"  Bank Account: %s",
		e.Name, e.HoursWorked, e.HourlyRate, e.TaxDeductions, e.BankAccount.AccountNumber,
	)
}

// Process Payroll
func ProcessPayroll(employees []Payment, maxConcurrent int) {
	fmt.Println("\n======================================")
	fmt.Println("     CRIXUS HOLDINGS MONTHLY PAYROLL REPORT     ")
	fmt.Println("======================================")

	var wg sync.WaitGroup
	resultChannel := make(chan PaymentResult)
	limiter := make(chan struct{}, maxConcurrent)

	for _, employee := range employees {
		wg.Add(1)
		go func(employee Payment) {
			defer wg.Done()
			time.Sleep(time.Duration(employee.CalculateMonthlySalary() * float64(time.Millisecond)))
			limiter <- struct{}{}
			defer func() { <-limiter }()
			startTime := time.Now()
			var err error

			// fmt.Println("────────────────────────────────────────")
			// Print employee details (uses String() method)
			// fmt.Println(employee)

			pay := employee.CalculateMonthlySalary()
			acc := employee.GetBankAccount()
			

			if acc == nil {
				err = fmt.Errorf("employee has no bank account")
			} else {
				err = acc.Deposit(pay)
			}

			
			resultChannel <- PaymentResult{
				EmployeeName: fmt.Sprintf("%v", employee),
				AccountNumber: func() string {
					if acc == nil {
						return "N/A"
					}
					return acc.AccountNumber
				}(),
				AmountPaid: pay,
				Error:      err,
				Duration:   time.Since(startTime),
			}

		}(employee)

	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect results neatly in ONE place (prevents messy concurrent printing)
	totalPayroll := 0.0
	successCount := 0
	failCount := 0

	for result := range resultChannel {
		if result.Error != nil {
			failCount++
			fmt.Printf("❌ FAILED: %s -> %s | $%.2f | %v\n", result.EmployeeName, result.AccountNumber, result.AmountPaid, result.Error.Error())
			continue
		} else {
			successCount++
			totalPayroll = totalPayroll + result.AmountPaid
			fmt.Printf("Finished Paying %s (%v) in %s\n", result.EmployeeName, result.AmountPaid, result.Duration)
		}
	}
	fmt.Println("\n--------------------------------------")
	fmt.Printf("Success: %d | Failed: %d\n", successCount, failCount)
	fmt.Printf("TOTAL MONTHLY PAYROLL PAID: $%.2f\n", totalPayroll)
	fmt.Println("======================================")

}
