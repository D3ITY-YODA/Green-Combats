// app/(public)/layout.tsx

import { AuthProvider } from "@/hooks/use-auth";
import { PlaceProvider } from "@/hooks/use-place";
import { serverApiFetch } from "@/lib/api/server-client";
import type { UserPlace } from "@/types/places";

export default async function PublicLayout({ children }: { children: React.ReactNode }) {
  // Optional: Fetch places on the server to prevent loading flash
  let userPlaces: UserPlace[] | undefined;
  try {
    const res = await serverApiFetch<UserPlace[]>("/api/v1/me/places");
    userPlaces = res.data;
  } catch {
    // User is not logged in or has no places
  }

  return (
    <AuthProvider>
      <PlaceProvider initialPlaces={userPlaces}>
        {/* ... rest of your layout (PublicHeader, main, PublicNavigation) ... */}
        {children}
      </PlaceProvider>
    </AuthProvider>
  );
}
