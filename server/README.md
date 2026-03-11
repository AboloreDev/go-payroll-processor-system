# go-payroll-processor-system
A Golang Salary Payment Processor
# 💼 Crixus PLC Payroll System

A real-world Go payroll system demonstrating interfaces, pointers, structs, and banking integration. Process salaries for 10 employees across different employment types.

## 🚀 Features

- **Multi-employee Support** - Manage 10+ employees with different employment types
- **Smart Payroll Processing** - Automatic salary calculation with tax deductions
- **Banking Integration** - Direct deposit to employee bank accounts
- **Type Safety** - Interface-based design for extensibility
- **Transaction Management** - Full deposit/withdrawal functionality with error handling

## 📋 Employee Types

| Type | Payment Model | Features |
|------|--------------|----------|
| **Full-time** | Annual salary | Transportation & feeding allowances |
| **Remote** | Hourly rate | Flexible hours tracking |
| **Hybrid** | Annual salary | Work-from-anywhere benefits |

## 🏗️ Architecture

```
├── BankAccount          # Banking system with deposit/withdraw
├── Payment Interface    # Contract for all employee types
│   ├── FulltimeEmployee
│   ├── RemoteEmployee
│   └── HybridEmployee
└── ProcessPayroll       # Automated payroll processing
```

## 💡 Key Concepts Demonstrated

- **Interfaces** - Polymorphic payment processing
- **Pointers** - Efficient memory management for account modifications
- **Structs** - Data modeling for employees and accounts
- **Methods** - Behavior tied to types
- **Error Handling** - Robust transaction validation
- **Composition** - Linking employees to bank accounts

## 🎯 Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/hebron-payroll.git
cd hebron-payroll

# Run the system
go run main.go
```

## 📊 Sample Output

```
╔════════════════════════════════════════╗
║    CRIXUS PLC MONTHLY PAYROLL REPORT   ║
╚════════════════════════════════════════╝

[Employee #1]
Full-Time Employee: Sarah Martinez
  Annual Salary: $85000.00
  Monthly Salary: $6712.50
✓ Deposited $6712.50 to Account ACC-001

...

TOTAL MONTHLY PAYROLL: $56,584.50
✓ Payroll processing completed.
```

## 📚 Project Structure

```
.
├── main.go                    # Main application
├── GUIDE.md                   # Comprehensive learning guide
├── EMPLOYEE_DIRECTORY.md      # Employee data reference
├── EXERCISES.md               # Practice challenges
├── EXAMPLE_OUTPUT.txt         # Expected program output
└── README.md                  # This file
```

## 🎓 Learning Resources

- **GUIDE.md** - Deep dive into Go concepts (structs, pointers, interfaces)
- **EXERCISES.md** - 10 hands-on challenges from beginner to expert
- **EMPLOYEE_DIRECTORY.md** - Complete employee data and calculations

## 💰 Payroll Breakdown

### Current Employees (10)

| Name | Type | Monthly Pay |
|------|------|-------------|
| Sarah Martinez | Full-time | $6,712.50 |
| James Chen | Remote | $7,290.00 |
| Maria Rodriguez | Hybrid | $5,400.00 |
| David Kim | Full-time | $4,395.00 |
| Jennifer Williams | Remote | $4,320.00 |
| Michael O'Brien | Hybrid | $5,100.00 |
| Aisha Patel | Full-time | $6,195.00 |
| Carlos Santos | Remote | $5,472.00 |
| Emily Thompson | Hybrid | $6,150.00 |
| Hassan Ahmed | Full-time | $5,550.00 |

**Total Monthly Payroll:** $56,584.50

## 🔧 How It Works

### 1. Create Bank Accounts
```go
employeeBankAccounts := []*BankAccount{
    {AccountNumber: "ACC-001", OwnerName: "Sarah Martinez", Balance: 0.0},
    // ... more accounts
}
```

### 2. Register Employees
```go
sarah := FulltimeEmployee{
    Name:         "Sarah Martinez",
    AnnualSalary: 85000.0,
    BankAccount:  employeeBankAccounts[0],
}
```

### 3. Process Payroll
```go
employees := []Payment{sarah, james, maria, ...}
ProcessPayroll(employees)
```

### 4. Employees Transact
```go
employeeBankAccounts[0].Withdraw(3000.0)  // Sarah pays rent
```

## 🧪 Example Usage

```go
// Add a new employee
newEmployee := RemoteEmployee{
    Name:         "Alex Johnson",
    HoursWorked:  160.0,
    HourlyRate:   50.0,
    TaxDeduction: "10%",
    BankAccount:  newAccount,
}

// Process single payroll
salary := newEmployee.CalculateMonthlyPay()
newAccount.Deposit(salary)

// Check balance
balance := newAccount.GetBalance()
fmt.Printf("Balance: $%.2f\n", balance)
```


## 🤝 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Code Highlights

### Interface-Based Design
```go
type Payment interface {
    CalculateMonthlyPay() float64
    String() string
}
```
One function processes all employee types:
```go
func ProcessPayroll(employees []Payment) {
    for _, employee := range employees {
        salary := employee.CalculateMonthlyPay()
        // Deposit to account
    }
}
```

### Pointer-Based Transactions
```go
func (account *BankAccount) Deposit(amount float64) error {
    account.Balance += amount  // Modifies actual account
    return nil
}
```

### Type-Specific Calculations
```go
// Full-time: Salary + allowances
func (e FulltimeEmployee) CalculateMonthlyPay() float64 {
    return ((e.AnnualSalary + e.FeedingAllowance + 
             e.TransportationAllowance) / 12) * 0.9
}

// Remote: Hourly rate × hours
func (e RemoteEmployee) CalculateMonthlyPay() float64 {
    return (e.HourlyRate * e.HoursWorked) * 0.9
}
```

## 🐛 Error Handling

```go
err := account.Withdraw(10000.0)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    // Handle insufficient funds, invalid amount, etc.
}
```

## 📖 Learn Go Concepts

This project teaches:

| Concept | Application |
|---------|-------------|
| **Structs** | Employee & account data modeling |
| **Interfaces** | Polymorphic employee types |
| **Pointers** | In-place account modifications |
| **Methods** | Type-specific behavior |
| **Slices** | Managing employee collections |
| **Error Handling** | Transaction validation |
| **Type Assertions** | Extracting specific employee types |


## 🙏 Acknowledgments

Built as a learning project to demonstrate Go's strengths in building type-safe, efficient business systems.

**⭐ Star this repo if you found it helpful!**

**🔗 Connect:** [Twitter](https://twitter.com/abolorreeeee) | [LinkedIn](https://linkedin.com/in/alabifathiu)
