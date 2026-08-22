"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Sprout, ArrowRight, Loader2 } from "lucide-react";
import { login } from "@/lib/api/auth";

export default function SignInPage() {
  const router = useRouter();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await login({ identifier, password });
      router.push("/today");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleDemoAccess = () => {
    // Bypass auth for demo — go straight to the dashboard
    router.push("/today");
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-6 bg-background">
      <div className="w-full max-w-md space-y-8">
        {/* Logo / Brand */}
        <div className="flex flex-col items-center gap-3">
          <div className="p-4 rounded-full bg-forest/10">
            <Sprout className="h-10 w-10 text-forest" />
          </div>
          <h1 className="text-page font-bold text-forest-deep">Green Compass</h1>
          <p className="text-body text-text-muted max-w-xs mx-auto">
            A calm, location-specific view of what is happening in your area.
          </p>
        </div>

        {/* Login Form */}
        <form onSubmit={handleLogin} className="space-y-4 pt-4">
          <input
            type="text"
            placeholder="Email or phone number"
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
            className="w-full rounded-xl border border-background-stone bg-background p-4 text-text-charcoal focus:border-forest focus:ring-1 focus:ring-forest"
            autoComplete="username"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full rounded-xl border border-background-stone bg-background p-4 text-text-charcoal focus:border-forest focus:ring-1 focus:ring-forest"
            autoComplete="current-password"
          />

          {error && (
            <p className="text-sm text-status-emergency">{error}</p>
          )}

          <button
            type="submit"
            disabled={loading || !identifier || !password}
            className="w-full py-4 rounded-xl bg-forest text-white font-medium flex items-center justify-center gap-2 hover:bg-forest-deep transition-colors shadow-sm disabled:opacity-50"
          >
            {loading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <>
                Sign in
                <ArrowRight className="h-4 w-4" />
              </>
            )}
          </button>
        </form>

        {/* Demo access */}
        <div className="space-y-3 pt-2">
          <button
            onClick={handleDemoAccess}
            className="w-full py-3 rounded-xl border border-background-stone bg-background text-text-charcoal font-medium hover:bg-background-mist transition-colors"
          >
            Continue without account
          </button>
          <p className="text-center text-metadata text-text-muted">
            No complex forms required. You can select your location and preferences quietly inside the app.
          </p>
        </div>
      </div>
    </main>
  );
}
