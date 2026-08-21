// hooks/use-today.ts

"use client";

import { useQuery } from "@tanstack/react-query";
import { getTodayData, type TodayData } from "@/lib/api/today-api";

interface UseTodayOptions {
  placeId?: string;
}

interface UseTodayResult {
  data: TodayData | null;
  isLoading: boolean;
  isError: boolean;
  refetch: () => void;
}

/**
 * React Query hook that fetches today's context from the Go backend.
 */
export function useToday({ placeId }: UseTodayOptions = {}): UseTodayResult {
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ["today", placeId],
    queryFn: () => getTodayData(placeId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 1,
  });

  return {
    data: data ?? null,
    isLoading,
    isError,
    refetch,
  };
}

// --- TodayClientWidget component (moved here to avoid circular import) ---

import { usePlace } from "@/hooks/place-context";
import { DelayedState } from "@/components/states/delayed-state";
import { UpdateCard } from "@/components/updates/update-card";

export function TodayClientWidget() {
  const { currentPlace } = usePlace();
  const { data, isLoading, isError, refetch } = useToday({
    placeId: currentPlace?.id,
  });

  if (isLoading)
    return <div className="animate-pulse h-40 bg-stone/20 rounded-2xl" />;
  if (isError)
    return (
      <button onClick={refetch}>Failed to load. Try again.</button>
    );
  if (!data) return null;

  if (data.hasImportantUpdates) {
    return <DelayedState onRetry={refetch} />;
  }

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-stone bg-white p-5">
        <p className="text-sm font-medium text-muted">{data.location}</p>
        <p className="mt-2 text-base text-charcoal">{data.localOutlook}</p>
        <p className="mt-1 text-sm text-muted">
          Water: {data.waterOutlook}
        </p>
        <p className="mt-1 text-xs text-muted">
          Last updated: {data.lastUpdated}
        </p>
      </div>
    </div>
  );
}
