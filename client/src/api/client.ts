import type {
  Account,
  Employee,
  PayrollRun,
  PayrollSummary,
  WithdrawalRun,
  WithdrawalSummary,
} from "../types/types";

const BASE_URL = "https://go-payroll-processor-system.onrender.com";

async function fetchJSON<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${url}`, options);
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}

export const api = {
  // Employees
  getEmployees: () => fetchJSON<Employee[]>("/api/employees"),

  getEmployeeCount: () => fetchJSON<{ count: number }>("/api/employees/count"),

  getEmployeeById: (id: number) => fetchJSON<Employee>(`/api/employees/${id}`),

  // Accounts
  getAccounts: () => fetchJSON<Account[]>("/api/accounts"),

  // Payroll
  runPayroll: () =>
    fetchJSON<PayrollSummary>("/api/payroll/run", { method: "POST" }),

  getPayrollHistory: () => fetchJSON<PayrollRun[]>("/api/payroll/history"),

  // Withdrawals
  runWithdrawals: () =>
    fetchJSON<WithdrawalSummary>("/api/withdrawals/run", { method: "POST" }),

  getWithdrawalHistory: () =>
    fetchJSON<WithdrawalRun[]>("/api/withdrawals/history"),
};
