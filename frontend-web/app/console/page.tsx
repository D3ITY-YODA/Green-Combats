"use client";

import { useState, useEffect } from "react";
import { Bell, FileText, Database, AlertCircle } from "lucide-react";
import { getDashboardClient } from "@/lib/api/reporting";
import type { DashboardData } from "@/lib/api/reporting";

export default function ConsoleOverview() {
  const [data, setData] = useState<DashboardData>({
    important_updates: 2,
    community_reports: 14,
    information_delayed: 1,
    pending_review: 4,
  });

  useEffect(() => {
    getDashboardClient().then(setData).catch(() => {});
  }, []);

  const stats = [
    { label: "Important updates", value: String(data.important_updates ?? 0), icon: Bell, color: "text-status-important" },
    { label: "Community reports", value: String(data.community_reports ?? 0), icon: FileText, color: "text-forest" },
    { label: "Information delayed", value: String(data.information_delayed ?? 0), icon: AlertCircle, color: "text-status-watch" },
    { label: "Pending review", value: String(data.pending_review ?? 0), icon: Database, color: "text-sky" },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-text-charcoal">Overview</h1>
        <p className="text-text-muted mt-1">Welcome back. Here is what needs your attention.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <div key={stat.label} className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
            <div className={`p-3 rounded-lg bg-background-mist ${stat.color}`}>
              <stat.icon className="h-5 w-5" />
            </div>
            <div>
              <p className="text-2xl font-bold text-text-charcoal">{stat.value}</p>
              <p className="text-xs text-text-muted">{stat.label}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="bg-background rounded-xl border border-background-stone p-6">
        <h2 className="text-lg font-semibold text-text-charcoal mb-4">Priority actions</h2>
        <div className="space-y-4">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-background-mist border border-background-stone">
            <FileText className="h-5 w-5 text-forest mt-0.5" />
            <div className="flex-1">
              <h3 className="text-sm font-semibold text-text-charcoal">Review community updates</h3>
              <p className="text-xs text-text-muted mt-1">{data.pending_review ?? 4} updates are waiting for review.</p>
            </div>
            <button className="text-xs font-medium text-forest hover:underline">Review</button>
          </div>
          <div className="flex items-start gap-3 p-4 rounded-lg bg-background-mist border border-background-stone">
            <AlertCircle className="h-5 w-5 text-status-watch mt-0.5" />
            <div className="flex-1">
              <h3 className="text-sm font-semibold text-text-charcoal">Information delayed</h3>
              <p className="text-xs text-text-muted mt-1">The latest water update is delayed.</p>
            </div>
            <button className="text-xs font-medium text-forest hover:underline">Check sources</button>
          </div>
        </div>
      </div>
    </div>
  );
}
