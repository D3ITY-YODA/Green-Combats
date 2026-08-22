import { UpdateCard } from "@/components/updates/update-card";
import { EmptyState } from "@/components/states/empty-state";
import { getUpdates } from "@/lib/api/updates-server";
import { adaptUpdateToPublicUpdate } from "@/lib/api/updates";
import type { PublicUpdate } from "@/types/updates";

// Fallback mock data
function getFallbackUpdates(): PublicUpdate[] {
  return [
    { id: "1", place_id: "1", topic_key: "water_outlook", type: "warning", priority: "important", title: "Water level rising", message: "Water level rising in the river.", valid_from: "2026-08-22", updated_at: "3h ago", status: "published", display_on_today: true, display_in_feed: true, locale: "en", place_name: "Lower Valley" },
    { id: "2", place_id: "1", topic_key: "local_outlook", type: "information", priority: "normal", title: "Heavy rainfall last night", message: "Reported in West Kano.", valid_from: "2026-08-22", updated_at: "8h ago", status: "published", display_on_today: false, display_in_feed: true, locale: "en", place_name: "Lower Valley" }
  ];
}

export default async function UpdatesPage() {
  let updates: PublicUpdate[];

  try {
    const result = await getUpdates({ page: 1, limit: 20 });
    updates = result.updates.map(adaptUpdateToPublicUpdate);
  } catch {
    updates = getFallbackUpdates();
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <header className="mb-8">
        <h1 className="text-3xl font-semibold text-text-charcoal">Updates</h1>
      </header>
      {updates.length === 0 ? (
        <EmptyState title="No important updates" message="There are no important updates for your selected places." />
      ) : (
        <div className="space-y-4">
          {updates.map((update) => (
            <UpdateCard key={update.id} update={update} />
          ))}
        </div>
      )}
    </div>
  );
}
