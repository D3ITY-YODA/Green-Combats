"use client";

import { createContext, useContext, useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import type { User } from "@/types/users";
import { fetchAPI } from "@/lib/api/client";

// --- Types ---

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  refreshUser: () => Promise<void>;
  signOut: () => Promise<void>;
}

// --- Context ---

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// --- Provider ---

interface AuthProviderProps {
  children: React.ReactNode;
  initialUser?: User | null; // Optional: Pass user from server component to prevent flash
}

export function AuthProvider({ children, initialUser = null }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(initialUser);
  const [isLoading, setIsLoading] = useState(!initialUser);
  const router = useRouter();

  const fetchUser = useCallback(async () => {
    try {
      const data = await fetchAPI<User>("/v1/auth/me");
      setUser(data);
    } catch {
      // If the API returns 401 or fails, the user is not authenticated
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    // Only fetch if we didn't receive an initial user from the server
    if (!initialUser) {
      fetchUser();
    }
  }, [initialUser, fetchUser]);

  const signOut = async () => {
    try {
      setUser(null);
      router.push("/sign-in");
      router.refresh(); // Clear server component cache
    } catch (error) {
      console.error("Failed to sign out:", error);
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated: !!user && user.status === "active",
    refreshUser: fetchUser,
    signOut,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// --- Hook ---

export function useAuth() {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
}
