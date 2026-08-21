// hooks/use-place.ts

"use client";

import { createContext, useContext, useEffect, useState, useCallback } from "react";
import type { Place, UserPlace } from "@/types/places"; // Adjust path to match your types directory
import { apiFetch } from "@/lib/api/client"; // Adjust path to match your API client

// --- Types ---

interface PlaceContextType {
  currentPlace: Place | null;
  places: UserPlace[];
  isLoading: boolean;
  setCurrentPlace: (placeId: string) => void;
  refetchPlaces: () => Promise<void>;
}

// --- Context ---

const PlaceContext = createContext<PlaceContextType | undefined>(undefined);

// --- Provider ---

interface PlaceProviderProps {
  children: React.ReactNode;
  initialPlaces?: UserPlace[]; // Optional: Pass from server component to prevent flash
}

export function PlaceProvider({ children, initialPlaces }: PlaceProviderProps) {
  const [places, setPlaces] = useState<UserPlace[]>(initialPlaces || []);
  const [isLoading, setIsLoading] = useState(!initialPlaces);

  // Determine initial current place (prefer primary, fallback to first)
  const [currentPlace, setCurrentPlaceState] = useState<Place | null>(
    initialPlaces?.find((p) => p.is_primary)?.place || initialPlaces?.[0]?.place || null
  );

  const fetchPlaces = useCallback(async () => {
    try {
      const response = await apiFetch<UserPlace[]>("/api/v1/me/places");
      setPlaces(response.data);

      const primary = response.data.find((p) => p.is_primary);
      const placeToSet = primary?.place || response.data[0]?.place || null;
      setCurrentPlaceState(placeToSet);
    } catch (error) {
      console.error("Failed to fetch places:", error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!initialPlaces) {
      fetchPlaces();
    }
  }, [initialPlaces, fetchPlaces]);

  const setCurrentPlace = (placeId: string) => {
    const selected = places.find((p) => p.place_id === placeId);
    if (selected) {
      setCurrentPlaceState(selected.place);
      
      // Optional: If your backend supports updating the primary place via API:
      // apiFetch(`/api/v1/me/places/${placeId}/primary`, { method: "POST" }).catch(console.error);
    }
  };

  const value: PlaceContextType = {
    currentPlace,
    places,
    isLoading,
    setCurrentPlace,
    refetchPlaces: fetchPlaces,
  };

  return <PlaceContext.Provider value={value}>{children}</PlaceContext.Provider>;
}

// --- Hook ---

export function usePlace() {
  const context = useContext(PlaceContext);

  if (context === undefined) {
    throw new Error("usePlace must be used within a PlaceProvider");
  }

  return context;
}
