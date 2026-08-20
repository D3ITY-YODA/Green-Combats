export default function Loading() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <div className="h-8 w-40 animate-pulse rounded bg-stone" />
      <div className="mt-8 space-y-4">
        <div className="h-40 animate-pulse rounded-2xl bg-stone" />
        <div className="h-40 animate-pulse rounded-2xl bg-stone" />
      </div>
    </div>
  );
}
