export default function ProfileDetailPage() {
  const title = "organizations".replace(/\b\w/g, l => l.toUpperCase());
  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-2xl font-bold text-text-charcoal mb-6"></h1>
      <div className="rounded-xl border border-background-stone bg-background p-6">
        <p className="text-text-muted">Manage your  preferences here. Changes are saved automatically.</p>
      </div>
    </div>
  );
}
