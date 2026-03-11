package main

import (
	"abah-go/projects/payment/helpers"
	"abah-go/projects/payment/payroll"
	"net/http"
)

func (app *Application) RunPayroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reset all balance to zero
	err := app.accountRepository.ResetAllBalance()
	if err != nil {
		app.errorLog.Println("could not reset balances:", err)
		http.Error(w, "could not reset balances", http.StatusInternalServerError)
		return
	}

	// Get all the employee
	employeeList, err := app.employeeRepository.GetAll()
	if err != nil {
		app.errorLog.Println("could not fetch employees:", err)
		http.Error(w, "could not fetch employees", http.StatusInternalServerError)
		return
	}

	// Run the process payroll function
	// The go routine that process payment for all employee
	summary, err := payroll.ProcessPayroll(employeeList, app.accountRepository, app.payrollRepository, 10)
	if err != nil {
		app.errorLog.Println("payroll processing failed:", err)
		http.Error(w, "payroll processing failed", http.StatusInternalServerError)
		return
	}

	app.infoLog.Printf("payroll run complete — paid $%.2f to %d employees",
		summary.TotalPaid, summary.SuccessCount)

	// Write the json
	helpers.WriteJSON(w, http.StatusOK, summary)
}

// Get history of payroll
func (app *Application) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history, err := app.payrollRepository.GetHistory()
	if err != nil {
		app.errorLog.Println("could not fetch payroll history:", err)
		http.Error(w, "could not fetch history", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, history)
}
