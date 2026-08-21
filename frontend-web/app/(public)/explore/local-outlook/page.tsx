export default function ExploreDetailPage() {
  const title = "local-outlook".replace('-', ' ').replace(/\b\w/g, l => l.toUpperCase());
  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <header className="mb-8">
        <p className="text-sm text-text-muted">Explore / </p>
        <h1 className="mt-1 text-3xl font-semibold text-text-charcoal"></h1>
      </header>
      <div className="rounded-2xl border border-background-stone bg-background p-6">
        <p className="text-text-muted">Detailed information and historical data for  will be displayed here. This view provides deep insights relevant to Lower Valley.</p>
      </div>
    </div>
  );
}
