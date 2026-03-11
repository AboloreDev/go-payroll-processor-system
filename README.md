# Crixus PLC — Payroll & Banking System

A concurrent payroll engine built with Go that processes 250 employees simultaneously using goroutines, channels, and semaphores. Includes a React dashboard to visualize payroll runs in real time.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go, net/http, SQLite (WAL mode) |
| Frontend | React, TypeScript, Vite, shadcn/ui, TanStack Query |
| Pattern | Repository pattern with interfaces |

---

## Key Features

- **Concurrent payroll engine** — 250 goroutines run simultaneously, limited to 10 at a time via semaphore channel
- **Three salary models** — fulltime, remote, and hybrid each with different calculation logic
- **Atomic DB writes** — SQLite WAL mode handles concurrent writes without a mutex
- **Live results** — frontend animates each employee result as goroutines complete
- **Withdrawals** — employees withdraw 10–50% of balance concurrently after payroll

---

## How Payroll Works
```
POST /api/payroll/run
  → Reset all 250 balances to zero
  → Spawn 250 goroutines
  → Semaphore limits to 10 concurrent
  → Each goroutine calculates salary + credits account via atomic SQL UPDATE
  → Results collected via channel → returned as JSON
  → Frontend animates rows in real time
```

---

## Salary Calculation
```
Fulltime → (annualSalary + transport + feeding) / 12 × (1 - taxRate)
Remote   → hourlyRate × hoursWorked × (1 - taxRate)
Hybrid   → annualSalary / 12 × (1 - taxRate)
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/employees` | All 250 employees |
| GET | `/api/accounts` | All bank balances |
| POST | `/api/payroll/run` | Run payroll |
| GET | `/api/payroll/history` | Past payroll runs |
| POST | `/api/withdrawals/run` | Run withdrawals |
| GET | `/api/withdrawals/history` | Past withdrawal runs |

---

## Getting Started
```bash
# Backend
cd projects/payment
go run ./cmd
# Runs on http://localhost:4040
# Seeds 250 employees automatically on first run

# Frontend
cd client
pnpm install && pnpm dev
# Runs on http://localhost:5173
```

> Requires GCC for go-sqlite3 CGO compilation

---

Built by [Your Name](https://github.com/yourusername)
