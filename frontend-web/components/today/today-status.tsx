// components/today/today-status.tsx

import type { DataStatus } from "@/types/common"; // Adjust path to match your types directory
import { formatDateTime } from "@/lib/formatters"; // Adjust path to your date formatting utility

interface TodayStatusProps {
  title: string;
  message: string;
  updatedAt?: string;
  dataStatus: DataStatus;
}

export function TodayStatus({
  title,
  message,
  updatedAt,
  dataStatus,
}: TodayStatusProps) {
  // Handle delayed or stale data states
  if (dataStatus === "delayed" || dataStatus === "stale") {
    return (
      <section 
        className="rounded-2xl border border-stone bg-white p-6"
        role="status"
        aria-live="polite"
      >
        <p className="text-sm font-medium text-muted">
          Information delayed
        </p>
        <h1 className="mt-2 text-2xl font-semibold text-charcoal">
          The latest update for this place is not available yet.
        </h1>
        {updatedAt && (
          <p className="mt-3 text-sm text-muted">
            Last reliable update: {formatDateTime(updatedAt)}
          </p>
        )}
      </section>
    );
  }

  // Normal / Current state
  return (
    <section className="rounded-2xl border border-stone bg-white p-6">
      {updatedAt && (
        <p className="text-sm text-muted">
          {formatDateTime(updatedAt)}
        </p>
      )}
      <h1 className="mt-2 text-2xl font-semibold text-charcoal">
        {title}
      </h1>
      <p className="mt-2 text-base leading-relaxed text-muted">
        {message}
      </p>
    </section>
  );
}
