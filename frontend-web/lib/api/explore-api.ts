// lib/api/explore-api.ts

export interface ExploreSection {
  id: string;
  title: string;
  description: string;
  icon: "local" | "water" | "seasonal" | "community";
}

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export async function getExploreSections(): Promise<ExploreSection[]> {
  await delay(400);

  // MOCK DATA: Only returning relevant sections as per PDF rules
  return [
    {
      id: "1",
      title: "Local outlook",
      description: "Conditions for the coming days are stable with mild temperatures.",
      icon: "local",
    },
    {
      id: "2",
      title: "Water outlook",
      description: "Water availability may decline slightly over the next two weeks.",
      icon: "water",
    },
    {
      id: "3",
      title: "Seasonal information",
      description: "Planting season is approaching. Soil moisture levels are adequate.",
      icon: "seasonal",
    },
  ];
}
