// components/layout/app-nav.tsx
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, Compass, Bell, FileText, User } from "lucide-react";
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
    <nav className="fixed bottom-0 left-0 right-0 z-50 border-t border-background-stone bg-background md:static md:border-t-0 md:border-r md:border-forest md:bg-forest-deep md:w-64 md:h-screen md:flex-col md:p-4 flex items-center justify-around py-3">
      {navItems.map((item) => {
        const isActive = pathname === item.href;
        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex flex-col md:flex-row items-center gap-1 md:gap-3 px-3 py-2 rounded-lg transition-colors w-full md:w-auto",
              isActive 
                // Active State: Green tint on mobile, White tint on desktop
                ? "bg-forest/10 text-forest md:bg-white/10 md:text-white" 
                // Inactive State: Muted text on mobile, Soft Sage text on desktop
                : "text-text-muted hover:bg-background-mist md:text-sage-soft md:hover:bg-white/10 md:hover:text-white"
            )}
          >
            <item.icon className="h-5 w-5" />
            <span className="text-metadata md:text-body font-medium">{item.label}</span>
          </Link>
        );
      })}
    </nav>
  );
}
