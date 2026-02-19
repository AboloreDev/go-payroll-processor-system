package withdrawals

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"fmt"
	"sync"
	"time"
)

// Result “receipt” for each withdrawal
type PaymentResult struct {
	EmployeeName    string
	Error           error
	AmountWithdrawn float64
	Duration        time.Duration
	AccountNumber   string
	BalanceAfter    float64
}

// Interface: anything we can withdraw from must give us a bank account + String()
type Payment interface {
	GetBankAccount() *bankaccount.BankAccount
	fmt.Stringer
}

// ProcessWithdrawal runs withdrawals concurrently.
// ✅ Uses a map: accountNumber -> amountToWithdraw
func ProcessWithdrawal(employees []Payment, withdrawByAccount map[string]float64, maxConcurrent int) {
	fmt.Println("\n======================================")
	fmt.Println("           NEPTUNE BANK               ")
	fmt.Println("        CONCURRENT WITHDRAWALS        ")
	fmt.Println("======================================")

	var wg sync.WaitGroup
	resultChannel := make(chan PaymentResult)
	limiter := make(chan struct{}, maxConcurrent)

	for _, employee := range employees {
		wg.Add(1)

		go func(emp Payment) {
			defer wg.Done()

			// limit how many withdrawals can run at once
			limiter <- struct{}{}
			defer func() { <-limiter }()

			startTime := time.Now()

			acc := emp.GetBankAccount()
			if acc == nil {
				resultChannel <- PaymentResult{
					EmployeeName:  fmt.Sprintf("%v", emp),
					AccountNumber: "N/A",
					Error:         fmt.Errorf("employee has no bank account"),
					Duration:      time.Since(startTime),
				}
				return
			}

			// ✅ get withdrawal amount from map
			amount, ok := withdrawByAccount[acc.AccountNumber]
			if !ok {
				resultChannel <- PaymentResult{
					EmployeeName:    fmt.Sprintf("%v", emp),
					AccountNumber:   acc.AccountNumber,
					AmountWithdrawn: 0,
					Error:           fmt.Errorf("no withdrawal request for account %s", acc.AccountNumber),
					Duration:        time.Since(startTime),
					BalanceAfter:    acc.AccountBalance,
				}
				return
			}

			err := acc.Withdraw(amount)

			resultChannel <- PaymentResult{
				EmployeeName:    fmt.Sprintf("%v", emp),
				AccountNumber:   acc.AccountNumber,
				AmountWithdrawn: amount,
				Error:           err,
				Duration:        time.Since(startTime),
				BalanceAfter:    acc.AccountBalance,
			}
		}(employee)
	}

	// close results when all goroutines finish
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var totalWithdrawn float64
	successCount := 0
	failCount := 0

	for result := range resultChannel {
		if result.Error != nil {
			failCount++
			fmt.Printf("❌ FAILED: %s -> %s | -$%.2f | %v\n",
				result.EmployeeName, result.AccountNumber, result.AmountWithdrawn, result.Error)
			continue
		}

		successCount++
		totalWithdrawn += result.AmountWithdrawn
		fmt.Printf("✅ WITHRAWAL SUCCESSFUL:     %s -> %s | -$%.2f | Balance: $%.2f | (%s)\n",
			result.EmployeeName, result.AccountNumber, result.AmountWithdrawn, result.BalanceAfter, result.Duration)
	}

	fmt.Println("\n--------------------------------------")
	fmt.Printf("Success: %d | Failed: %d\n", successCount, failCount)
	fmt.Printf("TOTAL WITHDRAWN: $%.2f\n", totalWithdrawn)
	fmt.Println("======================================")
}
