// middleware.ts

import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Paths that do not require authentication
const publicPaths = [
  "/",
  "/sign-in",
  "/sign-up",
  "/verify",
  "/privacy",
  "/terms",
];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  
  // Check for the session cookie (adjust the cookie name if your backend uses a different one)
  const hasSession = request.cookies.has("gc_session");

  // Check if the current path is a public path
  const isPublicPath = publicPaths.some(
    (path) => pathname === path || pathname.startsWith(`${path}/`),
  );

  // Redirect to sign-in if the path is not public and the user has no session
  if (!isPublicPath && !hasSession) {
    const url = request.nextUrl.clone();
    url.pathname = "/sign-in";
    url.searchParams.set("next", pathname);
    return NextResponse.redirect(url);
  }

  // Explicitly protect the console routes (redundant but safe, as console is not in publicPaths)
  if (pathname.startsWith("/console") && !hasSession) {
    const url = request.nextUrl.clone();
    url.pathname = "/sign-in";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

// Configure the middleware to run on specific paths
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
};
