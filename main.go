package main

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"abah-go/projects/payment/payroll"
	"abah-go/projects/payment/withdrawals"
	"fmt"
)

func init() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  CRIXUS PLC PAYROLL & BANKING SYSTEM   ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

func main() {
	fmt.Println("\n📋 Creating employee bank accounts...")

	employeeBankAccounts := []*bankaccount.BankAccount{
		{
			AccountNumber:  "ACC-001",
			AccountOwner:   "Sarah Martinez",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-002",
			AccountOwner:   "James Chen",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-003",
			AccountOwner:   "Maria Rodriguez",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-004",
			AccountOwner:   "David Kim",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-005",
			AccountOwner:   "Jennifer Williams",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-006",
			AccountOwner:   "Michael O'Brien",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-007",
			AccountOwner:   "Aisha Patel",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-008",
			AccountOwner:   "Carlos Santos",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-009",
			AccountOwner:   "Emily Thompson",
			AccountBalance: 0.0,
		},
		{
			AccountNumber:  "ACC-010",
			AccountOwner:   "Hassan Ahmed",
			AccountBalance: 0.0,
		},
	}

	fmt.Printf("✓ Created %d bank accounts\n", len(employeeBankAccounts))

	fmt.Println("\n📋 Registering employees bank account into the payroll...")

	hybrid := []payroll.HybridEmployee{
		{
			Name:          "Maria Rodriguez",
			AnnualSalary:  72000.0,
			TaxDeductions: "10%",
			BankAccount:   employeeBankAccounts[2],
		},

		{
			Name:          "Emily Thompson",
			AnnualSalary:  82000.0,
			TaxDeductions: "10%",
			BankAccount:   employeeBankAccounts[8],
		},

		{
			Name:          "Michael O'Brien",
			AnnualSalary:  68000.0,
			TaxDeductions: "10%",
			BankAccount:   employeeBankAccounts[5],
		},
	}

	remote := []payroll.RemoteEmployee{
		{
			Name:          "James Chen",
			HoursWorked:   180.0,
			HourlyRate:    45.0,
			TaxDeductions: "10%",
			BankAccount:   employeeBankAccounts[1],
		},
		{
			Name:          "Jennifer Williams",
			HoursWorked:   120.0,
			HourlyRate:    40.0,
			TaxDeductions: "10%",
			BankAccount:   employeeBankAccounts[4],
		},
		{
			Name:          "Carlos Santos",
			HoursWorked:   160.0,
			HourlyRate:    38.0,
			TaxDeductions: "10%",
			BankAccount:   employeeBankAccounts[7],
		},
	}

	fulltime := []payroll.FulltimeEmployee{
		{
			Name:                    "Sarah Martinez",
			AnnualSalary:            85000.0,
			TransportationAllowance: 1500.0,
			FeedingAllowance:        3000.0,
			TaxDeductions:           "10%",
			BankAccount:             employeeBankAccounts[0],
		},
		{
			Name:                    "David Kim",
			AnnualSalary:            55000.0,
			TransportationAllowance: 1200.0,
			FeedingAllowance:        2400.0,
			TaxDeductions:           "10%",
			BankAccount:             employeeBankAccounts[3],
		},
		{
			Name:                    "Aisha Patel",
			AnnualSalary:            78000.0,
			TransportationAllowance: 1800.0,
			FeedingAllowance:        2800.0,
			TaxDeductions:           "10%",
			BankAccount:             employeeBankAccounts[6],
		},
		{
			Name:                    "Hassan Ahmed",
			AnnualSalary:            70000.0,
			TransportationAllowance: 1400.0,
			FeedingAllowance:        2600.0,
			TaxDeductions:           "10%",
			BankAccount:             employeeBankAccounts[9]},
	}

	fmt.Println("✓ Employees registered")

	// put all the employee into a company variable
	companyPayroll := []payroll.Payment{}
	withdrawalPayroll := []withdrawals.Payment{}

	// Range through each employee and append to the payroll
	for _, employee := range hybrid {
		companyPayroll = append(companyPayroll, employee)
		withdrawalPayroll = append(withdrawalPayroll, employee)
	}
	for _, employee := range remote {
		companyPayroll = append(companyPayroll, employee)
		withdrawalPayroll = append(withdrawalPayroll, employee)
	}
	for _, employee := range fulltime {
		companyPayroll = append(companyPayroll, employee)
		withdrawalPayroll = append(withdrawalPayroll, employee)
	}

	payroll.ProcessPayroll(companyPayroll, 5)

	fmt.Println("Salary Payment successfull")

	withdrawByAccount := map[string]float64{
		"ACC-001": 3000,
		"ACC-002": 1200,
		"ACC-003": 800,
		"ACC-004": 450,
		"ACC-005": 300,
		"ACC-006": 15000,
	}

	withdrawals.ProcessWithdrawal(withdrawalPayroll, withdrawByAccount, 2)

	fmt.Println("\n\n╔════════════════════════════════════════╗")
	fmt.Println("║      FINAL ACCOUNT BALANCE           ║")
	fmt.Println("╚════════════════════════════════════════╝")

	// Display all accounts with their final balances
	for _, account := range employeeBankAccounts {
		fmt.Println(account)
	}

	fmt.Println("\n✓ System demonstration complete!")
}
