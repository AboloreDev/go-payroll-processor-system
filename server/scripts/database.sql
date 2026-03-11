-- CREATE A TABLE FOR THE EMPLOYEE TYPE
CREATE TABLE IF NOT EXISTS employee_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL
)

-- CREATE A TABLE FOR ALL EMPLOYEES BELONGING O THE COMPANY
CREATE TABLE IF NOT EXISTS employees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    employee_type_id INTEGER NOT NULL FOREIGN KEY REFERENCES employee_types(id),
    annual_salary REAL DEFAULT 0,
    transport_allowance REAL DEFAULT 0,
    feeding_allowance   REAL DEFAULT 0,
    hourly_rate         REAL DEFAULT 0,
    hours_worked        REAL DEFAULT 0,
    tax_rate            REAL NOT NULL DEFAULT 0.10,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP
)

-- CREATE A BANK ACCOUNT TABLE
CREATE TABLE IF NOT EXISTS bank_accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    account_number TEXT NOT NULL,
    account_owner TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    balance REAL NOT NULL DEFAULT 0.0
)

-- SUMMARY OF EACH PAYROLL RUN
-- one row is created every time /api/payroll/run is called
CREATE TABLE IF NOT EXISTS payroll_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    total_paid      REAL NOT NULL DEFAULT 0.0,
    success_count   INTEGER NOT NULL DEFAULT 0,
    fail_count      INTEGER NOT NULL DEFAULT 0,
    ran_at          DATETIME DEFAULT CURRENT_TIMESTAMP
)

-- One row per employee per payroll run
-- tells you exactly who got paid what in each run
CREATE TABLE IF NOT EXISTS payroll_logs (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id          INTEGER NOT NULL REFERENCES payroll_runs(id),
    employee_id     INTEGER NOT NULL REFERENCES employees(id),
    account_number  TEXT NOT NULL,
    amount_paid     REAL NOT NULL DEFAULT 0.0,
    status          TEXT NOT NULL,
    error_message   TEXT,
    duration_ms     INTEGER,
    paid_at         DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- Summary of each withdrawal run
CREATE TABLE IF NOT EXISTS withdrawal_runs (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    total_withdrawn     REAL NOT NULL DEFAULT 0.0,
    success_count       INTEGER NOT NULL DEFAULT 0,
    fail_count          INTEGER NOT NULL DEFAULT 0,
    ran_at              DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- One row per employee per withdrawal run
CREATE TABLE IF NOT EXISTS withdrawal_logs (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id              INTEGER NOT NULL REFERENCES withdrawal_runs(id),
    employee_id         INTEGER NOT NULL REFERENCES employees(id),
    account_number      TEXT NOT NULL,
    amount_withdrawn    REAL NOT NULL DEFAULT 0.0,
    balance_after       REAL NOT NULL DEFAULT 0.0,
    status              TEXT NOT NULL,
    error_message       TEXT,
    duration_ms         INTEGER,
    withdrawn_at        DATETIME DEFAULT CURRENT_TIMESTAMP
);
