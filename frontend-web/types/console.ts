// types/console.ts

/**
 * Institutional console overview returned by the backend.
 */
export interface ConsoleOverview {
  organization: {
    id: string;
    name: string;
  };
  important_updates: number;
  community_reports: number;
  delayed_sources: number;
  pending_reviews: number;
  priority_items: ConsolePriorityItem[];
}

/**
 * A priority action item in the console dashboard.
 */
export interface ConsolePriorityItem {
  id: string;
  title: string;
  message: string;
  priority: "normal" | "important" | "urgent";
  href: string;
}
