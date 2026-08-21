// app/today/page.tsx
import { getTodayData } from "@/lib/api/today-api";
import { InfoCard } from "@/components/ui/info-card";
import { StatusBadge } from "@/components/ui/status-badge";
import { AlertTriangle, Droplets, Sprout } from "lucide-react";

export default async function TodayPage() {
  // Fetch data from our mock API layer
  const data = await getTodayData();

  return (
    <main className="flex min-h-screen flex-col p-6 md:p-10 max-w-2xl mx-auto">
      {/* Header: Location and Timestamp */}
      <header className="mb-8">
        <h1 className="text-page font-bold text-forest-deep mb-1">Today</h1>
        <p className="text-metadata text-text-muted flex items-center gap-2">
          {data.location} · Updated {data.lastUpdated}
        </p>
      </header>

      {/* Important Updates Section (PDF: "become more visible only when an important update affects the user’s area") */}
      {data.hasImportantUpdates ? (
        <div className="mb-8 rounded-xl border border-status-emergency/30 bg-status-emergency/5 p-5">
          <div className="flex items-start gap-3">
            <AlertTriangle className="h-5 w-5 text-status-emergency mt-0.5 flex-shrink-0" />
            <div>
              <div className="flex items-center gap-2 mb-1">
                <StatusBadge status="emergency" label="Emergency" />
                <span className="text-card font-semibold text-text-charcoal">Flood warning</span>
              </div>
              <p className="text-body text-text-charcoal">
                Heavy rainfall may cause flooding within 24 hours.
              </p>
              <p className="text-metadata text-text-muted mt-2">Updated {data.lastUpdated}</p>
            </div>
          </div>
        </div>
      ) : (
        // Calm state (PDF: "The application should remain subtle during ordinary conditions")
        <div className="mb-8 rounded-xl border border-status-normal/20 bg-status-normal/5 p-5 flex items-center gap-3">
          <Sprout className="h-5 w-5 text-status-normal flex-shrink-0" />
          <p className="text-body text-status-normal font-medium">
            No important updates for your area.
          </p>
        </div>
      )}

      {/* Outlook Cards */}
      <div className="grid gap-4">
        <InfoCard title="Local outlook">
          {data.localOutlook}
        </InfoCard>

        <InfoCard title="Water outlook">
          <div className="flex items-start gap-3">
            <Droplets className="h-5 w-5 text-sky mt-0.5 flex-shrink-0" />
            <span>{data.waterOutlook}</span>
          </div>
        </InfoCard>
      </div>
    </main>
  );
}
