import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Green Compass",
  description: "Clear environmental updates for the places that matter to you.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased bg-warm-white text-charcoal">
        {children}
      </body>
    </html>
  );
}
