"use client";

import { useState } from 'react';
import type { Place } from '@/types/common';

export function usePlace() {
  const [place, setPlace] = useState<Place | null>(null);

  return { 
    place, 
    setPlace,
    placeName: place?.name ?? 'Select a place'
  };
}
