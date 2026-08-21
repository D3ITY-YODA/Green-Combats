import Link from "next/link";

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-background-stone bg-background px-6 py-4 flex items-center justify-between">
        <Link href="/" className="text-page font-bold text-forest-deep">Green Compass</Link>
        <nav className="flex gap-4">
          <Link href="/sign-in" className="text-body text-text-muted hover:text-forest">Sign In</Link>
          <Link href="/sign-up" className="text-body px-4 py-2 rounded-lg bg-forest text-white hover:bg-forest-deep">Get Started</Link>
        </nav>
      </header>
      <main className="flex-1">
        {children}
      </main>
    </div>
  );
}
