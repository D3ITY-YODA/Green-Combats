"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, Map, Bell, BookOpen, FileText, Database, FolderKanban, Users, Activity } from "lucide-react";
import { cn } from "@/lib/utils";

const navItems = [
  { label: "Overview", href: "/console", icon: LayoutDashboard },
  { label: "Local conditions", href: "/console/local-conditions", icon: Map },
  { label: "Updates", href: "/console/updates", icon: Bell },
  { label: "Guidance", href: "/console/guidance", icon: BookOpen },
  { label: "Community reports", href: "/console/reports", icon: FileText },
  { label: "Information", href: "/console/information", icon: Database },
  { label: "Projects", href: "/console/projects", icon: FolderKanban },
  { label: "People", href: "/console/people", icon: Users },
  { label: "Activity", href: "/console/activity", icon: Activity },
];

export function ConsoleSidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden lg:flex w-64 flex-col bg-forest-deep text-white h-screen sticky top-0">
      <div className="p-6 border-b border-white/10">
        <h1 className="text-lg font-bold">Green Compass</h1>
        <p className="text-xs text-sage-soft/80 mt-1">Institutional Console</p>
      </div>
      <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
                isActive ? "bg-white/10 text-white" : "text-sage-soft hover:bg-white/5 hover:text-white"
              )}
            >
              <item.icon className="h-4 w-4" />
              {item.label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
