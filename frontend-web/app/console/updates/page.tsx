"use client";

import { useState, useEffect } from "react";
import { Plus } from "lucide-react";
import { getUpdatesClient } from "@/lib/api/updates";
import { adaptUpdateToPublicUpdate } from "@/lib/api/updates";
import type { Update } from "@/types/updates";

const statusColors: Record<string, string> = {
  today: "bg-status-normal/10 text-status-normal",
  forecast: "bg-status-watch/10 text-status-watch",
  alert: "bg-status-important/10 text-status-important",
};

export default function ConsoleUpdates() {
  const [updates, setUpdates] = useState<Update[]>([]);

  useEffect(() => {
    getUpdatesClient({ page: 1, limit: 50 })
      .then(result => setUpdates(result.updates))
      .catch(() => {});
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-text-charcoal">Updates</h1>
        <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-forest text-white text-sm font-medium hover:bg-forest-deep">
          <Plus className="h-4 w-4" /> Create update
        </button>
      </div>

      <div className="bg-background rounded-xl border border-background-stone overflow-hidden">
        <table className="w-full text-left text-sm">
          <thead className="bg-background-mist border-b border-background-stone text-text-muted">
            <tr>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 font-medium">Type</th>
              <th className="p-4 font-medium">Place</th>
              <th className="p-4 font-medium">Updated</th>
              <th className="p-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background-stone">
            {updates.length === 0 ? (
              <tr>
                <td colSpan={5} className="p-8 text-center text-text-muted">
                  No updates found
                </td>
              </tr>
            ) : (
              updates.map((update) => (
                <tr key={update.id} className="hover:bg-background-mist/50">
                  <td className="p-4">
                    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${statusColors[update.content_type] || "bg-background-stone text-text-muted"}`}>
                      {update.content_type}
                    </span>
                  </td>
                  <td className="p-4 text-text-charcoal font-medium">{update.headline}</td>
                  <td className="p-4 text-text-muted">{update.place_name}</td>
                  <td className="p-4 text-text-muted">{new Date(update.generated_at).toLocaleDateString()}</td>
                  <td className="p-4 text-right">
                    <button className="text-forest hover:underline text-xs font-medium">Edit</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
