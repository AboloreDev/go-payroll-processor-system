import { useQuery } from "@tanstack/react-query";
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
import { Input } from "../components/ui/input";
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../components/ui/select";
import { UsersIcon } from "lucide-react";

const ITEMS_PER_PAGE = 10;

type EmployeeTypeFilter = "all" | "fulltime" | "remote" | "hybrid";

export default function Employees() {
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState<EmployeeTypeFilter>("all");
  const [currentPage, setCurrentPage] = useState(1);

  const {
    data: employees,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["employees"],
    queryFn: api.getEmployees,
  });

  const formatCurrency = (amount: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: 0,
    }).format(amount);

  // Filter by search and type
  const filtered = employees?.filter((e) => {
    const matchesSearch =
      e.name.toLowerCase().includes(search.toLowerCase()) ||
      e.account_number.toLowerCase().includes(search.toLowerCase());
    const matchesType = typeFilter === "all" || e.employee_type === typeFilter;
    return matchesSearch && matchesType;
  });

  const handleSearch = (value: string) => {
    setSearch(value);
    setCurrentPage(1);
  };

  const handleTypeFilter = (value: EmployeeTypeFilter) => {
    setTypeFilter(value);
    setCurrentPage(1);
  };

  // Counts per type for summary cards
  const fulltimeCount =
    employees?.filter((e) => e.employee_type === "fulltime").length ?? 0;
  const remoteCount =
    employees?.filter((e) => e.employee_type === "remote").length ?? 0;
  const hybridCount =
    employees?.filter((e) => e.employee_type === "hybrid").length ?? 0;

  // Pagination
  const totalItems = filtered?.length ?? 0;
  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const paginated = filtered?.slice(startIndex, endIndex);

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

  const typeBadgeVariant = (type: string) => {
    switch (type) {
      case "fulltime":
        return "default";
      case "remote":
        return "secondary";
      case "hybrid":
        return "outline";
      default:
        return "outline";
    }
  };

  const monthlySalary = (
    e: Awaited<ReturnType<typeof api.getEmployees>>[number],
  ) => {
    switch (e.employee_type) {
      case "fulltime":
        return (
          ((e.annual_salary + e.transport_allowance + e.feeding_allowance) /
            12) *
          (1 - e.tax_rate)
        );
      case "remote":
        return e.hourly_rate * e.hours_worked * (1 - e.tax_rate);
      case "hybrid":
        return (e.annual_salary / 12) * (1 - e.tax_rate);
      default:
        return 0;
    }
  };

  return (
    <div className="flex flex-col gap-6 px-4 lg:px-6">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-bold">Employees</h1>
        <p className="text-muted-foreground text-sm">
          All registered employees across Crixus PLC
        </p>
      </div>

      {/* Summary cards */}
      <div className="grid grid-cols-1 gap-4 @xl/main:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Full-time</CardDescription>
            <CardTitle className="text-2xl flex items-center gap-2">
              {isLoading ? <Skeleton className="h-8 w-16" /> : fulltimeCount}
              <Badge variant="default" className="text-xs">
                Fulltime
              </Badge>
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Remote</CardDescription>
            <CardTitle className="text-2xl flex items-center gap-2">
              {isLoading ? <Skeleton className="h-8 w-16" /> : remoteCount}
              <Badge variant="secondary" className="text-xs">
                Remote
              </Badge>
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Hybrid</CardDescription>
            <CardTitle className="text-2xl flex items-center gap-2">
              {isLoading ? <Skeleton className="h-8 w-16" /> : hybridCount}
              <Badge variant="outline" className="text-xs">
                Hybrid
              </Badge>
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Table */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between gap-4 flex-wrap">
            <div>
              <CardTitle className="flex items-center gap-2">
                <UsersIcon className="size-4" />
                All Employees
              </CardTitle>
              <CardDescription>
                {totalItems > 0
                  ? `Showing ${startIndex + 1}–${Math.min(endIndex, totalItems)} of ${totalItems} employees`
                  : "Search or filter employees"}
              </CardDescription>
            </div>

            <div className="flex items-center gap-2">
              <Select
                value={typeFilter}
                onValueChange={(v) => handleTypeFilter(v as EmployeeTypeFilter)}
              >
                <SelectTrigger className="w-36 cursor-pointer">
                  <SelectValue placeholder="Filter by type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  <SelectItem value="fulltime">Full-time</SelectItem>
                  <SelectItem value="remote">Remote</SelectItem>
                  <SelectItem value="hybrid">Hybrid</SelectItem>
                </SelectContent>
              </Select>

              <Input
                placeholder="Search employees..."
                value={search}
                onChange={(e) => handleSearch(e.target.value)}
                className="max-w-xs"
              />
            </div>
          </div>
        </CardHeader>

        <CardContent className="flex flex-col gap-4">
          {isError && (
            <p className="text-destructive text-sm">
              Failed to load employees. Is the Go server running?
            </p>
          )}

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Account</TableHead>
                <TableHead>Monthly Salary</TableHead>
                <TableHead>Balance</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              {isLoading ? (
                Array.from({ length: ITEMS_PER_PAGE }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell>
                      <Skeleton className="h-4 w-36" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-20" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-24" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-24" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-24" />
                    </TableCell>
                  </TableRow>
                ))
              ) : paginated?.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={5}
                    className="text-center text-muted-foreground py-8"
                  >
                    No employees found
                  </TableCell>
                </TableRow>
              ) : (
                paginated?.map((employee) => (
                  <TableRow key={employee.id}>
                    <TableCell className="font-medium">
                      {employee.name}
                    </TableCell>
                    <TableCell>
                      <Badge variant={typeBadgeVariant(employee.employee_type)}>
                        {employee.employee_type}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      {employee.account_number}
                    </TableCell>
                    <TableCell className="font-mono">
                      {formatCurrency(monthlySalary(employee))}
                    </TableCell>
                    <TableCell className="font-mono">
                      <Badge
                        variant={employee.balance > 0 ? "default" : "secondary"}
                      >
                        {formatCurrency(employee.balance)}
                      </Badge>
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
