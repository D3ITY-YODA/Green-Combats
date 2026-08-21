// app/updates/page.tsx
import { getUpdates } from "@/lib/api/updates-api";
import { StatusBadge } from "@/components/ui/status-badge";
import { Clock } from "lucide-react";

export default async function UpdatesPage() {
  const updates = await getUpdates();

  return (
    <main className="flex min-h-screen flex-col p-6 md:p-10 max-w-2xl mx-auto">
      {/* Header */}
      <header className="mb-8">
        <h1 className="text-page font-bold text-forest-deep mb-1">Updates</h1>
        <p className="text-metadata text-text-muted">
          Important changes affecting your area
        </p>
      </header>

      {/* Updates List */}
      <div className="space-y-4">
        {updates.map((update) => (
          <article 
            key={update.id} 
            className="rounded-xl border border-background-stone bg-background p-5 transition-shadow hover:shadow-sm"
          >
            <div className="flex items-center justify-between mb-2">
              <StatusBadge status={update.status} label={update.status.charAt(0).toUpperCase() + update.status.slice(1)} />
              <span className="text-metadata text-text-muted flex items-center gap-1">
                <Clock className="h-3 w-3" /> {update.updated}
              </span>
            </div>
            
            <h3 className="text-card font-semibold text-text-charcoal mb-1">
              {update.title}
            </h3>
            
            <p className="text-body text-text-muted leading-relaxed">
              {update.description}
            </p>
            
            <p className="text-metadata text-text-muted mt-3 pt-3 border-t border-background-stone">
              {update.location}
            </p>
          </article>
        ))}
      </div>

      {/* Empty State */}
      {updates.length === 0 && (
        <div className="rounded-xl border border-background-stone bg-background-mist p-8 text-center">
          <p className="text-body text-text-muted">
            There are no important updates for your selected locations right now.
          </p>
        </div>
      )}
    </main>
  );
}
