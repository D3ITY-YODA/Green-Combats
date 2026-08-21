export default function Loading() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-8 space-y-6">
      <div className="h-8 w-40 animate-pulse rounded bg-background-stone" />
      <div className="h-40 animate-pulse rounded-2xl bg-background-stone" />
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="h-32 animate-pulse rounded-2xl bg-background-stone" />
        <div className="h-32 animate-pulse rounded-2xl bg-background-stone" />
      </div>
    </div>
  );
}
