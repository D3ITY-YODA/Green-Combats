// components/updates/update-card.tsx

interface UpdateCardProps {
  update: {
    id: string;
    title?: string;
    message?: string;
    headline?: string;
    body_text?: string;
    status?: string;
    priority?: string;
  };
}

export function UpdateCard({ update }: UpdateCardProps) {
  return (
    <div className="rounded-2xl border border-stone bg-white p-5">
      <h3 className="font-semibold text-charcoal">
        {update.title ?? update.headline ?? "Update"}
      </h3>
      <p className="mt-2 text-sm leading-6 text-muted">
        {update.message ?? update.body_text ?? ""}
      </p>
    </div>
  );
}
