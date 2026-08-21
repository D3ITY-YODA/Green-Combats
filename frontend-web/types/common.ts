export type UUID = string;
export type Locale = "en" | "sw" | "fr" | string;

export type DataStatus = "current" | "delayed" | "stale" | "unavailable" | "not_configured";
export type RiskState = "normal" | "watch" | "warning" | "emergency" | "stale" | "unavailable";

export interface Place {
  id: UUID;
  name: string;
  type: string;
  country_code: string;
}

export interface TopicSection {
  key: string;
  title: string;
  description: string;
  available: boolean;
  data_status?: DataStatus;
  href?: string;
}
