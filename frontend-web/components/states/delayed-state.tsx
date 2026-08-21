// components/states/delayed-state.tsx

"use client";

import { Clock } from "lucide-react";

interface DelayedStateProps {
  title?: string;
  message?: string;
  onRetry?: () => void;
  isRetrying?: boolean;
}

export function DelayedState({
  title = "Information delayed",
  message = "The latest update for this place is not available yet.",
  onRetry,
  isRetrying = false,
}: DelayedStateProps) {
  return (
    <section 
      className="rounded-2xl border border-stone bg-white p-8 text-center"
      role="status"
      aria-live="polite"
    >
      {/* Icon */}
      <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-stone/20">
        <Clock 
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

      {/* Action Button */}
      {onRetry && (
        <div className="mt-5">
          <button
            type="button"
            onClick={onRetry}
            disabled={isRetrying}
            className="inline-flex items-center justify-center rounded-xl bg-forest px-5 py-2.5 text-sm font-medium text-white transition hover:bg-forest/90 focus:outline-none focus:ring-2 focus:ring-forest focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {isRetrying ? "Checking..." : "Try again"}
          </button>
        </div>
      )}
    </section>
  );
}
