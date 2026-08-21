"use client";

import { useState } from 'react';
import type { TodayResponse } from '@/types/updates';

export function useToday() {
  const [data, setData] = useState<TodayResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);

  return { data, loading, error, setData, setLoading, setError };
}
