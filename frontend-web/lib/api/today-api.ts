// lib/api/today-api.ts

export interface TodayData {
  location: string;
  lastUpdated: string;
  hasImportantUpdates: boolean;
  localOutlook: string;
  waterOutlook: string;
}

// Simulate network delay
const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export async function getTodayData(): Promise<TodayData> {
  await delay(300); // Simulate API latency
  
  // MOCK DATA: Toggle hasImportantUpdates to test the UI states
  return {
    location: "Lower Valley",
    lastUpdated: "10:00",
    hasImportantUpdates: false, 
    localOutlook: "Conditions for the coming days are stable with mild temperatures.",
    waterOutlook: "Water availability may decline slightly over the next two weeks.",
  };
}
