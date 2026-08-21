// components/explore/explore-card.tsx

import Link from "next/link";
import type { TopicSection } from "@/types/context"; // Adjust path if TopicSection is located elsewhere

interface ExploreCardProps {
  section: TopicSection;
  href: string;
}

export function ExploreCard({ section, href }: ExploreCardProps) {
  // Handle topics that are not yet available for the selected place
  if (!section.available) {
    return (
      <div className="block rounded-2xl border border-stone bg-white p-5 opacity-60">
        <h2 className="font-semibold text-charcoal">{section.title}</h2>
        <p className="mt-2 text-sm leading-6 text-muted">
          {section.description}
        </p>
        <span className="mt-4 inline-block text-sm font-medium text-muted">
          Not available yet
        </span>
      </div>
    );
  }

  // Use Next.js Link for internal navigation and performance
  return (
    <Link
      href={href}
      className="block rounded-2xl border border-stone bg-white p-5 transition hover:border-sage hover:bg-soft-sage focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2"
    >
      <h2 className="font-semibold text-charcoal">{section.title}</h2>
      <p className="mt-2 text-sm leading-6 text-muted">
        {section.description}
      </p>
      <span className="mt-4 inline-block text-sm font-medium text-forest">
        View
      </span>
    </Link>
  );
}
