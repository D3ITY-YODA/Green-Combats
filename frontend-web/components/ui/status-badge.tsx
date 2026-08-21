// components/ui/status-badge.tsx
import { cn } from "@/lib/utils";

type StatusLevel = "normal" | "watch" | "important" | "emergency" | "delayed" | "unavailable";

interface StatusBadgeProps {
  status: StatusLevel;
  label: string;
  className?: string;
}

const STATUS_STYLES: Record<StatusLevel, { bg: string; text: string; border: string }> = {
  normal: { bg: "bg-status-normal/10", text: "text-status-normal", border: "border-status-normal/20" },
  watch: { bg: "bg-status-watch/10", text: "text-status-watch", border: "border-status-watch/20" },
  important: { bg: "bg-status-important/10", text: "text-status-important", border: "border-status-important/20" },
  emergency: { bg: "bg-status-emergency/10", text: "text-status-emergency", border: "border-status-emergency/20" },
  delayed: { bg: "bg-status-delayed/10", text: "text-status-delayed", border: "border-status-delayed/20" },
  unavailable: { bg: "bg-status-unavailable/10", text: "text-status-unavailable", border: "border-status-unavailable/20" },
};

export function StatusBadge({ status, label, className }: StatusBadgeProps) {
  const styles = STATUS_STYLES[status];
  return (
    <span className={cn(
      "inline-flex items-center rounded-full border px-2.5 py-0.5 text-metadata font-medium",
      styles.bg, styles.text, styles.border, className
    )}>
      {label}
    </span>
  );
}
