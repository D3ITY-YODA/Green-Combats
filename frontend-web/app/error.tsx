"use client";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="mx-auto max-w-3xl px-4 py-16 text-center">
      <h2 className="text-2xl font-semibold text-text-charcoal">Something went wrong</h2>
      <p className="mt-2 text-text-muted">We could not load this information right now.</p>
      <button
        onClick={() => reset()}
        className="mt-6 rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-forest-deep transition-colors"
      >
        Try again
      </button>
    </div>
  );
}
