import { useQuery } from "@tanstack/react-query";
import {
  UsersIcon,
  BanknoteIcon,
  CheckCircleIcon,
  ArrowDownCircleIcon,
} from "lucide-react";
import { api } from "../api/client";
import {
  Card,
  CardAction,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Skeleton } from "./ui/skeleton";
import { Badge } from "./ui/badge";

export function SectionCards() {
  const { data: countData, isLoading: countLoading } = useQuery({
    queryKey: ["employeeCount"],
    queryFn: api.getEmployeeCount,
  });

  const { data: payrollHistory, isLoading: payrollLoading } = useQuery({
    queryKey: ["payrollHistory"],
    queryFn: api.getPayrollHistory,
  });

  const { data: withdrawalHistory, isLoading: withdrawalLoading } = useQuery({
    queryKey: ["withdrawalHistory"],
    queryFn: api.getWithdrawalHistory,
  });

  // Get the most recent run from history
  const lastPayroll = payrollHistory?.[0];
  const lastWithdrawal = withdrawalHistory?.[0];

  const totalPaid = lastPayroll?.total_paid ?? 0;
  const successCount = lastPayroll?.success_count ?? 0;
  const totalWithdrawn = lastWithdrawal?.total_withdrawn ?? 0;
  const employeeCount = countData?.count ?? 0;

  const formatCurrency = (amount: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: 0,
    }).format(amount);

  const isLoading = countLoading || payrollLoading || withdrawalLoading;

  return (
    <div className="grid bg-black text-white grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-linear-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card *:data-[slot=card]:shadow-xs lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4 dark:*:data-[slot=card]:bg-card">
      {/* Total Employees */}
      <Card className="@container/card ">
        <CardHeader>
          <CardDescription>Total Employees</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            {isLoading ? <Skeleton className="h-8 w-24" /> : employeeCount}
          </CardTitle>
          <CardAction>
            <Badge variant="outline">
              <UsersIcon className="size-3" />
              Active
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Fulltime, Remote & Hybrid
            <UsersIcon className="size-4" />
          </div>
          <div className="text-muted-foreground">Across all employee types</div>
        </CardFooter>
      </Card>

      {/* Last Payroll Total */}
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Last Payroll Total</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            {isLoading ? (
              <Skeleton className="h-8 w-32" />
            ) : lastPayroll ? (
              formatCurrency(totalPaid)
            ) : (
              "No runs yet"
            )}
          </CardTitle>
          <CardAction>
            <Badge variant="outline">
              <BanknoteIcon className="size-3" />
              Paid
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Most recent payroll run
            <BanknoteIcon className="size-4" />
          </div>
          <div className="text-muted-foreground">
            {lastPayroll
              ? `Run on ${new Date(lastPayroll.ran_at).toLocaleDateString()}`
              : "Run payroll to see results"}
          </div>
        </CardFooter>
      </Card>

      {/* Successful Payments */}
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Successful Payments</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            {isLoading ? (
              <Skeleton className="h-8 w-20" />
            ) : lastPayroll ? (
              `${successCount} / ${employeeCount}`
            ) : (
              "—"
            )}
          </CardTitle>
          <CardAction>
            <Badge variant="outline">
              <CheckCircleIcon className="size-3" />
              {lastPayroll
                ? `${Math.round((successCount / employeeCount) * 100)}%`
                : "0%"}
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Last run success rate
            <CheckCircleIcon className="size-4" />
          </div>
          <div className="text-muted-foreground">
            {lastPayroll
              ? `${lastPayroll.fail_count} failed payments`
              : "No payroll data yet"}
          </div>
        </CardFooter>
      </Card>

      {/* Total Withdrawn */}
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Last Withdrawal Total</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            {isLoading ? (
              <Skeleton className="h-8 w-32" />
            ) : lastWithdrawal ? (
              formatCurrency(totalWithdrawn)
            ) : (
              "No runs yet"
            )}
          </CardTitle>
          <CardAction>
            <Badge variant="outline">
              <ArrowDownCircleIcon className="size-3" />
              Withdrawn
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Most recent withdrawal run
            <ArrowDownCircleIcon className="size-4" />
          </div>
          <div className="text-muted-foreground">
            {lastWithdrawal
              ? `${lastWithdrawal.success_count} employees withdrew`
              : "Run withdrawals to see results"}
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
