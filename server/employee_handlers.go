package main

import (
	"abah-go/projects/payment/helpers"
	"net/http"
	"strconv"
	"strings"
)

// GetAll handles GET /api/employees
// returns all 250 employees as JSON
func (app *Application) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	employees, err := app.employeeRepository.GetAll()
	if err != nil {
		app.errorLog.Println("could not fetch employees:", err)
		http.Error(w, "could not fetch employees", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, employees)
}

// GetByID handles GET /api/employees/{id}
// returns one employee as JSON
func (app *Application) GetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the id from the URL
	// URL looks like /api/employees/42
	// We split by "/" and grab the last part
	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[len(parts)-1])

	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}
	employee, err := app.employeeRepository.GetById(id)
	if err != nil {
		app.errorLog.Println("could not fetch employee:", err)
		http.Error(w, "could not fetch employee", http.StatusInternalServerError)
		return
	}

	// Employee not found
	if employee == nil {
		app.errorLog.Println("could not find employee:", err)
		http.Error(w, "employee not found", http.StatusNotFound)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, employee)
}

// GetCount handles GET /api/employees/count
// returns total number of employees
func (app *Application) GetCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count, err := app.employeeRepository.GetCount()
	if err != nil {
		app.errorLog.Println("could not count employees:", err)
		http.Error(w, "could not count employee", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}
