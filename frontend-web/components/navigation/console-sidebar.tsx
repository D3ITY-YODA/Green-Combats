// components/navigation/console-sidebar.tsx

"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { 
  LayoutDashboard, 
  Map, 
  Megaphone, 
  BookOpen, 
  MessageSquare, 
  Database, 
  FolderKanban, 
  BarChart3, 
  Shield, 
  Activity,
  X,
  HelpCircle
} from "lucide-react";
import type { LucideIcon } from "lucide-react";

// --- Types & Configuration ---

interface NavItem {
  label: string;
  href: string;
  icon: LucideIcon;
  // In a real app, you would check permissions here
  // requiredPermission?: Permission; 
}

const consoleNavigation: NavItem[] = [
  { label: "Overview", href: "/console", icon: LayoutDashboard },
  { label: "Local conditions", href: "/console/local-conditions", icon: Map },
  { label: "Updates", href: "/console/updates", icon: Megaphone },
  { label: "Guidance", href: "/console/guidance", icon: BookOpen },
  { label: "Community reports", href: "/console/reports", icon: MessageSquare },
  { label: "Information", href: "/console/information", icon: Database },
  { label: "Projects", href: "/console/projects", icon: FolderKanban },
  { label: "Impact", href: "/console/impact", icon: BarChart3 },
  { label: "Administration", href: "/console/people", icon: Shield },
  { label: "Activity", href: "/console/activity", icon: Activity },
];

interface ConsoleSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

// --- Component ---

export function ConsoleSidebar({ isOpen, onClose }: ConsoleSidebarProps) {
  const pathname = usePathname();

  // Helper to determine if a link is currently active
  const isActive = (href: string) => {
    if (href === "/console") {
      return pathname === "/console";
    }
    return pathname.startsWith(href);
  };

  return (
    <>
      {/* Mobile Overlay Backdrop */}
      {isOpen && (
        <div 
          className="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm lg:hidden"
          onClick={onClose}
          aria-hidden="true"
        />
      )}

      {/* Sidebar Container */}
      <aside
        className={`
          fixed inset-y-0 left-0 z-50 w-72 transform border-r border-stone bg-white 
          transition-transform duration-200 ease-in-out
          ${isOpen ? "translate-x-0" : "-translate-x-full"} 
          lg:translate-x-0
        `}
        aria-label="Console navigation"
      >
        <div className="flex h-full flex-col">
          {/* Sidebar Header / Logo Area */}
          <div className="flex h-16 items-center justify-between border-b border-stone px-6">
            <Link href="/console" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-forest/10">
                <LayoutDashboard className="h-4 w-4 text-forest" aria-hidden="true" />
              </div>
              <span className="text-lg font-bold tracking-tight text-charcoal">
                Green Compass
              </span>
            </Link>
            
            {/* Close button for mobile */}
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg p-2 text-muted hover:bg-stone/50 hover:text-charcoal lg:hidden focus:outline-none focus:ring-2 focus:ring-forest"
              aria-label="Close navigation menu"
            >
              <X className="h-5 w-5" aria-hidden="true" />
            </button>
          </div>

          {/* Navigation Links */}
          <nav className="flex-1 overflow-y-auto px-4 py-6">
            <ul className="space-y-1">
              {consoleNavigation.map((item) => {
                const active = isActive(item.href);
                const Icon = item.icon;
                
                return (
                  <li key={item.href}>
                    <Link
                      href={item.href}
                      onClick={onClose} // Close sidebar on mobile when a link is clicked
                      className={`
                        flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors
                        focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2
                        ${active 
                          ? "bg-forest/10 text-forest" 
                          : "text-muted hover:bg-stone/50 hover:text-charcoal"
                        }
                      `}
                      aria-current={active ? "page" : undefined}
                    >
                      <Icon className="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                      {item.label}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </nav>

          {/* Sidebar Footer */}
          <div className="border-t border-stone p-4">
            <Link
              href="/help"
              className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted transition-colors hover:bg-stone/50 hover:text-charcoal focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2"
            >
              <HelpCircle className="h-5 w-5 flex-shrink-0" aria-hidden="true" />
              Help & Support
            </Link>
            <p className="mt-4 px-3 text-xs text-muted">
              Green Compass v1.0.0
            </p>
          </div>
        </div>
      </aside>
    </>
  );
}
