import Link from "next/link";

export default function SignInPage() {
  return (
    <div className="rounded-2xl border border-stone bg-white p-8 shadow-sm">
      <div className="mb-6 text-center">
        <h1 className="text-2xl font-semibold text-charcoal">Welcome back</h1>
        <p className="mt-2 text-sm text-muted">
          Clear environmental updates for the places that matter to you.
        </p>
      </div>

      <form className="space-y-4">
        <div>
          <label htmlFor="phone-email" className="block text-sm font-medium text-charcoal">
            Phone or email
          </label>
          <input
            id="phone-email"
            type="text"
            className="mt-1 block w-full rounded-xl border border-stone bg-warm-white p-3 text-charcoal placeholder-muted focus:border-forest focus:outline-none focus:ring-1 focus:ring-forest"
            placeholder="Enter phone or email"
          />
        </div>

        <button
          type="submit"
          className="w-full rounded-xl bg-forest px-4 py-3 font-medium text-white transition-colors hover:bg-deep-forest"
        >
          Continue
        </button>
      </form>

      <div className="mt-6 text-center text-sm text-muted">
        <Link href="/sign-up" className="font-medium text-forest hover:underline">
          Create an account
        </Link>
        <span className="mx-2">·</span>
        <Link href="/organization-search" className="font-medium text-forest hover:underline">
          Join an organization
        </Link>
      </div>
    </div>
  );
}
