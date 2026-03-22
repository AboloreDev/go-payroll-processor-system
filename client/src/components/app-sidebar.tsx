"use client";

import * as React from "react";

import {
  CommandIcon,
  Users,
  CreditCard,
  Banknote,
  ArrowDownCircle,
  History,
  LayoutDashboardIcon,
} from "lucide-react";
import { NavMain } from "./nav-main";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar";

const navItems = [
  { title: "Dashboard", url: "/", icon: LayoutDashboardIcon },
  { title: "Employees", url: "/employees", icon: Users },
  { title: "Accounts", url: "/accounts", icon: CreditCard },
  { title: "Payroll", url: "/payroll", icon: Banknote },
  { title: "Withdrawals", url: "/withdrawals", icon: ArrowDownCircle },
  { title: "History", url: "/history", icon: History },
];

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader className="px-2 py-4 bg-black text-white">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton className="data-[slot=sidebar-menu-button]:p-1.5! cursor-pointer hover:bg-sidebar-accent rounded-lg transition-colors">
              <CommandIcon className="size-5! shrink-0" />
              <div className="px-2 py-1">
                <p className="text-xs font-semibold text-muted-foreground uppercase tracking-widest">
                  Crixus PLC
                </p>
                <h1 className="text-sm font-bold">Payroll System</h1>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent className="px-2 bg-black text-white">
        <NavMain items={navItems} />
      </SidebarContent>
    </Sidebar>
  );
}
