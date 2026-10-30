import { useState, useEffect, useCallback } from "react";
import type { AuthUser, LoginCredentials } from "../models/Auth";
import { login, logout, getMe } from "../services/authApi";

export interface UseAuthReturn {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  handleLogin: (credentials: LoginCredentials) => Promise<void>;
  handleLogout: () => Promise<void>;
}

/**
 * Manages admin session state.
 * - Calls /auth/me on mount to restore an existing session.
 * - Provides handleLogin and handleLogout that update state reactively.
 * - Never stores credentials or secrets in state.
 */
export function useAuth(): UseAuthReturn {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Restore session on mount
  useEffect(() => {
    let cancelled = false;

    const checkSession = async () => {
      try {
        const currentUser = await getMe();
        if (!cancelled) {
          setUser(currentUser);
        }
      } catch {
        // Unexpected error: treat as unauthenticated, don't surface to UI
        if (!cancelled) {
          setUser(null);
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    void checkSession();

    return () => {
      cancelled = true;
    };
  }, []);

  const handleLogin = useCallback(async (credentials: LoginCredentials) => {
    setError(null);
    try {
      await login(credentials);
      // Fetch user info after successful login
      const currentUser = await getMe();
      setUser(currentUser);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Login failed";
      setError(message);
      throw err; // Re-throw so callers can detect failure (e.g. keep modal open)
    }
  }, []);

  const handleLogout = useCallback(async () => {
    setError(null);
    try {
      await logout();
    } catch {
      // Best-effort logout: clear local state even if the request fails
    } finally {
      setUser(null);
    }
  }, []);

  return {
    user,
    isAuthenticated: user !== null,
    isLoading,
    error,
    handleLogin,
    handleLogout,
  };
}
