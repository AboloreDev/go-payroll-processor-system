package bankaccount

import (
	"errors"
	"fmt"
	"sync"
)

// Bank Account struct definition
type BankAccount struct {
	AccountNumber string
	AccountOwner string
	AccountBalance float64
	Mutex sync.Mutex
}

// Methods
// Deposit
func (account *BankAccount) Deposit(amount float64) error {
	account.Mutex.Lock()
	defer account.Mutex.Unlock()
	if (amount < 0) {
		return errors.New("You can't deposit a negative number")
	}
	// Add the amount deposited into the back account balance
	account.AccountBalance = account.AccountBalance + amount
	fmt.Printf("✓ Deposited $%.2f to Account %s belonging to %s. New Balance: $%.2f\n", 
		amount, account.AccountNumber,  account.AccountOwner, account.AccountBalance)
	return nil
}

func (account *BankAccount) Withdraw(amount float64) error {
	account.Mutex.Lock()
	defer account.Mutex.Unlock()
	if (amount < 0) {
		return errors.New("You can't withdraw a negative number")
	}
	if account.AccountBalance < amount {
		return fmt.Errorf("insufficient funds: balance $%.2f, attempted withdrawal $%.2f",
			account.AccountBalance, amount)
	}
	account.AccountBalance = account.AccountBalance - amount
	fmt.Printf("✓ Withdrew $%.2f from Account %s belonging to %s. New Balance: $%.2f\n", 
		amount, account.AccountNumber,  account.AccountOwner, account.AccountBalance,)
	
	return nil
}

func (account *BankAccount) GetBalance() float64{
	account.Mutex.Lock()
	defer account.Mutex.Unlock()
	return account.AccountBalance
}


func (account *BankAccount) String() string {
	return fmt.Sprintf(
		"\n┌─────────────────────────────────────┐\n"+
		"│ Account Number: %s │\n"+
		"│ Account Owner:  %s │\n"+
		"│ Account Balance: $%.2f │\n"+
		"└─────────────────────────────────────┘",
		account.AccountNumber,
		account.AccountOwner,
		account.AccountBalance,
	)
}

