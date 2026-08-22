import Link from "next/link";

export default function MarketingPage() {
  return (
    <section className="flex flex-col items-center justify-center px-6 py-24 text-center max-w-4xl mx-auto">
      <h1 className="text-display font-bold text-forest-deep mb-4">
        Know Your Place.<br />Move with Change.
      </h1>
      <p className="text-body text-text-muted max-w-xl mb-8">
        Clear environmental updates for the places that matter to you.
      </p>
      <div className="flex gap-4">
        <Link href="/sign-up" className="px-6 py-3 rounded-xl bg-forest text-white font-medium hover:bg-forest-deep transition-colors">
          Get the App
        </Link>
        <Link href="/today" className="px-6 py-3 rounded-xl border border-background-stone bg-background text-text-charcoal font-medium hover:bg-background-mist transition-colors">
          Explore Web App
        </Link>
      </div>
    </section>
  );
}
