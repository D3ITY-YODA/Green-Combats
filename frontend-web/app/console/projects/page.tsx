"use client";

import { useState, useEffect } from "react";
import { FolderKanban, Plus, Clock, CheckCircle, Archive } from "lucide-react";
import { getProjects } from "@/lib/api/console";
import type { Project } from "@/lib/api/console";

const STATUS_META: Record<string, { label: string; color: string; bg: string; icon: React.ElementType }> = {
  active: { label: "Active", color: "text-status-normal", bg: "bg-status-normal/10 border-status-normal/20", icon: CheckCircle },
  archived: { label: "Archived", color: "text-text-muted", bg: "bg-background-stone text-text-muted", icon: Archive },
};

export default function ConsoleProjects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    // Note: requires org_id. Without an org, we show an empty state with guidance.
    getProjects("00000000-0000-0000-0000-000000000000", 1, 50)
      .then((result) => setProjects(result.projects))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-text-charcoal">Projects</h1>
          <p className="text-text-muted mt-1">Manage institutional projects and programs.</p>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-forest text-white text-sm font-medium hover:bg-forest-deep">
          <Plus className="h-4 w-4" /> New Project
        </button>
      </div>

      {loading ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Clock className="h-6 w-6 text-text-muted mx-auto mb-2 animate-pulse" />
          <p className="text-text-muted text-sm">Loading projects...</p>
        </div>
      ) : projects.length === 0 ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <FolderKanban className="h-10 w-10 text-text-muted mx-auto mb-3" />
          <h2 className="text-lg font-semibold text-text-charcoal">No projects yet</h2>
          <p className="text-sm text-text-muted mt-2 max-w-md mx-auto">
            Projects help you organize monitoring campaigns, community programs, and institutional initiatives.
          </p>
          <button className="mt-4 inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-forest text-white text-sm font-medium hover:bg-forest-deep">
            <Plus className="h-4 w-4" /> Create your first project
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          {projects.map((project) => {
            const meta = STATUS_META[project.status] || STATUS_META.active;
            const Icon = meta.icon;
            return (
              <div key={project.id} className="rounded-xl border border-background-stone bg-background p-5 hover:bg-background-mist transition-colors">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="text-sm font-semibold text-text-charcoal">{project.name}</h3>
                      <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium border ${meta.bg} ${meta.color}`}>
                        <Icon className="h-3 w-3" /> {meta.label}
                      </span>
                    </div>
                    {project.description && (
                      <p className="text-sm text-text-muted mt-1 line-clamp-2">{project.description}</p>
                    )}
                    <div className="flex items-center gap-4 mt-3 text-xs text-text-muted">
                      {project.start_date && (
                        <span>Start: {new Date(project.start_date).toLocaleDateString()}</span>
                      )}
                      {project.end_date && (
                        <span>End: {new Date(project.end_date).toLocaleDateString()}</span>
                      )}
                      <span>Created: {new Date(project.created_at).toLocaleDateString()}</span>
                    </div>
                  </div>
                  <button className="text-xs font-medium text-forest hover:underline">View</button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
