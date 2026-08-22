import Link from "next/link";
import type { TopicSection } from "@/types/common";

interface ExploreCardProps {
  section: TopicSection;
}

export function ExploreCard({ section }: ExploreCardProps) {
  return (
    <Link
      href={section.href || "#"}
      className="block rounded-2xl border border-background-stone bg-background p-5 transition hover:border-sage hover:bg-sage-soft"
    >
      <h2 className="font-semibold text-text-charcoal">{section.title}</h2>
      <p className="mt-2 text-sm leading-6 text-text-muted">{section.description}</p>
      <span className="mt-4 inline-block text-sm font-medium text-forest">View details &rarr;</span>
    </Link>
  );
}
