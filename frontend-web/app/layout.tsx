// app/layout.tsx
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AppNav } from "@/components/layout/app-nav";
import { OfflineBanner } from "@/components/layout/offline-banner";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });

export const metadata: Metadata = {
  title: "Green Compass",
  description: "A calm, location-specific view of what is happening.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={`${inter.variable} flex flex-col md:flex-row min-h-screen bg-background text-text-charcoal`}>
        {/* Offline Banner (Shows only when disconnected) */}
        <OfflineBanner />
        
        {/* Navigation Sidebar / Bottom Bar */}
        <AppNav />
        
        {/* Main Content Area (Added pt-8 to push content down when offline banner shows) */}
        <main className="flex-1 pb-20 md:pb-0 md:p-0 pt-8 md:pt-0">
          {children}
        </main>
      </body>
    </html>
  );
}
