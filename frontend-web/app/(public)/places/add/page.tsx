export default function AddPlacePage() {
  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-2xl font-bold text-text-charcoal mb-6">Add a new place</h1>
      <div className="rounded-xl border border-background-stone bg-background p-6 space-y-4">
        <div>
          <label className="block text-sm font-medium text-text-charcoal mb-2">Search for a place</label>
          <input type="text" className="w-full rounded-lg border border-background-stone bg-background p-3 text-text-charcoal focus:border-forest focus:ring-1 focus:ring-forest" placeholder="Enter village, district, or region..." />
        </div>
        <button className="w-full rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-forest-deep">Search</button>
      </div>
    </div>
  );
}
