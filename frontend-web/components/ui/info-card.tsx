// components/ui/info-card.tsx
import { cn } from "@/lib/utils";

export function InfoCard({ title, children, className }: { title: string; children: React.ReactNode; className?: string }) {
  return (
    <div className={cn(
      "rounded-xl border border-background-stone bg-background-mist p-5",
      className
    )}>
      <h3 className="text-card font-semibold text-forest-deep mb-2">{title}</h3>
      <div className="text-body text-text-charcoal leading-relaxed">{children}</div>
    </div>
  );
}
