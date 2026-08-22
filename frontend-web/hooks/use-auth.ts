"use client";

import { createContext, useContext } from 'react';

// Minimal auth context stub to satisfy TypeScript and unblock the build
export const AuthContext = createContext<{ user: any; isLoading: boolean }>({ 
  user: null, 
  isLoading: false 
});

export function useAuth() {
  return useContext(AuthContext);
}
