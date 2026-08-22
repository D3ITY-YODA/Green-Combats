"use client";

import { useState, useEffect } from "react";
import { Activity, Clock, Filter, CheckCircle, XCircle, FileText, UserPlus, Settings, Eye } from "lucide-react";
import { getAuditEntries } from "@/lib/api/console";
import type { AuditEntry } from "@/lib/api/console";

const ACTION_META: Record<string, { label: string; color: string; bg: string; icon: React.ElementType }> = {
  observation_verified: { label: "Verified", color: "text-status-normal", bg: "bg-status-normal/10", icon: CheckCircle },
  observation_rejected: { label: "Rejected", color: "text-status-emergency", bg: "bg-status-emergency/10", icon: XCircle },
  content_published: { label: "Published", color: "text-forest", bg: "bg-forest/10", icon: FileText },
  user_created: { label: "User Created", color: "text-sky", bg: "bg-sky/10", icon: UserPlus },
  user_deactivated: { label: "User Deactivated", color: "text-text-muted", bg: "bg-background-stone", icon: UserPlus },
  org_created: { label: "Org Created", color: "text-forest", bg: "bg-forest/10", icon: Settings },
  org_member_added: { label: "Member Added", color: "text-sky", bg: "bg-sky/10", icon: UserPlus },
  org_member_removed: { label: "Member Removed", color: "text-status-watch", bg: "bg-status-watch/10", icon: UserPlus },
  source_enabled: { label: "Source Enabled", color: "text-status-normal", bg: "bg-status-normal/10", icon: Settings },
  source_disabled: { label: "Source Disabled", color: "text-status-watch", bg: "bg-status-watch/10", icon: Settings },
  config_changed: { label: "Config Changed", color: "text-status-watch", bg: "bg-status-watch/10", icon: Settings },
};

function relativeTime(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 1) return "just now";
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffH = Math.floor(diffMin / 60);
  if (diffH < 24) return `${diffH}h ago`;
  const diffD = Math.floor(diffH / 24);
  return `${diffD}d ago`;
}

export default function ConsoleActivity() {
  const [entries, setEntries] = useState<AuditEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<string>("");

  useEffect(() => {
    setLoading(true);
    const params: { page?: number; limit?: number; action?: string } = { page: 1, limit: 50 };
    if (filter) params.action = filter;
    getAuditEntries(params)
      .then((result) => setEntries(result.entries))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [filter]);

  const filters = [
    { key: "", label: "All" },
    { key: "observation_verified", label: "Verified" },
    { key: "observation_rejected", label: "Rejected" },
    { key: "content_published", label: "Published" },
    { key: "user_created", label: "User Created" },
    { key: "config_changed", label: "Config" },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text-charcoal">Activity</h1>
        <p className="text-text-muted mt-1">Audit log of all system and user actions.</p>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-2 flex-wrap">
        <Filter className="h-4 w-4 text-text-muted" />
        {filters.map((f) => (
          <button
            key={f.key}
            onClick={() => setFilter(f.key)}
            className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
              filter === f.key
                ? "bg-forest text-white"
                : "bg-background-stone text-text-muted hover:text-text-charcoal"
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {/* Activity feed */}
      {loading ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Clock className="h-6 w-6 text-text-muted mx-auto mb-2 animate-pulse" />
          <p className="text-text-muted text-sm">Loading activity...</p>
        </div>
      ) : entries.length === 0 ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Activity className="h-10 w-10 text-text-muted mx-auto mb-3" />
          <h2 className="text-lg font-semibold text-text-charcoal">No activity yet</h2>
          <p className="text-sm text-text-muted mt-2">
            Audit entries will appear here as users interact with the platform.
          </p>
        </div>
      ) : (
        <div className="bg-background rounded-xl border border-background-stone overflow-hidden">
          <div className="divide-y divide-background-stone">
            {entries.map((entry) => {
              const meta = ACTION_META[entry.action] || { label: entry.action, color: "text-text-muted", bg: "bg-background-stone", icon: Activity };
              const Icon = meta.icon;
              return (
                <div key={entry.id} className="px-5 py-3 flex items-center gap-3 hover:bg-background-mist/50">
                  <div className={`p-2 rounded-lg ${meta.bg}`}>
                    <Icon className={`h-4 w-4 ${meta.color}`} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className={`text-xs font-medium ${meta.color}`}>{meta.label}</span>
                      <span className="text-xs text-text-muted">•</span>
                      <span className="text-xs text-text-muted">{entry.resource_type}</span>
                    </div>
                    {entry.details && Object.keys(entry.details).length > 0 && (
                      <p className="text-xs text-text-muted mt-0.5 truncate">
                        {JSON.stringify(entry.details)}
                      </p>
                    )}
                  </div>
                  <div className="text-right flex-shrink-0">
                    <p className="text-xs text-text-muted">{relativeTime(entry.created_at)}</p>
                    <p className="text-xs text-text-muted/50">{new Date(entry.created_at).toLocaleDateString()}</p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
