interface EmptyStateProps {
  title: string;
  message: string;
  action?: React.ReactNode;
}

export function EmptyState({ title, message, action }: EmptyStateProps) {
  return (
    <section className="rounded-2xl border border-background-stone bg-background p-8 text-center">
      <h2 className="text-lg font-semibold text-text-charcoal">{title}</h2>
      <p className="mt-2 text-sm text-text-muted">{message}</p>
      {action && <div className="mt-5">{action}</div>}
    </section>
  );
}
