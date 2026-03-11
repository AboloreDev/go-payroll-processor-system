package main

import (
	"net/http"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// Employee routes
	mux.HandleFunc("/api/employees", app.GetAllEmployees)
	mux.HandleFunc("/api/employees/count", app.GetCount)
	mux.HandleFunc("/api/employees/", app.GetById)

	// Payroll routes
	mux.HandleFunc("/api/payroll/run", app.RunPayroll)
	mux.HandleFunc("/api/payroll/history", app.GetHistory)

	// Account routes
	mux.HandleFunc("/api/accounts", app.GetAllAccounts)

	// Withdrawal routes
	mux.HandleFunc("/api/withdrawals/run", app.RunWithdrawal)
	mux.HandleFunc("/api/withdrawals/history", app.GetWithdrawalHistory)

	// Wrap with CORS so React can talk to us
	return CORS(mux)
}
