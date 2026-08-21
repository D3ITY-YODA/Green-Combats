export default function PublicHelpPage() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <h1 className="text-2xl font-bold text-text-charcoal mb-6">Help & Support</h1>
      <div className="space-y-4">
        <div className="rounded-xl border border-background-stone bg-background p-6">
          <h2 className="text-lg font-semibold text-forest-deep">Using the Today View</h2>
          <p className="text-text-muted mt-2">The Today view shows the most critical updates for your selected location. If there are no alerts, you will see a calm "No important updates" message.</p>
        </div>
        <div className="rounded-xl border border-background-stone bg-background p-6">
          <h2 className="text-lg font-semibold text-forest-deep">Offline Access</h2>
          <p className="text-text-muted mt-2">Green Compass caches your latest updates. If you lose connectivity, a banner will appear, and you can still view your last downloaded data.</p>
        </div>
      </div>
    </div>
  );
}
