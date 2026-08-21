// components/states/empty-state.tsx

import type { ReactNode } from "react";
import { Inbox } from "lucide-react";
import type { LucideIcon } from "lucide-react";

interface EmptyStateProps {
  title: string;
  message: string;
  action?: ReactNode;
  icon?: LucideIcon;
}

export function EmptyState({
  title,
  message,
  action,
  icon: Icon = Inbox,
}: EmptyStateProps) {
  return (
    <section className="rounded-2xl border border-stone bg-white p-8 text-center">
      {/* Icon */}
      <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-stone/20">
        <Icon 
          className="h-6 w-6 text-muted" 
          aria-hidden="true" 
        />
      </div>

      {/* Text Content */}
      <h2 className="mt-4 text-lg font-semibold text-charcoal">
        {title}
      </h2>
      <p className="mt-2 text-sm text-muted">
        {message}
      </p>

      {/* Optional Action (e.g., a button or link) */}
      {action && (
        <div className="mt-5 flex justify-center">
          {action}
        </div>
      )}
    </section>
  );
}
