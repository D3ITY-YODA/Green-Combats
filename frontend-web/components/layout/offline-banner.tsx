// components/layout/offline-banner.tsx
"use client";

import { useOnlineStatus } from "@/hooks/use-online-status";
import { WifiOff } from "lucide-react";

export function OfflineBanner() {
  const isOnline = useOnlineStatus();

  if (isOnline) return null;

  return (
    <div className="fixed top-0 left-0 right-0 z-[100] bg-status-delayed text-white px-4 py-2 flex items-center justify-center gap-2 text-metadata font-medium shadow-sm">
      <WifiOff className="h-4 w-4" />
      <span>You are currently offline. Showing cached information.</span>
    </div>
    // Note: We use 'status-delayed' (grey) to keep it calm, not alarming like 'emergency' (red)
  );
}
