import type { AuthUser, LoginCredentials } from "../models/Auth";
import { API_BASE_URL } from "../core/config";

/**
 * POST /api/v1/auth/login
 * Sends credentials; the browser stores the HttpOnly session cookie automatically.
 * Throws Error("Invalid credentials") on 401.
 * Throws generic Error on other failures.
 */
export async function login(credentials: LoginCredentials): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    body: JSON.stringify(credentials),
  });

  if (response.status === 401) {
    throw new Error("Invalid credentials");
  }

  if (!response.ok) {
    throw new Error(`Login failed (${response.status})`);
  }
}

/**
 * POST /api/v1/auth/logout
 * Asks the server to invalidate the session and expire the cookie.
 * Throws Error on failure.
 */
export async function logout(): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
    headers: { Accept: "application/json" },
  });

  if (!response.ok) {
    throw new Error(`Logout failed (${response.status})`);
  }
}

/**
 * GET /api/v1/auth/me
 * Returns the current user if a valid session cookie exists.
 * Returns null (without throwing) on 401 — normal "not logged in" state.
 * Throws Error for unexpected server errors.
 */
export async function getMe(): Promise<AuthUser | null> {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/me`, {
      method: "GET",
      credentials: "include",
      headers: { Accept: "application/json" },
    });

    if (!response.ok) {
      return null;
    }

    const data = (await response.json()) as { authenticated?: boolean; username?: string };
    if (!data.authenticated || !data.username) {
      return null;
    }

    return { username: data.username };
  } catch {
    return null;
  }
}
