import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { Skeleton } from "../components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../components/ui/table";
import { Badge } from "../components/ui/badge";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationEllipsis,
} from "../components/ui/pagination";
import {
  BanknoteIcon,
  ArrowDownCircleIcon,
  CheckCircleIcon,
  XCircleIcon,
  HistoryIcon,
} from "lucide-react";
import { useState } from "react";

const ITEMS_PER_PAGE = 10;

// Combined event type for the unified timeline
type RunType = "payroll" | "withdrawal";

interface TimelineEvent {
  id: string; // unique key: "payroll-1" or "withdrawal-1"
  runId: number;
  type: RunType;
  totalAmount: number;
  successCount: number;
  failCount: number;
  ranAt: string;
}

export default function History() {
  const [currentPage, setCurrentPage] = useState(1);

  const { data: payrollHistory, isLoading: payrollLoading } = useQuery({
    queryKey: ["payrollHistory"],
    queryFn: api.getPayrollHistory,
  });

  const { data: withdrawalHistory, isLoading: withdrawalLoading } = useQuery({
    queryKey: ["withdrawalHistory"],
    queryFn: api.getWithdrawalHistory,
  });

  const isLoading = payrollLoading || withdrawalLoading;

  const formatCurrency = (amount: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: 0,
    }).format(amount);

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    if (!dateStr || date.getFullYear() < 2000) return "—";
    return date.toLocaleString();
  };

  // Merge payroll and withdrawal runs into one sorted timeline
  const timeline: TimelineEvent[] = [
    ...(payrollHistory?.map((run) => ({
      id: `payroll-${run.id}`,
      runId: run.id,
      type: "payroll" as RunType,
      totalAmount: run.total_paid,
      successCount: run.success_count,
      failCount: run.fail_count,
      ranAt: run.ran_at,
    })) ?? []),
    ...(withdrawalHistory?.map((run) => ({
      id: `withdrawal-${run.id}`,
      runId: run.id,
      type: "withdrawal" as RunType,
      totalAmount: run.total_withdrawn,
      successCount: run.success_count,
      failCount: run.fail_count,
      ranAt: run.ran_at,
    })) ?? []),
  ].sort((a, b) => {
    // Sort by date descending — most recent first
    const dateA = new Date(a.ranAt).getTime();
    const dateB = new Date(b.ranAt).getTime();
    if (dateA < 0 || dateB < 0) return 0;
    return dateB - dateA;
  });

  // Summary stats
  const totalPayrollRuns = payrollHistory?.length ?? 0;
  const totalWithdrawalRuns = withdrawalHistory?.length ?? 0;
  const totalPayrollPaid =
    payrollHistory?.reduce((sum, r) => sum + r.total_paid, 0) ?? 0;
  const totalWithdrawn =
    withdrawalHistory?.reduce((sum, r) => sum + r.total_withdrawn, 0) ?? 0;

  // Pagination
  const totalItems = timeline.length;
  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const paginated = timeline.slice(startIndex, endIndex);

  const getPageNumbers = () => {
    const pages: (number | "ellipsis")[] = [];
    if (totalPages <= 5) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
    } else {
      pages.push(1);
      if (currentPage > 3) pages.push("ellipsis");
      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);
      for (let i = start; i <= end; i++) pages.push(i);
      if (currentPage < totalPages - 2) pages.push("ellipsis");
      pages.push(totalPages);
    }
    return pages;
  };

  return (
    <div className="flex flex-col gap-6 px-4 lg:px-6">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-bold">History</h1>
        <p className="text-muted-foreground text-sm">
          Combined timeline of all payroll and withdrawal runs
        </p>
      </div>

      {/* Summary cards */}
      <div className="grid grid-cols-1 gap-4 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Payroll Runs</CardDescription>
            <CardTitle className="text-2xl">
              {isLoading ? <Skeleton className="h-8 w-16" /> : totalPayrollRuns}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Paid Out</CardDescription>
            <CardTitle className="text-2xl">
              {isLoading ? (
                <Skeleton className="h-8 w-32" />
              ) : (
                formatCurrency(totalPayrollPaid)
              )}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Withdrawal Runs</CardDescription>
            <CardTitle className="text-2xl">
              {isLoading ? (
                <Skeleton className="h-8 w-16" />
              ) : (
                totalWithdrawalRuns
              )}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Withdrawn</CardDescription>
            <CardTitle className="text-2xl">
              {isLoading ? (
                <Skeleton className="h-8 w-32" />
              ) : (
                formatCurrency(totalWithdrawn)
              )}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Timeline table */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <HistoryIcon className="size-4" />
            All Runs
          </CardTitle>
          <CardDescription>
            {totalItems > 0
              ? `Showing ${startIndex + 1}–${Math.min(endIndex, totalItems)} of ${totalItems} runs`
              : "No runs yet — go to Payroll or Withdrawals to run"}
          </CardDescription>
        </CardHeader>

        <CardContent className="flex flex-col gap-4">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Type</TableHead>
                <TableHead>Run ID</TableHead>
                <TableHead>Total Amount</TableHead>
                <TableHead>Successful</TableHead>
                <TableHead>Failed</TableHead>
                <TableHead>Date</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              {isLoading ? (
                Array.from({ length: ITEMS_PER_PAGE }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell>
                      <Skeleton className="h-4 w-24" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-8" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-28" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-16" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-16" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-32" />
                    </TableCell>
                  </TableRow>
                ))
              ) : paginated.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="text-center text-muted-foreground py-8"
                  >
                    No runs yet — go to Payroll or Withdrawals to get started
                  </TableCell>
                </TableRow>
              ) : (
                paginated.map((event) => (
                  <TableRow key={event.id}>
                    <TableCell>
                      {event.type === "payroll" ? (
                        <Badge variant="default" className="gap-1.5">
                          <BanknoteIcon className="size-3" />
                          Payroll
                        </Badge>
                      ) : (
                        <Badge variant="secondary" className="gap-1.5">
                          <ArrowDownCircleIcon className="size-3" />
                          Withdrawal
                        </Badge>
                      )}
                    </TableCell>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      #{event.runId}
                    </TableCell>
                    <TableCell className="font-mono font-medium">
                      {formatCurrency(event.totalAmount)}
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="gap-1">
                        <CheckCircleIcon className="size-3 text-green-500" />
                        {event.successCount}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          event.failCount > 0 ? "destructive" : "outline"
                        }
                        className="gap-1"
                      >
                        {event.failCount > 0 && (
                          <XCircleIcon className="size-3" />
                        )}
                        {event.failCount}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {formatDate(event.ranAt)}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>

          {totalPages > 1 && (
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    size="default"
                    onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                    className={
                      currentPage === 1
                        ? "pointer-events-none opacity-50"
                        : "cursor-pointer"
                    }
                  />
                </PaginationItem>

                {getPageNumbers().map((page, index) =>
                  page === "ellipsis" ? (
                    <PaginationItem key={`ellipsis-${index}`}>
                      <PaginationEllipsis />
                    </PaginationItem>
                  ) : (
                    <PaginationItem key={page}>
                      <PaginationLink
                        size="default"
                        isActive={currentPage === page}
                        onClick={() => setCurrentPage(page)}
                        className="cursor-pointer"
                      >
                        {page}
                      </PaginationLink>
                    </PaginationItem>
                  ),
                )}

                <PaginationItem>
                  <PaginationNext
                    size="default"
                    onClick={() =>
                      setCurrentPage((p) => Math.min(totalPages, p + 1))
                    }
                    className={
                      currentPage === totalPages
                        ? "pointer-events-none opacity-50"
                        : "cursor-pointer"
                    }
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
