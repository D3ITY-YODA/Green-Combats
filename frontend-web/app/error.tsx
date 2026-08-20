"use client";

import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-warm-white p-4">
      <h2 className="text-2xl font-semibold text-charcoal">Something went wrong</h2>
      <p className="mt-2 text-muted">We could not load this information right now.</p>
      <button
        className="mt-6 rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-deep-forest transition-colors"
        onClick={() => reset()}
      >
        Try again
      </button>
    </div>
  );
}
