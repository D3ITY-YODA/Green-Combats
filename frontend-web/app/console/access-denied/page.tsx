import Link from "next/link";
export default function AccessDeniedPage() {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
      <h1 className="text-2xl font-bold text-text-charcoal">Access not available</h1>
      <p className="mt-2 text-text-muted">You do not have permission to view this organization's information.</p>
      <Link href="/console" className="mt-6 rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-forest-deep">Return to Overview</Link>
    </div>
  );
}
