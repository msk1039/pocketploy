"use client";

import React, { createContext, useContext, useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import {
  login as apiLogin,
  signup as apiSignup,
  logout as apiLogout,
  getCurrentUser,
  refreshAccessToken,
  clearTokens,
  getAccessToken,
  ApiError,
} from "@/lib/api";
import { User, AuthContextType } from "@/types/auth";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const refreshTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const router = useRouter();

  // Schedule token refresh before expiry
  const scheduleTokenRefresh = (expiresAt: string) => {
    // Clear existing timeout
    if (refreshTimeoutRef.current) {
      clearTimeout(refreshTimeoutRef.current);
    }

    // Calculate time until refresh (1 minute before expiry)
    const expiryTime = new Date(expiresAt).getTime();
    const currentTime = Date.now();
    const timeUntilRefresh = expiryTime - currentTime - 60000; // 1 minute before expiry

    if (timeUntilRefresh > 0) {
      refreshTimeoutRef.current = setTimeout(async () => {
        try {
          await refreshAccessToken();
          // Schedule next refresh (assume 15 minutes expiry)
          const newExpiryTime = new Date(Date.now() + 15 * 60 * 1000).toISOString();
          scheduleTokenRefresh(newExpiryTime);
        } catch (error) {
          console.error("Failed to refresh token:", error);
          // If refresh fails, clear session
          clearTokens();
          setUser(null);
          router.push("/login");
        }
      }, timeUntilRefresh);
    }
  };

  // Load user on mount
  useEffect(() => {
    async function loadUser() {
      const refreshToken = typeof window !== "undefined" 
        ? localStorage.getItem("refresh_token") 
        : null;

      if (!refreshToken) {
        setIsLoading(false);
        return;
      }

      try {
        console.log("Loading user session...");
        
        // Try to refresh the access token first
        const newAccessToken = await refreshAccessToken();
        console.log("Access token refreshed successfully");

        // Now fetch current user with the new access token
        const response = await getCurrentUser();
        console.log("User data fetched:", response.data.user.username);
        
        setUser(response.data.user);

        // Schedule token refresh
        const expiryTime = new Date(Date.now() + 15 * 60 * 1000).toISOString();
        scheduleTokenRefresh(expiryTime);
      } catch (error) {
        console.error("Failed to load user:", error);
        clearTokens();
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    }

    loadUser();

    // Cleanup timeout on unmount
    return () => {
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
      }
    };
  }, []); // Only run on mount

  const login = async (email: string, password: string) => {
    try {
      const response = await apiLogin({ email, password });
      setUser(response.data.user);
      scheduleTokenRefresh(response.data.expires_at);
      router.push("/dashboard");
    } catch (error) {
      if (error instanceof ApiError) {
        throw new Error(error.message);
      }
      throw new Error("Login failed. Please try again.");
    }
  };

  const signup = async (username: string, email: string, password: string) => {
    try {
      const response = await apiSignup({ username, email, password });
      setUser(response.data.user);
      scheduleTokenRefresh(response.data.expires_at);
      router.push("/dashboard");
    } catch (error) {
      if (error instanceof ApiError) {
        if (error.details) {
          // Format validation errors
          const messages = Object.values(error.details).join(", ");
          throw new Error(messages);
        }
        throw new Error(error.message);
      }
      throw new Error("Signup failed. Please try again.");
    }
  };

  const logout = async () => {
    try {
      await apiLogout();
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      setUser(null);
      clearTokens();
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
        refreshTimeoutRef.current = null;
      }
      router.push("/login");
    }
  };

  const refreshUser = async () => {
    try {
      const response = await getCurrentUser();
      setUser(response.data.user);
    } catch (error) {
      console.error("Failed to refresh user:", error);
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    signup,
    logout,
    refreshUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
