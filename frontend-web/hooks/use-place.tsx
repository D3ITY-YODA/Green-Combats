// hooks/use-place.tsx (layout wrapper — re-exports for backward compatibility)

import { AuthProvider } from "@/hooks/use-auth";
import { PlaceProvider } from "@/hooks/place-context";
import { serverApiFetch } from "@/lib/api/server-client";
import type { UserPlace } from "@/types/places";

export { PlaceProvider, usePlace } from "@/hooks/place-context";

export default async function PublicLayout({ children }: { children: React.ReactNode }) {
  // Fetch places on the server to prevent loading flash
  let userPlaces: UserPlace[] | undefined;
  try {
    const res = await serverApiFetch<{ places: UserPlace[] }>("/api/v1/me/places");
    userPlaces = res.data.places;
  } catch {
    // User is not logged in or has no places
  }

  return (
    <AuthProvider>
      <PlaceProvider initialPlaces={userPlaces}>
        {children}
      </PlaceProvider>
    </AuthProvider>
  );
}
