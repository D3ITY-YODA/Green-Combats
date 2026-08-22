import { TodayStatus } from "@/components/today/today-status";
import { UpdateCard } from "@/components/updates/update-card";
import { ExploreCard } from "@/components/explore/explore-card";
import type { TodayResponse } from "@/types/updates";

// Mock data fetcher (Replace with real API call later)
async function getTodayData(): Promise<TodayResponse> {
  return {
    place: { id: "1", name: "Lower Valley", type: "ward", country_code: "KE" },
    status: { title: "Good morning", message: "Here's what's happening in Lower Valley.", updated_at: "10:00", data_status: "current" },
    updates: [
      { id: "1", place_id: "1", topic_key: "water_outlook", type: "information", priority: "normal", title: "Water conditions normal", message: "No significant changes in the last 24h.", valid_from: "2026-08-22", updated_at: "2h ago", source_name: "Local Water Authority", status: "published", display_on_today: true, display_in_feed: true, locale: "en", place_name: "Lower Valley" }
    ],
    sections: [
      { key: "local_outlook", title: "Local outlook", description: "Conditions for the coming days.", available: true, href: "/explore/local-outlook" },
      { key: "water_outlook", title: "Water outlook", description: "Water levels and availability.", available: true, href: "/explore/water" }
    ],
    updated_at: "10:00"
  };
}

export default async function TodayPage() {
  const today = await getTodayData();

  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <header className="mb-8">
        <p className="text-sm text-text-muted">Today</p>
        <h1 className="mt-1 text-3xl font-semibold text-text-charcoal">{today.place.name}</h1>
      </header>

      <TodayStatus
        title={today.status.title}
        message={today.status.message}
        updatedAt={today.status.updated_at}
        dataStatus={today.status.data_status}
      />

      {today.updates.length > 0 && (
        <section className="mt-8 space-y-4">
          <h2 className="text-lg font-semibold text-text-charcoal">Key updates</h2>
          {today.updates.map((update) => (
            <UpdateCard key={update.id} update={update} />
          ))}
        </section>
      )}

      <section className="mt-8 grid gap-4 md:grid-cols-2">
        {today.sections.map((section) => (
          <ExploreCard key={section.key} section={section} />
        ))}
      </section>
    </div>
  );
}
