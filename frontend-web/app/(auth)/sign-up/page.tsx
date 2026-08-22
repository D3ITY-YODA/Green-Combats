"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Loader2 } from "lucide-react";
import { register } from "@/lib/api/auth";

export default function SignUpPage() {
  const router = useRouter();
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await register({
        display_name: displayName,
        email: email || undefined,
        password,
        language: "en",
      });
      router.push("/choose-place");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold text-forest-deep">Create your account</h1>
        <p className="text-text-muted mt-2">Sign up to save your places and preferences.</p>
      </div>
      <form onSubmit={handleRegister} className="space-y-4">
        <input
          type="text"
          placeholder="Full name"
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          required
          className="w-full rounded-lg border border-background-stone bg-background p-3"
        />
        <input
          type="email"
          placeholder="Email address (optional)"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full rounded-lg border border-background-stone bg-background p-3"
        />
        <input
          type="password"
          placeholder="Password (min 8 characters)"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={8}
          className="w-full rounded-lg border border-background-stone bg-background p-3"
        />
        {error && (
          <p className="text-sm text-status-emergency">{error}</p>
        )}
        <button
          type="submit"
          disabled={loading || !displayName || !password}
          className="w-full rounded-xl bg-forest py-3 font-medium text-white hover:bg-forest-deep disabled:opacity-50 flex items-center justify-center gap-2"
        >
          {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : "Create account"}
        </button>
      </form>
      <p className="text-center text-sm text-text-muted">
        Already have an account?{" "}
        <Link href="/sign-in" className="text-forest font-medium">
          Sign in
        </Link>
      </p>
    </div>
  );
}
