// app/console/layout.tsx

import { ConsoleSidebar } from "@/components/navigation/console-sidebar";
import { ConsoleHeader } from "@/components/layout/console-header";

export default function ConsoleLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-mist">
      {/* Persistent left sidebar for organization navigation */}
      <ConsoleSidebar isOpen={true} onClose={() => {}} />
      
      {/* Main content wrapper with left padding on large screens to clear the fixed sidebar */}
      <div className="lg:pl-72">
        {/* Top header with organization context, user info, and actions */}
        <ConsoleHeader />
        
        {/* Page-specific content with responsive padding */}
        <main className="p-4 md:p-8">
          {children}
        </main>
      </div>
    </div>
  );
}
