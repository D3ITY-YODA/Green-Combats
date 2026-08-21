// components/today/today-client-widget.tsx

"use client";

import { usePlace } from "@/hooks/use-place";
import { useToday } from "@/hooks/use-today";
import { DelayedState } from "@/components/states/delayed-state";
import { UpdateCard } from "@/components/updates/update-card";

export function TodayClientWidget() {
  const { currentPlace } = usePlace();
  const { data, isLoading, isError, refetch } = useToday({ 
    placeId: currentPlace?.id 
  });

  if (isLoading) return <div className="animate-pulse h-40 bg-stone/20 rounded-2xl" />;
  if (isError) return <button onClick={refetch}>Failed to load. Try again.</button>;
  if (!data) return null;

  // Handle the delayed state defined in the blueprint
  if (data.status.data_status === "delayed" || data.status.data_status === "stale") {
    return <DelayedState onRetry={refetch} />;
  }

  return (
    <div className="space-y-4">
      {data.updates.map((update) => (
        <UpdateCard key={update.id} update={update} />
      ))}
    </div>
  );
}
