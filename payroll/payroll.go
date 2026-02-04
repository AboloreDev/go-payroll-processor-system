package payroll

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"fmt"
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
	TaxDeductions  string
	AccountNumber string
	BankAccount   *bankaccount.BankAccount
}

type HybridEmployee struct {
	Name          string
	AnnualSalary  float64
	TaxDeductions  string
	AccountNumber string
	BankAccount   *bankaccount.BankAccount
}

// Interface
type Payment interface {
	CalculateMonthlySalary() float64
	fmt.Stringer
}

// Functions
func (fulltime FulltimeEmployee) CalculateMonthlySalary() float64 {
	annualGrossSalary := fulltime.AnnualSalary + fulltime.FeedingAllowance + fulltime.TransportationAllowance
	monthlyGrossSalary := annualGrossSalary / 12
	monthlyNetSalary := monthlyGrossSalary * 0.10

	return monthlyNetSalary
}

func (remote RemoteEmployee) CalculateMonthlySalary() float64 {
	grossSalary := remote.HourlyRate * remote.HoursWorked
	netSalary := grossSalary / 0.10

	return netSalary
}

func (hybrid HybridEmployee) CalculateMonthlySalary() float64 {
	monthlyGrossSalary := hybrid.AnnualSalary / 12
	monthlyNetSalary := monthlyGrossSalary * 0.10

	return monthlyNetSalary
}

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
func ProcessPayroll(employees []Payment) {
	fmt.Println("\n======================================")
	fmt.Println("     CRIXUS HOLDINGS MONTHLY PAYROLL REPORT     ")
	fmt.Println("======================================")

	totalPayroll := 0.0

	for i, employee := range employees {
		fmt.Printf("\n[Employee #%d]\n", i+1)
		fmt.Println("────────────────────────────────────────")

		// Print employee details (uses String() method)
		fmt.Println(employee)

		pay := employee.CalculateMonthlySalary()
		fmt.Printf("\n💰 Monthly Salary: $%.2f\n", pay)
		fmt.Println("-----------------------")
		fmt.Println("Depositing the monthly salary into employee bank account")

		var bankaccount *bankaccount.BankAccount
		// Using switch statement to know the employee type before depositing
		switch e := employee.(type) {
		case FulltimeEmployee:
			bankaccount = e.BankAccount
		case RemoteEmployee:
			bankaccount = e.BankAccount
		case HybridEmployee:
			bankaccount = e.BankAccount
		}

		if bankaccount != nil {
			err := bankaccount.Deposit(pay)

			if err != nil {
				fmt.Printf("❌ Error depositing salary: %v\n", err)
			}
		}
		totalPayroll = totalPayroll + pay
	}
	fmt.Println("\n════════════════════════════════════════")
	fmt.Printf("TOTAL MONTHLY PAYROLL: $%.2f\n", totalPayroll)
	fmt.Println("════════════════════════════════════════")
	fmt.Println("✓ Payroll processing completed.")
}
