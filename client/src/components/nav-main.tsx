import { Link, useLocation } from "react-router-dom";
import { type LucideIcon } from "lucide-react";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar";

interface NavItem {
  title: string;
  url: string;
  icon: LucideIcon;
}

export function NavMain({ items }: { items: NavItem[] }) {
  const location = useLocation();

  return (
    <SidebarGroup className="py-2">
      <SidebarGroupLabel className="px-2 mb-1">Navigation</SidebarGroupLabel>
      <SidebarMenu className="space-y-1">
        {items.map((item) => (
          <SidebarMenuItem key={item.title}>
            <Link to={item.url} className="w-full">
              <SidebarMenuButton
                isActive={location.pathname === item.url}
                className="w-full cursor-pointer px-3 py-2.5 rounded-lg transition-colors gap-3 hover:bg-sidebar-accent"
              >
                <item.icon className="size-4 shrink-0" />
                <span className="text-sm font-medium">{item.title}</span>
              </SidebarMenuButton>
            </Link>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}
