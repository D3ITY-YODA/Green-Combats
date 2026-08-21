import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-warm-white p-4">
      <h2 className="text-2xl font-semibold text-charcoal">Page not found</h2>
      <p className="mt-2 text-muted">The page you are looking for does not exist.</p>
      <Link
        href="/today"
        className="mt-6 rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-deep-forest transition-colors"
      >
        Return to Today
      </Link>
    </div>
  );
}
