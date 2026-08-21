import Link from "next/link";

export default function NotFound() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-16 text-center">
      <h2 className="text-2xl font-semibold text-text-charcoal">Page not found</h2>
      <p className="mt-2 text-text-muted">The page you are looking for does not exist.</p>
      <Link
        href="/today"
        className="mt-6 inline-block rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-forest-deep transition-colors"
      >
        Return to Today
      </Link>
    </div>
  );
}
