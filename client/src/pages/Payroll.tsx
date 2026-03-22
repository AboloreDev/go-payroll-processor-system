import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
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
import { Button } from "../components/ui/button";
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
  PlayIcon,
  CheckCircleIcon,
  XCircleIcon,
} from "lucide-react";
import type { PayrollResult } from "../types/types";

const ITEMS_PER_PAGE = 10;

export default function Payroll() {
  const queryClient = useQueryClient();
  const [currentPage, setCurrentPage] = useState(1);
  const [visibleRows, setVisibleRows] = useState<PayrollResult[]>([]);

  // Fetch payroll history
  const { data: history, isLoading: historyLoading } = useQuery({
    queryKey: ["payrollHistory"],
    queryFn: api.getPayrollHistory,
  });

  // Run payroll mutation
  const { mutate: runPayroll, isPending } = useMutation({
    mutationFn: api.runPayroll,
    onSuccess: (summary) => {
      // Reset visible rows
      setVisibleRows([]);
      setCurrentPage(1);

      // Animate rows one by one every 20ms
      summary.results.forEach((result, index) => {
        setTimeout(() => {
          setVisibleRows((prev) => [...prev, result]);
        }, index * 20);
      });

      // Refetch history and accounts after run
      queryClient.invalidateQueries({ queryKey: ["payrollHistory"] });
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["employeeCount"] });
    },
  });

  const formatCurrency = (amount: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: 2,
    }).format(amount);

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    if (!dateStr || date.getFullYear() < 2000) return "—";
    return date.toLocaleString();
  };

  // Pagination for live results
  const totalItems = visibleRows.length;
  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const paginated = visibleRows.slice(startIndex, endIndex);

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

  const lastRun = history?.[0];

  return (
    <div className="flex flex-col gap-6 px-4 lg:px-6">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-bold">Payroll</h1>
        <p className="text-muted-foreground text-sm">
          Run and monitor payroll for all 250 employees
        </p>
      </div>

      {/* Summary cards */}
      <div className="grid grid-cols-1 gap-4 @xl/main:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Last Run Total</CardDescription>
            <CardTitle className="text-2xl">
              {historyLoading ? (
                <Skeleton className="h-8 w-32" />
              ) : lastRun ? (
                formatCurrency(lastRun.total_paid)
              ) : (
                "No runs yet"
              )}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Last Run Date</CardDescription>
            <CardTitle className="text-2xl">
              {historyLoading ? (
                <Skeleton className="h-8 w-32" />
              ) : lastRun ? (
                formatDate(lastRun.ran_at)
              ) : (
                "—"
              )}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Runs</CardDescription>
            <CardTitle className="text-2xl">
              {historyLoading ? (
                <Skeleton className="h-8 w-16" />
              ) : (
                (history?.length ?? 0)
              )}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Run Payroll button */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BanknoteIcon className="size-5" />
            Run Payroll
          </CardTitle>
          <CardDescription>
            Resets all balances to zero then processes all 250 employees
            concurrently. Results appear below in real time.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button
            onClick={() => runPayroll()}
            disabled={isPending}
            size="lg"
            className="cursor-pointer gap-2"
          >
            {isPending ? (
              <>
                <span className="size-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Processing 250 employees...
              </>
            ) : (
              <>
                <PlayIcon className="size-4" />
                Run Payroll
              </>
            )}
          </Button>
        </CardContent>
      </Card>

      {/* Live results table — only shows after a run */}
      {visibleRows.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Live Results</CardTitle>
            <CardDescription>
              {visibleRows.length} / 250 employees processed
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Employee</TableHead>
                  <TableHead>Account</TableHead>
                  <TableHead>Amount Paid</TableHead>
                  <TableHead>Duration</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {paginated.map((result) => (
                  <TableRow
                    key={result.employee_id}
                    className="animate-in fade-in slide-in-from-bottom-1 duration-200"
                  >
                    <TableCell className="font-medium">
                      {result.employee_name}
                    </TableCell>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      {result.account_number}
                    </TableCell>
                    <TableCell className="font-mono">
                      {formatCurrency(result.amount_paid)}
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {result.duration_ms}ms
                    </TableCell>
                    <TableCell>
                      {result.status === "success" ? (
                        <Badge variant="default" className="gap-1">
                          <CheckCircleIcon className="size-3" />
                          Success
                        </Badge>
                      ) : (
                        <Badge variant="destructive" className="gap-1">
                          <XCircleIcon className="size-3" />
                          Failed
                        </Badge>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
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
      )}

      {/* Payroll history */}
      <Card>
        <CardHeader>
          <CardTitle>Payroll History</CardTitle>
          <CardDescription>All past payroll runs</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Run ID</TableHead>
                <TableHead>Total Paid</TableHead>
                <TableHead>Successful</TableHead>
                <TableHead>Failed</TableHead>
                <TableHead>Date</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {historyLoading ? (
                Array.from({ length: 3 }).map((_, i) => (
                  <TableRow key={i}>
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
              ) : history?.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={5}
                    className="text-center text-muted-foreground py-8"
                  >
                    No payroll runs yet — click Run Payroll to start
                  </TableCell>
                </TableRow>
              ) : (
                history?.map((run) => (
                  <TableRow key={run.id}>
                    <TableCell className="font-mono text-sm">
                      #{run.id}
                    </TableCell>
                    <TableCell className="font-mono font-medium">
                      {formatCurrency(run.total_paid)}
                    </TableCell>
                    <TableCell>
                      <Badge variant="default" className="gap-1">
                        <CheckCircleIcon className="size-3" />
                        {run.success_count}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          run.fail_count > 0 ? "destructive" : "secondary"
                        }
                        className="gap-1"
                      >
                        {run.fail_count > 0 && (
                          <XCircleIcon className="size-3" />
                        )}
                        {run.fail_count}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {formatDate(run.ran_at)}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
