"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Search, MapPin, Loader2 } from "lucide-react";
import { getSavedPlacesClient, searchNearbyPlaces, savePlace, setPrimaryPlace } from "@/lib/api/places";
import type { SavedPlace, NearbyPlace } from "@/types/places";

export default function ChoosePlacePage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");
  const [savedPlaces, setSavedPlaces] = useState<SavedPlace[]>([]);
  const [nearbyPlaces, setNearbyPlaces] = useState<NearbyPlace[]>([]);
  const [selectedPlaceId, setSelectedPlaceId] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [searching, setSearching] = useState(false);
  const [saving, setSaving] = useState(false);

  // Load saved places on mount
  useEffect(() => {
    getSavedPlacesClient()
      .then((places) => {
        setSavedPlaces(places);
        if (places.length > 0) {
          const primary = places.find((p) => p.is_primary) || places[0];
          setSelectedPlaceId(primary.id);
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  // Search nearby places when typing (debounced)
  useEffect(() => {
    if (!searchQuery.trim()) {
      setNearbyPlaces([]);
      return;
    }
    setSearching(true);
    const timer = setTimeout(() => {
      // Try to use geolocation for nearby search, fallback to a default location (Nairobi)
      const defaultLat = -1.2921;
      const defaultLon = 36.8219;
      searchNearbyPlaces(defaultLat, defaultLon, 50000, 20)
        .then((places) => setNearbyPlaces(places))
        .catch(() => setNearbyPlaces([]))
        .finally(() => setSearching(false));
    }, 500);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  const handleSelect = async (placeId: string) => {
    setSelectedPlaceId(placeId);
    setSaving(true);
    try {
      // If it's not already saved, save it first
      const isAlreadySaved = savedPlaces.some((p) => p.id === placeId);
      if (!isAlreadySaved) {
        await savePlace(placeId);
      }
      await setPrimaryPlace(placeId);
      router.push("/today");
    } catch {
      // User can still continue even if save fails
      router.push("/today");
    } finally {
      setSaving(false);
    }
  };

  const handleSkip = () => {
    router.push("/today");
  };

  // Filter nearby places by search query
  const filteredNearby = nearbyPlaces.filter(
    (p) => p.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold text-forest-deep">Choose a place</h1>
        <p className="text-text-muted mt-2">See updates for a place that matters to you.</p>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-text-muted" />
        <input
          type="text"
          placeholder="Search for a place..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-lg border border-background-stone bg-background p-3 pl-10 text-text-charcoal focus:border-forest focus:ring-1 focus:ring-forest"
        />
      </div>

      {/* Saved places */}
      {savedPlaces.length > 0 && (
        <div className="space-y-2">
          <p className="text-xs font-medium text-text-muted uppercase tracking-wide">Your saved places</p>
          {savedPlaces.map((place) => (
            <button
              key={place.id}
              onClick={() => handleSelect(place.id)}
              className={`w-full text-left p-4 rounded-xl border transition-colors ${
                selectedPlaceId === place.id
                  ? "border-forest bg-forest/5"
                  : "border-background-stone bg-background hover:bg-background-mist"
              }`}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <MapPin className={`h-4 w-4 ${selectedPlaceId === place.id ? "text-forest" : "text-text-muted"}`} />
                  <span className="text-sm font-medium text-text-charcoal">{place.name}</span>
                </div>
                <div className="flex items-center gap-2">
                  {place.is_primary && (
                    <span className="text-xs text-status-normal font-medium">Current</span>
                  )}
                  <span className="text-xs text-text-muted capitalize">{place.place_type}</span>
                </div>
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Search results / nearby */}
      {searchQuery && (
        <div className="space-y-2">
          <p className="text-xs font-medium text-text-muted uppercase tracking-wide">
            {searching ? "Searching..." : `${filteredNearby.length} place${filteredNearby.length !== 1 ? "s" : ""} found`}
          </p>
          {filteredNearby.map((place) => {
            const isSaved = savedPlaces.some((p) => p.id === place.id);
            return (
              <button
                key={place.id}
                onClick={() => handleSelect(place.id)}
                className="w-full text-left p-4 rounded-xl border border-background-stone bg-background hover:bg-background-mist transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <MapPin className="h-4 w-4 text-text-muted" />
                    <div>
                      <span className="text-sm font-medium text-text-charcoal">{place.name}</span>
                      <span className="text-xs text-text-muted ml-2 capitalize">{place.place_type}</span>
                    </div>
                  </div>
                  <span className="text-xs text-text-muted">
                    {place.distance_m ? `${(place.distance_m / 1000).toFixed(1)} km` : ""}
                  </span>
                </div>
              </button>
            );
          })}
          {filteredNearby.length === 0 && !searching && (
            <p className="text-sm text-text-muted text-center py-4">No places match your search.</p>
          )}
        </div>
      )}

      {/* Actions */}
      <div className="space-y-3 pt-2">
        <button
          onClick={handleSkip}
          className="w-full py-3 rounded-xl border border-background-stone bg-background text-text-charcoal font-medium hover:bg-background-mist transition-colors"
        >
          Continue without selecting
        </button>
        <p className="text-center text-metadata text-text-muted">
          You can always change your place later from the profile.
        </p>
      </div>

      {saving && (
        <div className="fixed inset-0 bg-background/80 flex items-center justify-center z-50">
          <div className="bg-background rounded-xl border border-background-stone p-6 flex items-center gap-3 shadow-lg">
            <Loader2 className="h-5 w-5 text-forest animate-spin" />
            <span className="text-sm text-text-charcoal font-medium">Setting your place...</span>
          </div>
        </div>
      )}
    </div>
  );
}
