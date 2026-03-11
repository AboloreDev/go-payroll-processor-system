package main

import (
	"abah-go/projects/payment/helpers"
	"net/http"
)

// GetAll handles GET /api/accounts
// returns all account balances as JSON
func (app *Application) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accounts, err := app.employeeRepository.GetAll()
	if err != nil {
		app.errorLog.Println("could not fetch accounts:", err)
		http.Error(w, "could not fetch accounts", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, accounts)
}
