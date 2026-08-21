import { ConsoleSidebar } from "@/components/navigation/console-sidebar";

export default function ConsoleLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen bg-background-mist">
      <ConsoleSidebar />
      <div className="flex-1 flex flex-col min-w-0">
        {/* Console Top Header */}
        <header className="bg-background border-b border-background-stone px-6 py-4 flex items-center justify-between sticky top-0 z-10">
          <div>
            <h2 className="text-lg font-semibold text-text-charcoal">Lower Valley Water Authority</h2>
            <p className="text-xs text-text-muted">Organization Administrator</p>
          </div>
          <div className="h-8 w-8 rounded-full bg-forest text-white flex items-center justify-center text-sm font-medium">
            LV
          </div>
        </header>
        <main className="flex-1 p-6 md:p-8 overflow-y-auto">
          {children}
        </main>
      </div>
    </div>
  );
}
