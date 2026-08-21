import { ExploreCard } from "@/components/explore/explore-card";
import type { TopicSection } from "@/types/common";

export default async function ExplorePage() {
  const sections: TopicSection[] = [
    { key: "local_outlook", title: "Local outlook", description: "Overview of current conditions.", available: true, href: "/explore/local-outlook" },
    { key: "water_outlook", title: "Water outlook", description: "Water levels, availability and quality.", available: true, href: "/explore/water" },
    { key: "seasonal_information", title: "Seasonal information", description: "What to expect in the coming months.", available: true, href: "/explore/seasonal" },
    { key: "community_updates", title: "Community updates", description: "Reports from your neighbors.", available: true, href: "/explore/community" }
  ];

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <header className="mb-8">
        <p className="text-sm text-text-muted">Explore</p>
        <h1 className="mt-1 text-3xl font-semibold text-text-charcoal">Lower Valley</h1>
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
