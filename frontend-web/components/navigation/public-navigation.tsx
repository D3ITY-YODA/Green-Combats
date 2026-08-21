// components/navigation/public-navigation.tsx

"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Sun, Compass, Bell, PlusCircle, User } from "lucide-react";
import type { LucideIcon } from "lucide-react";

// --- Configuration ---

interface NavItem {
  label: string;
  href: string;
  icon: LucideIcon;
}

const publicNavigation: NavItem[] = [
  { label: "Today", href: "/today", icon: Sun },
  { label: "Explore", href: "/explore", icon: Compass },
  { label: "Updates", href: "/updates", icon: Bell },
  { label: "Report", href: "/report", icon: PlusCircle },
  { label: "Profile", href: "/profile", icon: User },
];

// --- Component ---

export function PublicNavigation() {
  const pathname = usePathname();

  // Helper to determine if a link is currently active
  const isActive = (href: string) => {
    // Exact match for root-level tabs, or prefix match for nested routes (e.g., /report/new)
    if (href === "/today") return pathname === "/today";
    return pathname.startsWith(href);
  };

  return (
    <nav 
      className="fixed bottom-0 left-0 right-0 z-40 border-t border-stone bg-white pb-safe pt-2"
      aria-label="Main public navigation"
    >
      <ul className="mx-auto flex max-w-7xl items-center justify-around px-2 pb-2 sm:pb-4">
        {publicNavigation.map((item) => {
          const active = isActive(item.href);
          const Icon = item.icon;

          return (
            <li key={item.href} className="flex-1">
              <Link
                href={item.href}
                className={`
                  flex flex-col items-center gap-1 rounded-xl px-2 py-2 text-xs font-medium transition-colors
                  focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2
                  ${active 
                    ? "text-forest" 
                    : "text-muted hover:text-charcoal"
                  }
                `}
                aria-current={active ? "page" : undefined}
              >
                <Icon 
                  className={`h-6 w-6 ${active ? "stroke-[2.5px]" : "stroke-2"}`} 
                  aria-hidden="true" 
                />
                <span>{item.label}</span>
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
