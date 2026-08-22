import { ExploreCard } from "@/components/explore/explore-card";
import { getExploreIndicators } from "@/lib/api/updates-server";
import { adaptExploreToTopicSections } from "@/lib/api/updates";
import type { TopicSection } from "@/types/common";

export const dynamic = "force-dynamic";

// Fallback mock data
function getFallbackSections(): TopicSection[] {
  return [
    { key: "local_outlook", title: "Local outlook", description: "Overview of current conditions.", available: true, href: "/explore/local-outlook" },
    { key: "water_outlook", title: "Water outlook", description: "Water levels, availability and quality.", available: true, href: "/explore/water" },
    { key: "seasonal_information", title: "Seasonal information", description: "What to expect in the coming months.", available: true, href: "/explore/seasonal" },
    { key: "community_updates", title: "Community updates", description: "Reports from your neighbors.", available: true, href: "/explore/community" }
  ];
}

export default async function ExplorePage() {
  let sections: TopicSection[];
  let placeName = "Lower Valley";

  try {
    const result = await getExploreIndicators({ limit: 50 });
    sections = adaptExploreToTopicSections(result.indicators);
    // Use the first indicator's place name if available
    if (result.indicators.length > 0 && result.indicators[0].place_name) {
      placeName = result.indicators[0].place_name;
    }
  } catch {
    sections = getFallbackSections();
  }

  // Always ensure at least the core sections are available
  const fallbackSections = getFallbackSections();
  for (const fb of fallbackSections) {
    if (!sections.find(s => s.key === fb.key)) {
      sections.push(fb);
    }
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <header className="mb-8">
        <p className="text-sm text-text-muted">Explore</p>
        <h1 className="mt-1 text-3xl font-semibold text-text-charcoal">{placeName}</h1>
        <p className="mt-3 text-text-muted">Information relevant to your selected place.</p>
      </header>
      <div className="grid gap-4 md:grid-cols-2">
        {sections.map((section) => (
          <ExploreCard key={section.key} section={section} />
        ))}
      </div>
    </div>
  );
}
