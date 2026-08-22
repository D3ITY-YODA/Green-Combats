"use client";

import { useState, useEffect } from "react";
import { getPendingObservations } from "@/lib/api/observations";
import type { Observation } from "@/types/reports";

const statusColors: Record<string, string> = {
  pending: "bg-status-watch/10 text-status-watch",
  verified: "bg-status-normal/10 text-status-normal",
  rejected: "bg-status-emergency/10 text-status-emergency",
  flagged: "bg-status-important/10 text-status-important",
};

const categoryLabels: Record<string, string> = {
  flood: "Flooding visible",
  drought: "Unusually dry",
  water_quality: "Water has changed",
  crop_damage: "Vegetation stress",
  air_quality: "Air quality concern",
  other: "Other observation",
};

export default function ConsoleReports() {
  const [reports, setReports] = useState<Observation[]>([]);
  const [activeTab, setActiveTab] = useState("pending");

  useEffect(() => {
    getPendingObservations(50).then(setReports).catch(() => {});
  }, []);

  const tabs = ["pending", "verified", "rejected", "flagged"];
  const filtered = reports.filter(r => r.status === activeTab);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-text-charcoal">Community reports</h1>

      <div className="flex gap-2 border-b border-background-stone">
        {tabs.map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab
                ? "border-forest text-forest"
                : "border-transparent text-text-muted hover:text-text-charcoal"
            }`}
          >
            {tab.charAt(0).toUpperCase() + tab.slice(1)}
          </button>
        ))}
      </div>

      <div className="bg-background rounded-xl border border-background-stone overflow-hidden">
        <table className="w-full text-left text-sm">
          <thead className="bg-background-mist border-b border-background-stone text-text-muted">
            <tr>
              <th className="p-4 font-medium">Type</th>
              <th className="p-4 font-medium">Description</th>
              <th className="p-4 font-medium">Submitted</th>
              <th className="p-4 font-medium">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background-stone">
            {filtered.length === 0 ? (
              <tr>
                <td colSpan={4} className="p-8 text-center text-text-muted">
                  No {activeTab} reports
                </td>
              </tr>
            ) : (
              filtered.map((report) => (
                <tr key={report.id} className="hover:bg-background-mist/50">
                  <td className="p-4 text-text-charcoal font-medium">
                    {categoryLabels[report.category] || report.category}
                  </td>
                  <td className="p-4 text-text-muted max-w-xs truncate">{report.description}</td>
                  <td className="p-4 text-text-muted">
                    {new Date(report.created_at).toLocaleDateString()}
                  </td>
                  <td className="p-4">
                    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${statusColors[report.status] || "bg-background-stone text-text-muted"}`}>
                      {report.status}
                    </span>
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
