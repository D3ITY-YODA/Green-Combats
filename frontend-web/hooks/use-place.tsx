// hooks/use-place.tsx

"use client";

import { createContext, useContext, useState, useCallback, ReactNode } from "react";
import type { UserPlace } from "@/types/places";

interface PlaceContextType {
  currentPlace: UserPlace | null;
  places: UserPlace[];
  setCurrentPlace: (place: UserPlace | null) => void;
  setPlaces: (places: UserPlace[]) => void;
  addPlace: (place: UserPlace) => void;
  removePlace: (placeId: string) => void;
}

const PlaceContext = createContext<PlaceContextType | undefined>(undefined);

interface PlaceProviderProps {
  children: ReactNode;
  initialPlaces?: UserPlace[];
}

export function PlaceProvider({ children, initialPlaces = [] }: PlaceProviderProps) {
  const [places, setPlaces] = useState<UserPlace[]>(initialPlaces);
  const [currentPlace, setCurrentPlace] = useState<UserPlace | null>(
    initialPlaces.find((p) => p.is_primary) || initialPlaces[0] || null
  );

  const setPlacesList = useCallback((newPlaces: UserPlace[]) => {
    setPlaces(newPlaces);
    const primary = newPlaces.find((p) => p.is_primary);
    if (primary) {
      setCurrentPlace(primary);
    } else if (newPlaces.length > 0) {
      setCurrentPlace(newPlaces[0]);
    } else {
      setCurrentPlace(null);
    }
  }, []);

  const addPlace = useCallback((place: UserPlace) => {
    setPlaces((prev) => {
      const exists = prev.some((p) => p.id === place.id);
      if (exists) return prev;
      return [...prev, place];
    });
  }, []);

  const removePlace = useCallback((placeId: string) => {
    setPlaces((prev) => prev.filter((p) => p.id !== placeId));
    setCurrentPlace((prev) => (prev?.id === placeId ? null : prev));
  }, []);

  const value: PlaceContextType = {
    currentPlace,
    places,
    setCurrentPlace,
    setPlaces: setPlacesList,
    addPlace,
    removePlace,
  };

  return <PlaceContext.Provider value={value}>{children}</PlaceContext.Provider>;
}

export function usePlace() {
  const context = useContext(PlaceContext);

  if (context === undefined) {
    throw new Error("usePlace must be used within a PlaceProvider");
  }

  return context;
}