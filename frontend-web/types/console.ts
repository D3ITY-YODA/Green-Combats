// types/console.ts

export interface ConsoleOrganization {
  name: string;
}

export interface ConsolePriorityItem {
  id: string;
  title: string;
  message: string;
  priority: "normal" | "important" | "urgent";
  href: string;
}

export interface ConsoleOverview {
  organization: ConsoleOrganization;
  important_updates: number;
  community_reports: number;
  delayed_sources: number;
  pending_reviews: number;
  priority_items: ConsolePriorityItem[];
}

export interface ConsoleActivityItem {
  id: string;
  type: "update_created" | "update_approved" | "update_rejected" | "report_received";
  description: string;
  user_name: string;
  timestamp: string;
}

export interface ConsoleFilters {
  status?: string[];
  type?: string[];
  place_ids?: string[];
  date_from?: string;
  date_to?: string;
  page?: number;
  page_size?: number;
}
