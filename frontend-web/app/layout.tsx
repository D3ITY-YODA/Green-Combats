// app/layout.tsx
import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { AppNav } from "@/components/layout/app-nav";

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
        {/* Navigation Sidebar / Bottom Bar */}
        <AppNav />
        
        {/* Main Content Area */}
        <main className="flex-1 pb-20 md:pb-0 md:p-0">
          {children}
        </main>
      </body>
    </html>
  );
}