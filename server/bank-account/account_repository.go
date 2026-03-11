package bankaccount

import (
	"database/sql"
	"fmt"
	"time"
)

type Account struct {
	ID            int `json:"id"`
	EmployeeId    string `json:"employee_id"`
	AccountNumber string `json:"account_number"`
	AccountOwner  string `json:"account_owner"`
	Balance       float64 `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

// Interface
type AccountRepository interface {
	GetAllAccounts() ([]Account, error)
	GetAccountsById(employeeId int) (*Account, error)
	UpdateBalance(employeeId int, amount float64) error
	ResetAllBalance() error
}

type SQLAccountRepository struct {
	db *sql.DB
}

func NewSQLAccountRepository(db *sql.DB) *SQLAccountRepository {
	return &SQLAccountRepository{db: db}
}

func (r *SQLAccountRepository) GetAllAccounts() ([]Account, error) {
	// Prepare db statement
	queryStatement := `SELECT
			id,
			employee_id,
			account_number,
			account_owner,
			balance,
			created_at
		FROM bank_accounts
		ORDER BY id ASC 
		`
	rows, err := r.db.Query(queryStatement)
	if err != nil {
		return nil, fmt.Errorf("could not query bank accounts: %w", err)
	}

	var account []Account
	for rows.Next() {
		var a Account
		err := rows.Scan(
			&a.ID,
			&a.EmployeeId,
			&a.AccountNumber,
			&a.AccountOwner,
			&a.Balance,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan account row: %w", err)
		}

		account = append(account, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account rows: %w", err)
	}

	return account, nil
}

func (r *SQLAccountRepository) GetAccountsById(employeeId int) (*Account, error) {
	// Prepare statement
	queryStatement := `SELECT id, employee_id, account_number, account_owner, balance FROM bank_accounts
	WHERE id = ?`

	rows := r.db.QueryRow(queryStatement, employeeId)
	var a Account 
	err := rows.Scan(&a.ID, &a.EmployeeId, &a.AccountNumber, &a.AccountOwner, &a.Balance)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to query row %w", err)
	}
	 return &a, nil
}


// UpdateBalance adds amount to an employee's account balance
// This is called after each successful salary deposit
// or withdrawal during payroll processing
func (r *SQLAccountRepository) UpdateBalance(employeeId int, amount float64) error {
	// Prepare statement
	query := `
		UPDATE bank_accounts
		SET balance = balance + ?
		WHERE employee_id = ?
	`
	rows, err := r.db.Exec(query, amount, employeeId)
	if err != nil {
		return fmt.Errorf("could not update balance for employee %d: %w", employeeId, err)
	}

	// Make sure the update actually affected a row
	// If rowsAffected is 0 it means no account was found
	// for that employee id
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no account found for employee %d", employeeId)
	}

	return nil
}

func (r *SQLAccountRepository) ResetAllBalance() error {
	_, err := r.db.Exec("UPDATE bank_accounts SET balance = 0")
	if err != nil {
		return fmt.Errorf("could not reset balances: %w", err)
	}
	return nil
}