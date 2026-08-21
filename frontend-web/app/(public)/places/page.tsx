import Link from "next/link";
import { Plus, MapPin } from "lucide-react";

export default function PlacesPage() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-semibold text-text-charcoal">Saved Places</h1>
        <Link href="/places/add" className="flex items-center gap-2 px-4 py-2 rounded-lg bg-forest text-white text-sm font-medium hover:bg-forest-deep">
          <Plus className="h-4 w-4" /> Add Place
        </Link>
      </div>
      <div className="space-y-3">
        <Link href="/places/1" className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4 hover:bg-background-mist transition-colors">
          <div className="flex items-center gap-3"><MapPin className="h-5 w-5 text-forest" /><span className="text-base font-medium text-text-charcoal">Lower Valley</span></div>
          <span className="text-xs text-status-normal font-medium">Primary</span>
        </Link>
      </div>
    </div>
  );
}
