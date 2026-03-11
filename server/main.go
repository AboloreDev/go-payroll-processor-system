package main

import (
	bankaccount "abah-go/projects/payment/bank-account"
	"abah-go/projects/payment/employees"
	"abah-go/projects/payment/payroll"
	"abah-go/projects/payment/withdrawals"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	errorLog             *log.Logger
	infoLog              *log.Logger
	seedRepository       SeedRepository
	employeeRepository   employees.EmployeeRepository
	accountRepository    bankaccount.AccountRepository
	withdrawalRepository withdrawals.WithdrawalRepository
	payrollRepository    payroll.PayrollRepository
}

func init() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  CRIXUS PLC PAYROLL & BANKING SYSTEM   ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)
	dbName := "./database.db"
	db, err := connectDb(dbName)
	if err != nil {
		errorLog.Fatal("failed to connect to database:", err)
	}
	defer db.Close()
	infoLog.Println("Database Connected successfully!")

	app := &Application{
		errorLog:             errorLog,
		infoLog:              infoLog,
		seedRepository:       NewSeedRepository(db),
		employeeRepository:   employees.NewSQLEmployeeRepository(db),
		payrollRepository:    payroll.NewSQLPayrollRepository(db),
		accountRepository:    bankaccount.NewSQLAccountRepository(db),
		withdrawalRepository: withdrawals.NewSQLWithdrawalRepository(db),
	}

	// Lets run the scripts
	err = app.seedRepository.Setup()
	if err != nil {
		app.errorLog.Println("setup failed:", err)
	}

	// only seed if DB is empty
	count, err := app.employeeRepository.GetCount()
	if err != nil {
		app.errorLog.Fatal("could not check employee count:", err)
	}

	if count == 0 {
		app.infoLog.Println("empty database — seeding 250 employees...")
		if err := app.SeedingIntoDB(); err != nil {
			app.errorLog.Fatal("seeding failed:", err)
		}
	} else {
		app.infoLog.Printf("database already has %d employees — skipping seed", count)
	}

	app.infoLog.Println("server starting on port 4040")
	if err := app.Serve(); err != nil {
		app.errorLog.Fatal("opening port failed:", err)
	}

}

func connectDb(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func (app *Application) SeedingIntoDB() error {
	// Cear old data
	err := app.seedRepository.Clear()
	if err != nil {
		app.errorLog.Println("clear failed:", err)
		return err
	}

	app.infoLog.Println("seeding 250 employees...")

	// Seed employees into the database
	for i := 1; i <= 250; i++ {
		employee := RandomEmployeeGenerator(i)

		employeeID, err := app.seedRepository.InsertEmployee(employee)
		if err != nil {
			app.errorLog.Printf("failed to insert employee %s: %v", employee.Name, err)
			return err
		}

		accountNumber := fmt.Sprintf("ACC-%04d", i)
		err = app.seedRepository.InsertBankAccount(employeeID, accountNumber, employee.Name)
		if err != nil {
			app.errorLog.Printf("failed to insert bank account for %s: %v", employee.Name, err)
			return err
		}
	}

	app.infoLog.Println("seed complete — 250 employees ready")
	return nil
}
