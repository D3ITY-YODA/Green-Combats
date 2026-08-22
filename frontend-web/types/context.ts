// types/context.ts

import type { UUID } from "./common";

/**
 * Topic keys for categorizing updates and reports
 * Based on the Green Compass domain model (Section 5.6 of blueprint)
 */
export type TopicKey =
  | "water.availability"
  | "water.quality"
  | "water.infrastructure"
  | "weather.forecast"
  | "weather.alerts"
  | "weather.historical"
  | "agriculture.crop_health"
  | "agriculture.pest_disease"
  | "agriculture.market_prices"
  | "environment.air_quality"
  | "environment.biodiversity"
  | "environment.land_use"
  | "community.announcements"
  | "community.events"
  | "community.reports"
  | "infrastructure.roads"
  | "infrastructure.power"
  | "infrastructure.connectivity"
  | "health.facilities"
  | "health.outbreaks"
  | "health.services"
  | "emergency.incidents"
  | "emergency.evacuation"
  | "emergency.relief";

export interface Topic {
  key: TopicKey;
  label: Record<string, string>; // Localized labels
  icon: string;
  color: string;
  category: "water" | "weather" | "agriculture" | "environment" | "community" | "infrastructure" | "health" | "emergency";
}

/**
 * Context filters for API queries
 */
export interface ContextFilters {
  place_id?: UUID;
  topic_keys?: TopicKey[];
  types?: string[];
  status?: string[];
  valid_from?: string;
  valid_until?: string;
  locale?: string;
}

export interface TopicSection {
  topic_key: TopicKey;
  label: string;
  icon: string;
  color: string;
  updates: string[]; // Update IDs
  available?: boolean;
  title?: string;
  description?: string;
}