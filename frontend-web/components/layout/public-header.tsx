// components/layout/public-header.tsx

"use client";

import Link from "next/link";
import { MapPin, Bell, User, ChevronDown } from "lucide-react";

// TODO: Replace with actual hooks from your state management (Zustand/Context)
// import { usePlace } from "@/hooks/use-place";
// import { useAuth } from "@/hooks/use-auth";

export function PublicHeader() {
  // Mock data for demonstration. In production, pull this from your hooks.
  const currentPlaceName = "Lower Valley"; 
  const hasUnreadUpdates = true;

  return (
    <header className="sticky top-0 z-40 w-full border-b border-stone bg-white">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        
        {/* Left Side: Current Place Selector */}
        <Link
          href="/places"
          className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-charcoal transition hover:bg-stone/50 focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2"
          aria-label={`Change current place. Currently viewing: ${currentPlaceName}`}
        >
          <MapPin className="h-4 w-4 text-forest" aria-hidden="true" />
          <span className="truncate max-w-[150px] sm:max-w-none">{currentPlaceName}</span>
          <ChevronDown className="h-4 w-4 text-muted" aria-hidden="true" />
        </Link>

        {/* Right Side: Notifications & Profile */}
        <div className="flex items-center gap-1 sm:gap-2">
          {/* Notifications / Updates */}
          <Link
            href="/updates"
            className="relative rounded-lg p-2 text-muted transition hover:bg-stone/50 hover:text-charcoal focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2"
            aria-label="View updates and notifications"
          >
            <Bell className="h-5 w-5" aria-hidden="true" />
            
            {/* Unread Indicator */}
            {hasUnreadUpdates && (
              <span 
                className="absolute right-1.5 top-1.5 flex h-2.5 w-2.5"
                aria-hidden="true"
              >
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-important opacity-75"></span>
                <span className="relative inline-flex h-2.5 w-2.5 rounded-full bg-important"></span>
              </span>
            )}
          </Link>

          {/* User Profile */}
          <Link
            href="/profile"
            className="rounded-lg p-2 text-muted transition hover:bg-stone/50 hover:text-charcoal focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2"
            aria-label="View profile and settings"
          >
            <User className="h-5 w-5" aria-hidden="true" />
          </Link>
        </div>
      </div>
    </header>
  );
}
