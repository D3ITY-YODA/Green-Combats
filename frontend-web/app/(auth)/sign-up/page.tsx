import Link from "next/link";
export default function SignUpPage() {
  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold text-forest-deep">Create your account</h1>
        <p className="text-text-muted mt-2">Sign up to save your places and preferences.</p>
      </div>
      <form className="space-y-4">
        <input type="text" placeholder="Full name" className="w-full rounded-lg border border-background-stone bg-background p-3" />
        <input type="email" placeholder="Email address" className="w-full rounded-lg border border-background-stone bg-background p-3" />
        <input type="password" placeholder="Password" className="w-full rounded-lg border border-background-stone bg-background p-3" />
        <button type="button" className="w-full rounded-xl bg-forest py-3 font-medium text-white hover:bg-forest-deep">Create account</button>
      </form>
      <p className="text-center text-sm text-text-muted">Already have an account? <Link href="/sign-in" className="text-forest font-medium">Sign in</Link></p>
    </div>
  );
}
