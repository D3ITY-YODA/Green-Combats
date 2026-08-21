import Link from "next/link";
export default function ChoosePlacePage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-forest-deep text-center">Choose a place</h1>
      <p className="text-text-muted text-center">See updates for a place that matters to you.</p>
      <div className="space-y-3">
        <input type="text" placeholder="Search for a place..." className="w-full rounded-lg border border-background-stone bg-background p-3" />
        <div className="space-y-2">
          {["Lower Valley", "Riverside", "Upper Highland"].map(p => (
            <button key={p} className="w-full text-left p-4 rounded-xl border border-background-stone bg-background hover:bg-background-mist">{p}</button>
          ))}
        </div>
      </div>
      <Link href="/setup-complete" className="block w-full text-center rounded-xl bg-forest py-3 font-medium text-white hover:bg-forest-deep">Continue</Link>
    </div>
  );
}
