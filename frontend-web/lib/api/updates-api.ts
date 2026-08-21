// lib/api/updates-api.ts

export type StatusLevel = "normal" | "watch" | "important" | "emergency" | "delayed" | "unavailable";

export interface UpdateItem {
  id: string;
  status: StatusLevel;
  title: string;
  description: string;
  location: string;
  updated: string;
}

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export async function getUpdates(): Promise<UpdateItem[]> {
  await delay(400);

  // MOCK DATA: Showing a mix of statuses to test the UI
  return [
    {
      id: "1",
      status: "emergency",
      title: "Flood warning",
      description: "Heavy rainfall may cause flooding within 24 hours in low-lying areas.",
      location: "Lower Valley",
      updated: "10:00",
    },
    {
      id: "2",
      status: "watch",
      title: "Water outlook",
      description: "Water availability may decline slightly over the next two weeks.",
      location: "Lower Valley",
      updated: "Today",
    },
    {
      id: "3",
      status: "normal",
      title: "Seasonal update",
      description: "Planting season is approaching. Soil moisture levels are currently adequate.",
      location: "Lower Valley",
      updated: "Yesterday",
    },
  ];
}
