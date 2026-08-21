// hooks/place-context.tsx

"use client";

import { createContext, useContext, useState, type ReactNode } from "react";
import type { UserPlace } from "@/types/places";

interface PlaceContextType {
  places: UserPlace[];
  currentPlace: UserPlace | null;
  setCurrentPlace: (place: UserPlace | null) => void;
}

const PlaceContext = createContext<PlaceContextType | undefined>(undefined);

interface PlaceProviderProps {
  children: ReactNode;
  initialPlaces?: UserPlace[];
}

export function PlaceProvider({ children, initialPlaces = [] }: PlaceProviderProps) {
  const [places] = useState<UserPlace[]>(initialPlaces);
  const [currentPlace, setCurrentPlace] = useState<UserPlace | null>(
    initialPlaces.find((p) => p.is_primary) ?? initialPlaces[0] ?? null,
  );

  return (
    <PlaceContext.Provider value={{ places, currentPlace, setCurrentPlace }}>
      {children}
    </PlaceContext.Provider>
  );
}

export function usePlace(): PlaceContextType {
  const context = useContext(PlaceContext);
  if (context === undefined) {
    throw new Error("usePlace must be used within a PlaceProvider");
  }
  return context;
}
