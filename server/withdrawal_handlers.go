package main

import (
	"abah-go/projects/payment/helpers"
	"abah-go/projects/payment/withdrawals"
	"net/http"
)

func (app *Application) RunWithdrawal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	employeesList, err := app.employeeRepository.GetAll()
	if err != nil {
		app.errorLog.Println("could not fetch employees:", err)
		http.Error(w, "could not reset balances", http.StatusInternalServerError)
		return
	}

	withdrawals, err := withdrawals.ProcessWithdrawal(
		employeesList,
		app.accountRepository,
		app.withdrawalRepository,
		10,
	)
	if err != nil {
		app.errorLog.Println("withdrawal processing failed:", err)
		http.Error(w, "withdrawal processing failed", http.StatusInternalServerError)
		return
	}

	app.infoLog.Printf("withdrawal run complete — $%.2f withdrawn by %d employees",
		withdrawals.TotalWithdrawn, withdrawals.SuccessCount)

	helpers.WriteJSON(w, http.StatusOK, withdrawals)
}

func (app *Application) GetWithdrawalHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history, err := app.withdrawalRepository.GetHistory()
	if err != nil {
		app.errorLog.Println("could not fetch withdrawal history:", err)
		http.Error(w, "could not fetch history", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, history)
}
