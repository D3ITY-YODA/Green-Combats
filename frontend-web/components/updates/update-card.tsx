import type { PublicUpdate } from "@/types/updates";

interface UpdateCardProps {
  update: PublicUpdate;
}

export function UpdateCard({ update }: UpdateCardProps) {
  return (
    <article className="rounded-2xl border border-background-stone bg-background p-5">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h2 className="text-base font-semibold text-text-charcoal">{update.title}</h2>
          <p className="mt-2 text-sm leading-6 text-text-muted">{update.message}</p>
        </div>
        {(update.type === "warning" || update.type === "emergency") && (
          <span aria-label="Important update" className="h-2.5 w-2.5 rounded-full bg-status-important flex-shrink-0 mt-2" />
        )}
      </div>
      <div className="mt-4 flex items-center justify-between text-xs text-text-muted">
        <span>{update.place_name}</span>
        <span>{update.updated_at}</span>
      </div>
    </article>
  );
}
