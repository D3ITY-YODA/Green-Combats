"use client";

import { useState, useEffect } from "react";
import { BookOpen, AlertTriangle, Info, Clock, ChevronRight } from "lucide-react";
import { getUpdatesClient } from "@/lib/api/updates";
import type { Update } from "@/types/updates";

const TYPE_META: Record<string, { label: string; color: string; bg: string; icon: React.ElementType }> = {
  alert: { label: "Emergency", color: "text-status-emergency", bg: "bg-status-emergency/10 border-status-emergency/20", icon: AlertTriangle },
  today: { label: "Today", color: "text-forest", bg: "bg-forest/10 border-forest/20", icon: Info },
  forecast: { label: "Forecast", color: "text-status-watch", bg: "bg-status-watch/10 border-status-watch/20", icon: Clock },
};

export default function ConsoleGuidance() {
  const [updates, setUpdates] = useState<Update[]>([]);
  const [activeTab, setActiveTab] = useState<"all" | "alert" | "today" | "forecast">("all");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getUpdatesClient({ page: 1, limit: 50 })
      .then((result) => setUpdates(result.updates))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const filtered = activeTab === "all" ? updates : updates.filter((u) => u.content_type === activeTab);
  const alerts = updates.filter((u) => u.content_type === "alert");

  const tabs = [
    { key: "all" as const, label: "All", count: updates.length },
    { key: "alert" as const, label: "Alerts", count: alerts.length },
    { key: "today" as const, label: "Today", count: updates.filter((u) => u.content_type === "today").length },
    { key: "forecast" as const, label: "Forecast", count: updates.filter((u) => u.content_type === "forecast").length },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text-charcoal">Guidance</h1>
        <p className="text-text-muted mt-1">Advisory content and actionable guidance for all monitored places.</p>
      </div>

      {/* Active alerts banner */}
      {alerts.length > 0 && (
        <div className="rounded-xl border border-status-emergency/30 bg-status-emergency/5 p-4">
          <div className="flex items-start gap-3">
            <AlertTriangle className="h-5 w-5 text-status-emergency mt-0.5 flex-shrink-0" />
            <div>
              <h3 className="text-sm font-semibold text-status-emergency">
                {alerts.length} active alert{alerts.length !== 1 ? "s" : ""}
              </h3>
              <p className="text-xs text-text-muted mt-1">
                Emergency advisories are in effect. Review and ensure timely delivery.
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 border-b border-background-stone">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab.key
                ? "border-forest text-forest"
                : "border-transparent text-text-muted hover:text-text-charcoal"
            }`}
          >
            {tab.label}
            {tab.count > 0 && (
              <span className="ml-1.5 text-xs bg-background-stone rounded-full px-1.5 py-0.5">{tab.count}</span>
            )}
          </button>
        ))}
      </div>

      {/* Guidance cards */}
      {loading ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Clock className="h-6 w-6 text-text-muted mx-auto mb-2 animate-pulse" />
          <p className="text-text-muted text-sm">Loading guidance...</p>
        </div>
      ) : filtered.length === 0 ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <BookOpen className="h-8 w-8 text-text-muted mx-auto mb-2" />
          <p className="text-text-muted">No guidance content available.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filtered.map((update) => {
            const meta = TYPE_META[update.content_type] || TYPE_META.today;
            const Icon = meta.icon;
            return (
              <div key={update.id} className={`rounded-xl border bg-background p-5 ${meta.bg}`}>
                <div className="flex items-start gap-3">
                  <Icon className={`h-5 w-5 mt-0.5 flex-shrink-0 ${meta.color}`} />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className={`text-xs font-medium ${meta.color}`}>{meta.label}</span>
                      <span className="text-xs text-text-muted">•</span>
                      <span className="text-xs text-text-muted">{update.place_name}</span>
                    </div>
                    <h3 className="text-sm font-semibold text-text-charcoal">{update.headline}</h3>
                    <p className="mt-2 text-sm text-text-muted leading-relaxed">{update.body_text}</p>
                    {update.call_to_action && (
                      <div className="mt-3 p-3 rounded-lg bg-background-mist border border-background-stone">
                        <p className="text-xs font-medium text-text-charcoal">
                          📋 <strong>Action required:</strong> {update.call_to_action}
                        </p>
                      </div>
                    )}
                    <div className="mt-3 flex items-center gap-4 text-xs text-text-muted">
                      <span>Period: {new Date(update.period_start).toLocaleDateString()} — {new Date(update.period_end).toLocaleDateString()}</span>
                      <span>Generated: {new Date(update.generated_at).toLocaleString()}</span>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
