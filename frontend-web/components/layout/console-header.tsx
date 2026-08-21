// components/layout/console-header.tsx

"use client";

import { useState } from "react";
import Link from "next/link";
import { 
  Menu, 
  Bell, 
  Search, 
  ChevronDown, 
  LogOut, 
  Settings, 
  User 
} from "lucide-react";

// In a real app, you would import these from your auth context or fetch them via server props
// For this example, we assume a mock user and organization context.
const mockUser = {
  name: "Amina Yusuf",
  role: "Community Observer",
  avatarUrl: null, // Fallback to initials
};

const mockOrganization = {
  name: "Lower Valley Water Authority",
};

export function ConsoleHeader() {
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false); // Toggle for sidebar on mobile

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-stone bg-white px-4 md:px-8">
      {/* Left Side: Mobile Menu & Organization Context */}
      <div className="flex items-center gap-4">
        {/* Mobile Sidebar Toggle */}
        <button
          type="button"
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          className="rounded-lg p-2 text-muted hover:bg-stone/50 hover:text-charcoal lg:hidden focus:outline-none focus:ring-2 focus:ring-forest"
          aria-label="Toggle navigation menu"
          aria-expanded={isMobileMenuOpen}
        >
          <Menu className="h-5 w-5" aria-hidden="true" />
        </button>

        {/* Organization Name / Breadcrumb */}
        <div className="flex flex-col">
          <span className="text-xs font-medium uppercase tracking-wider text-muted">
            Organization
          </span>
          <h2 className="text-sm font-semibold text-charcoal md:text-base">
            {mockOrganization.name}
          </h2>
        </div>
      </div>

      {/* Right Side: Actions & User Menu */}
      <div className="flex items-center gap-2 md:gap-4">
        {/* Search (Hidden on small mobile) */}
        <button
          type="button"
          className="hidden rounded-lg p-2 text-muted hover:bg-stone/50 hover:text-charcoal sm:block focus:outline-none focus:ring-2 focus:ring-forest"
          aria-label="Search console"
        >
          <Search className="h-5 w-5" aria-hidden="true" />
        </button>

        {/* Notifications */}
        <button
          type="button"
          className="relative rounded-lg p-2 text-muted hover:bg-stone/50 hover:text-charcoal focus:outline-none focus:ring-2 focus:ring-forest"
          aria-label="View notifications"
        >
          <Bell className="h-5 w-5" aria-hidden="true" />
          {/* Notification Badge */}
          <span className="absolute right-1.5 top-1.5 flex h-2 w-2">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-important opacity-75"></span>
            <span className="relative inline-flex h-2 w-2 rounded-full bg-important"></span>
          </span>
        </button>

        {/* User Profile Dropdown */}
        <div className="relative">
          <button
            type="button"
            onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
            className="flex items-center gap-2 rounded-lg border border-stone bg-white px-3 py-1.5 text-sm font-medium text-charcoal transition hover:bg-stone/50 focus:outline-none focus:ring-2 focus:ring-forest"
            aria-expanded={isUserMenuOpen}
            aria-haspopup="true"
          >
            {/* Avatar / Initials */}
            <div className="flex h-6 w-6 items-center justify-center rounded-full bg-forest/10 text-xs font-bold text-forest">
              {mockUser.name.charAt(0)}
            </div>
            <span className="hidden sm:inline">{mockUser.name}</span>
            <ChevronDown className="h-4 w-4 text-muted" aria-hidden="true" />
          </button>

          {/* Dropdown Menu */}
          {isUserMenuOpen && (
            <div className="absolute right-0 mt-2 w-56 overflow-hidden rounded-xl border border-stone bg-white shadow-lg ring-1 ring-black/5">
              <div className="border-b border-stone px-4 py-3">
                <p className="text-sm font-semibold text-charcoal">{mockUser.name}</p>
                <p className="text-xs text-muted">{mockUser.role}</p>
              </div>
              
              <nav className="py-1">
                <Link
                  href="/profile"
                  className="flex items-center gap-3 px-4 py-2 text-sm text-charcoal hover:bg-stone/50 focus:bg-stone/50 focus:outline-none"
                  onClick={() => setIsUserMenuOpen(false)}
                >
                  <User className="h-4 w-4 text-muted" aria-hidden="true" />
                  My Profile
                </Link>
                <Link
                  href="/profile/settings"
                  className="flex items-center gap-3 px-4 py-2 text-sm text-charcoal hover:bg-stone/50 focus:bg-stone/50 focus:outline-none"
                  onClick={() => setIsUserMenuOpen(false)}
                >
                  <Settings className="h-4 w-4 text-muted" aria-hidden="true" />
                  Settings
                </Link>
              </nav>

              <div className="border-t border-stone py-1">
                <button
                  type="button"
                  className="flex w-full items-center gap-3 px-4 py-2 text-sm text-important hover:bg-stone/50 focus:bg-stone/50 focus:outline-none"
                  // onClick={handleSignOut}
                >
                  <LogOut className="h-4 w-4" aria-hidden="true" />
                  Sign out
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Mobile Sidebar Overlay (Simplified for this component) */}
      {isMobileMenuOpen && (
        <div 
          className="fixed inset-0 z-20 bg-black/20 lg:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
          aria-hidden="true"
        />
      )}
    </header>
  );
}
