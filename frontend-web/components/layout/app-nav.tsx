"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, Compass, Bell, FileText, User } from "lucide-react";

import { Logo } from "@/components/layout/logo";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/today", label: "Today", icon: Home },
  { href: "/explore", label: "Explore", icon: Compass },
  { href: "/updates", label: "Updates", icon: Bell },
  { href: "/report", label: "Report", icon: FileText },
  { href: "/profile", label: "Profile", icon: User },
];

export function AppNav() {
  const pathname = usePathname();

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 flex items-center justify-around border-t border-background-stone bg-background py-3 md:static md:h-screen md:w-64 md:flex-col md:items-stretch md:justify-start md:border-r md:border-t-0 md:border-forest md:bg-forest-deep md:p-4">
      {/* Logo - Only visible on desktop sidebar */}
      <Logo />

      {navItems.map((item) => {
        const isActive = pathname === item.href;

        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex w-full flex-col items-center gap-1 rounded-lg px-3 py-2 transition-colors md:w-auto md:flex-row md:gap-3",
              isActive
                ? "bg-forest/10 text-forest md:bg-white/10 md:text-white"
                : "text-text-muted hover:bg-background-mist md:text-sage-soft md:hover:bg-white/10 md:hover:text-white"
            )}
          >
            <item.icon className="h-5 w-5" />
            <span className="text-metadata font-medium md:text-body">
              {item.label}
            </span>
          </Link>
        );
      })}
    </nav>
  );
}