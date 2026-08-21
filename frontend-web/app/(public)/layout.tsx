// app/(public)/layout.tsx

import { PublicHeader } from "@/components/layout/public-header";
import { PublicNavigation } from "@/components/navigation/public-navigation";

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-warm-white">
      {/* Top header with place selector, profile link, and notifications */}
      <PublicHeader />
      
      {/* Main content area with bottom padding to prevent content from being hidden behind the fixed bottom navigation */}
      <main className="pb-24">
        {children}
      </main>
      
      {/* Fixed bottom navigation for primary public routes */}
      <PublicNavigation />
    </div>
  );
}
