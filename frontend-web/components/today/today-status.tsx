import type { DataStatus } from "@/types/common";

interface TodayStatusProps {
  title: string;
  message: string;
  updatedAt?: string;
  dataStatus: DataStatus;
}

export function TodayStatus({ title, message, updatedAt, dataStatus }: TodayStatusProps) {
  if (dataStatus === "delayed" || dataStatus === "stale") {
    return (
      <section className="rounded-2xl border border-background-stone bg-background p-6">
        <p className="text-sm font-medium text-status-delayed">Information delayed</p>
        <h1 className="mt-2 text-2xl font-semibold text-text-charcoal">
          The latest update for this place is not available yet.
        </h1>
        {updatedAt && <p className="mt-3 text-sm text-text-muted">Last reliable update: {updatedAt}</p>}
      </section>
    );
  }

  return (
    <section className="rounded-2xl border border-background-stone bg-background p-6">
      <p className="text-sm text-text-muted">{updatedAt && `Updated ${updatedAt}`}</p>
      <h1 className="mt-2 text-2xl font-semibold text-text-charcoal">{title}</h1>
      <p className="mt-2 text-base text-text-muted">{message}</p>
    </section>
  );
}
