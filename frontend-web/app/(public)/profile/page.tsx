"use client";
import { useState, useEffect } from "react";
import { MapPin, Bell, Globe, Shield, HelpCircle, ChevronRight, Plus, Accessibility } from "lucide-react";
import { getMe } from "@/lib/api/auth";
import { getSavedPlacesClient } from "@/lib/api/places";
import type { SessionUser } from "@/types/auth";
import type { SavedPlace } from "@/types/common";

export default function ProfilePage() {
  const [user, setUser] = useState<SessionUser | null>(null);
  const [places, setPlaces] = useState<SavedPlace[]>([]);
  const [language, setLanguage] = useState("English");
  const [notifications, setNotifications] = useState(true);

  useEffect(() => {
    getMe().then(setUser).catch(() => {});
    getSavedPlacesClient().then(setPlaces).catch(() => {});
  }, []);

  const displayName = user?.display_name || "Guest User";
  const primaryPlace = places.find(p => p.is_primary) || places[0];

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <header className="mb-8">
        <p className="text-sm text-text-muted">Profile</p>
        <h1 className="mt-1 text-3xl font-semibold text-text-charcoal">{displayName}</h1>
        <p className="mt-1 text-text-muted">{primaryPlace?.name || "No place selected"}</p>
      </header>

      <div className="space-y-8">
        <section>
          <h2 className="text-lg font-semibold text-text-charcoal mb-3">Saved places</h2>
          <div className="space-y-2">
            {places.length > 0 ? places.map((place) => (
              <div key={place.id} className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
                <div className="flex items-center gap-3">
                  <MapPin className="h-5 w-5 text-forest" />
                  <span className="text-base text-text-charcoal">{place.name}</span>
                </div>
                {place.is_primary && <span className="text-xs text-status-normal font-medium">Current place</span>}
              </div>
            )) : (
              <div className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
                <div className="flex items-center gap-3"><MapPin className="h-5 w-5 text-forest" /><span className="text-base text-text-charcoal">Lower Valley</span></div>
                <span className="text-xs text-status-normal font-medium">Current place</span>
              </div>
            )}
            <button className="w-full flex items-center justify-center gap-2 rounded-xl border border-dashed border-background-stone bg-background-mist p-4 text-base text-text-muted hover:bg-background-stone">
              <Plus className="h-4 w-4" /> Add place
            </button>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-text-charcoal mb-3">Preferences</h2>
          <div className="space-y-2">
            <div className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
              <div className="flex items-center gap-3"><Globe className="h-5 w-5 text-forest" /><span className="text-base text-text-charcoal">Language</span></div>
              <select value={language} onChange={(e) => setLanguage(e.target.value)} className="text-base text-text-muted bg-transparent focus:outline-none cursor-pointer text-right">
                <option>English</option><option>Swahili</option><option>French</option>
              </select>
            </div>
            <div className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
              <div className="flex items-center gap-3"><Bell className="h-5 w-5 text-forest" /><span className="text-base text-text-charcoal">Notifications</span></div>
              <button onClick={() => setNotifications(!notifications)} className={`w-12 h-6 rounded-full p-1 transition-colors ${notifications ? 'bg-forest' : 'bg-background-stone'}`} aria-label="Toggle notifications">
                <div className={`w-4 h-4 rounded-full bg-white transition-transform ${notifications ? 'translate-x-6' : 'translate-x-0'}`} />
              </button>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-text-charcoal mb-3">Support</h2>
          <div className="space-y-2">
            {[
              { icon: Shield, label: "Privacy & Security" },
              { icon: HelpCircle, label: "Help & About" },
              { icon: Accessibility, label: "Accessibility" }
            ].map((item) => (
              <button key={item.label} className="w-full flex items-center justify-between rounded-xl border border-background-stone bg-background p-4 hover:bg-background-mist transition-colors text-left">
                <div className="flex items-center gap-3"><item.icon className="h-5 w-5 text-forest" /><span className="text-base text-text-charcoal">{item.label}</span></div>
                <ChevronRight className="h-4 w-4 text-text-muted" />
              </button>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
