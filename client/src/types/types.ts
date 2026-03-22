export type EmployeeType = "fulltime" | "remote" | "hybrid";

export interface Employee {
  id: number;
  name: string;
  employee_type: EmployeeType;
  annual_salary: number;
  transport_allowance: number;
  feeding_allowance: number;
  hourly_rate: number;
  hours_worked: number;
  tax_rate: number;
  account_number: string;
  account_owner: string;
  balance: number;
  created_at: string;
}

export interface Account {
  id: number;
  employee_id: number;
  account_number: string;
  account_owner: string;
  balance: number;
  created_at: string;
}

export interface PayrollResult {
  employee_id: number;
  employee_name: string;
  account_number: string;
  amount_paid: number;
  status: "success" | "failed";
  error_message?: string;
  duration_ms: number;
}

export interface PayrollSummary {
  runId: number;
  total_paid: number;
  success_count: number;
  fail_count: number;
  ran_at: string;
  results: PayrollResult[];
}

export interface PayrollRun {
  id: number;
  total_paid: number;
  success_count: number;
  fail_count: number;
  ran_at: string;
}

export interface WithdrawalResult {
  employee_id: number;
  employee_name: string;
  account_number: string;
  withdrawn_amount: number;
  balance_after_withdrawal: number;
  status: "success" | "failed";
  error_message?: string;
  duration_ms: number;
}

export interface WithdrawalSummary {
  runId: number;
  total_withdrawn: number;
  success_count: number;
  fail_count: number;
  results: WithdrawalResult[];
  ran_at: string;
}

export interface WithdrawalRun {
  id: number;
  total_withdrawn: number;
  success_count: number;
  fail_count: number;
  ran_at: string;
}
