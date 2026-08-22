// hooks/use-today.tsx

"use client";

import { useQuery, UseQueryOptions, UseQueryResult } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api/client";
import type { TodayData } from "@/lib/api/today-api";
import type { UserPlace } from "@/types/places";

interface UseTodayOptions {
  placeId?: string;
  enabled?: boolean;
}

export function useToday(
  options: UseTodayOptions = {}
): UseQueryResult<TodayData, Error> {
  const { placeId, enabled = true } = options;

  return useQuery<TodayData, Error>({
    queryKey: ["today", placeId],
    queryFn: async () => {
      const url = placeId ? `/v1/context?place_id=${placeId}` : "/v1/context";
      return fetchAPI<TodayData>(url);
    },
    enabled: enabled && !!placeId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

export function useTodayWithPlace(
  place: UserPlace | null,
  queryOptions?: Omit<UseQueryOptions<TodayData, Error>, "queryKey" | "queryFn" | "enabled">
): UseQueryResult<TodayData, Error> {
  return useToday({
    placeId: place?.id,
    enabled: !!place,
    ...queryOptions,
  });
}