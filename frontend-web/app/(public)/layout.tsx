import { AppNav } from "@/components/layout/app-nav";
import { OfflineBanner } from "@/components/layout/offline-banner";

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col md:flex-row min-h-screen">
      <OfflineBanner />
      <AppNav />
      <main className="flex-1 pb-20 md:pb-0 md:p-0 pt-8 md:pt-0">
        {children}
      </main>
    </div>
  );
}
