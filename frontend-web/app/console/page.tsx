// app/console/page.tsx

import Link from "next/link";
import { fetchServer } from "@/lib/api/server-client";
import type { ConsoleOverview, ConsolePriorityItem } from "@/types/console";

// Fetch data on the server using the secure server client
async function getConsoleOverview(): Promise<ConsoleOverview> {
  try {
    return await fetchServer<ConsoleOverview>("/v1/console/overview");
  } catch {
    // Backend may be unavailable; render with safe (empty) defaults rather
    // than fabricated data.
    return {
      organization: { name: "Your Organization" },
      important_updates: 0,
      community_reports: 0,
      delayed_sources: 0,
      pending_reviews: 0,
      priority_items: [],
    };
  }
}

export default async function ConsoleOverviewPage() {
  const overview = await getConsoleOverview();

  return (
    <div className="mx-auto max-w-7xl space-y-8">
      {/* Page Header */}
      <header>
        <h1 className="text-2xl font-semibold text-charcoal">Overview</h1>
        <p className="mt-1 text-lg font-medium text-muted">
          {overview.organization.name}
        </p>
      </header>

      {/* Key Metrics Grid */}
      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <MetricCard label="Important updates" value={overview.important_updates} variant="important" />
        <MetricCard label="Community updates" value={overview.community_reports} variant="neutral" />
        <MetricCard label="Information delayed" value={overview.delayed_sources} variant="warning" />
        <MetricCard label="Pending review" value={overview.pending_reviews} variant="action" />
      </section>

      {/* Priority Actions */}
      <section className="space-y-4">
        <h2 className="text-xl font-semibold text-charcoal">Priority actions</h2>

        {overview.priority_items.length === 0 ? (
          <div className="rounded-2xl border border-stone bg-white p-8 text-center text-muted">
            All caught up. There are no priority actions at this time.
          </div>
        ) : (
          <ul className="space-y-3">
            {overview.priority_items.map((item) => (
              <PriorityActionItem key={item.id} item={item} />
            ))}
          </ul>
        )}
      </section>

      {/* Quick Actions / CTAs */}
      <section className="flex flex-wrap gap-4 pt-4">
        <Link
          href="/console/local-conditions"
          className="inline-flex items-center justify-center rounded-xl bg-forest px-5 py-3 text-sm font-medium text-white transition hover:bg-forest-dark focus:outline-none focus:ring-4 focus:ring-forest/20"
        >
          View local conditions
        </Link>

        <Link
          href="/console/updates/new"
          className="inline-flex items-center justify-center rounded-xl border border-stone bg-white px-5 py-3 text-sm font-medium text-charcoal transition hover:bg-stone/50 focus:outline-none focus:ring-4 focus:ring-stone/20"
        >
          Create update
        </Link>
      </section>
    </div>
  );
}

// --- Internal Components for cleaner page structure ---

interface MetricCardProps {
  label: string;
  value: number;
  variant: "important" | "neutral" | "warning" | "action";
}

function MetricCard({ label, value, variant }: MetricCardProps) {
  const valueColor = {
    important: "text-important",
    neutral: "text-charcoal",
    warning: "text-warning",
    action: "text-forest",
  }[variant];

  return (
    <div className="rounded-2xl border border-stone bg-white p-6">
      <p className="text-sm font-medium text-muted">{label}</p>
      <p className={`mt-2 text-4xl font-bold ${valueColor}`}>
        {value}
      </p>
    </div>
  );
}

interface PriorityActionItemProps {
  item: ConsolePriorityItem;
}

function PriorityActionItem({ item }: PriorityActionItemProps) {
  const priorityIndicator = {
    normal: "bg-stone",
    important: "bg-important",
    urgent: "bg-urgent",
  }[item.priority];

  return (
    <li>
      <Link
        href={item.href}
        className="flex items-start gap-4 rounded-2xl border border-stone bg-white p-5 transition hover:border-sage hover:bg-soft-sage"
      >
        <span
          className={`mt-1.5 h-2.5 w-2.5 flex-shrink-0 rounded-full ${priorityIndicator}`}
          aria-hidden="true"
        />
        <div className="flex-1">
          <h3 className="text-base font-semibold text-charcoal">{item.title}</h3>
          <p className="mt-1 text-sm text-muted">{item.message}</p>
        </div>
        <span className="text-sm font-medium text-forest">View &rarr;</span>
      </Link>
    </li>
  );
}
