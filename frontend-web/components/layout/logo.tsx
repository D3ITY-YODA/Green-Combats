// components/layout/logo.tsx
import { Sprout } from "lucide-react";

export function Logo() {
  return (
    <div className="hidden md:flex items-center gap-3 px-4 py-6 mb-4 border-b border-white/10">
      <div className="p-2 rounded-lg bg-white/10">
        <Sprout className="h-8 w-8 text-white" />
      </div>
      <div>
        <h1 className="text-lg font-bold text-white leading-tight">Green Compass</h1>
        <p className="text-xs text-sage-soft/80">Environmental Guidance</p>
      </div>
    </div>
  );
}
