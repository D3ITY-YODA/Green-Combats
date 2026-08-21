// types/context.ts

/**
 * Topic keys matching the backend's topic_key enum.
 */
export type TopicKey =
  | "local_outlook"
  | "water_outlook"
  | "seasonal_information"
  | "land_and_ecosystems"
  | "food_and_agriculture"
  | "community_updates";

/**
 * A topic section used in the Explore grid view.
 */
export interface TopicSection {
  key: TopicKey;
  title: string;
  description: string;
  available: boolean;
}
