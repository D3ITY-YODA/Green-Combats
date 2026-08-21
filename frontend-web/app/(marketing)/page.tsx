// app/(marketing)/page.tsx

import Link from "next/link";
import { Globe } from "lucide-react";

export default function MarketingPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-emerald-50 to-white px-4 py-12">
      {/* Main Content */}
      <main className="flex w-full max-w-md flex-col items-center text-center">
        {/* Logo */}
        <div className="mb-8">
          <div className="flex h-24 w-24 items-center justify-center rounded-full bg-emerald-100 ring-8 ring-emerald-50">
            <Globe 
              className="h-12 w-12 text-emerald-700" 
              aria-hidden="true"
            />
          </div>
        </div>

        {/* Brand Name */}
        <h1 className="text-3xl font-bold tracking-tight text-emerald-900">
          GREEN COMPASS
        </h1>

        {/* Tagline */}
        <p className="mt-3 text-xl font-medium text-emerald-800">
          Know Your Place.<br />Move with Change.
        </p>

        {/* Description */}
        <p className="mt-4 max-w-sm text-base leading-relaxed text-stone-600">
          Clear environmental updates for the places that matter to you.
        </p>

        {/* CTA Button */}
        <div className="mt-10 w-full">
          <Link
            href="/choose-language"
            className="flex w-full items-center justify-center rounded-2xl bg-emerald-700 px-8 py-4 text-lg font-semibold text-white shadow-lg shadow-emerald-200 transition-all duration-200 hover:bg-emerald-800 hover:shadow-xl hover:shadow-emerald-200 focus:outline-none focus:ring-4 focus:ring-emerald-300 active:scale-95"
          >
            Get started
          </Link>
        </div>

        {/* Language Selector Link */}
        <div className="mt-6">
          <Link
            href="/choose-language"
            className="inline-flex items-center gap-2 text-sm font-medium text-emerald-700 transition-colors hover:text-emerald-800 focus:outline-none focus:ring-2 focus:ring-emerald-300 focus:ring-offset-2 rounded-md px-2 py-1"
          >
            <Globe className="h-4 w-4" aria-hidden="true" />
            Choose language
          </Link>
        </div>
      </main>

      {/* Footer */}
      <footer className="mt-auto pt-12 text-center">
        <p className="text-sm text-stone-500">
          Built with care for people and the environment
        </p>
      </footer>
    </div>
  );
}
