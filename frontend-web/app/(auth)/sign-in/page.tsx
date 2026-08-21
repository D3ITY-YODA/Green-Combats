// app/sign-in/page.tsx
"use client";

import { useRouter } from "next/navigation";
import { Sprout, ArrowRight } from "lucide-react";

export default function SignInPage() {
  const router = useRouter();

  const handleContinue = () => {
    // Bypass complex auth and go straight to the core experience
    // In a real app, this might set a lightweight guest session
    router.push("/today");
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-6 bg-background">
      <div className="w-full max-w-md text-center space-y-8">
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

        {/* Frictionless Entry (PDF: "should not force users to identify themselves through complex forms") */}
        <div className="space-y-4 pt-4">
          <button
            onClick={handleContinue}
            className="w-full py-4 rounded-xl bg-forest text-white font-medium flex items-center justify-center gap-2 hover:bg-forest-deep transition-colors shadow-sm"
          >
            Continue to Dashboard
            <ArrowRight className="h-4 w-4" />
          </button>
          
          <p className="text-metadata text-text-muted">
            No complex forms required. You can select your location and preferences quietly inside the app.
          </p>
        </div>
      </div>
    </main>
  );
}
