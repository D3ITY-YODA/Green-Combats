// app/explore/page.tsx
import { getExploreSections } from "@/lib/api/explore-api";
import { InfoCard } from "@/components/ui/info-card";
import { Sprout, Droplets, Calendar, Users, Shield } from "lucide-react";

// Helper to map the mock API icon strings to actual Lucide icons
const iconMap = {
  local: Sprout,
  water: Droplets,
  seasonal: Calendar,
  community: Users,
  preparedness: Shield,
};

export default async function ExplorePage() {
  // Fetch the relevant sections from our mock API
  const sections = await getExploreSections();

  return (
    <main className="flex min-h-screen flex-col p-6 md:p-10 max-w-4xl mx-auto">
      {/* Header */}
      <header className="mb-8">
        <h1 className="text-page font-bold text-forest-deep mb-1">Explore</h1>
        <p className="text-metadata text-text-muted">
          Relevant information for your area
        </p>
      </header>

      {/* Grid of Relevant Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {sections.map((section) => {
          const IconComponent = iconMap[section.icon] || Sprout;
          
          return (
            <InfoCard key={section.id} title={section.title} className="flex flex-col gap-3">
              <div className="flex items-center gap-3">
                <IconComponent className="h-5 w-5 text-forest" />
                <p className="text-body text-text-charcoal leading-relaxed">
                  {section.description}
                </p>
              </div>
            </InfoCard>
          );
        })}
      </div>

      {/* Empty State (PDF: "Users do not see features that are irrelevant") */}
      {sections.length === 0 && (
        <div className="rounded-xl border border-background-stone bg-background-mist p-8 text-center">
          <p className="text-body text-text-muted">
            No additional information is currently available for your selected location.
          </p>
        </div>
      )}
    </main>
  );
}
