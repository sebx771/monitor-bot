const memoryFallback = new Map<string, string>();

/**
 * Safe storage utility that wraps localStorage in a try/catch block.
 * Gracefully falls back to an in-memory Map if localStorage is restricted,
 * sandboxed, or throws a QuotaExceededError / SecurityError.
 */
export const safeStorage = {
  getItem(key: string): string | null {
    try {
      if (typeof window !== "undefined") {
        return window.localStorage.getItem(key);
      }
    } catch {
      // In private browsing or restricted iframe
    }
    return memoryFallback.get(key) ?? null;
  },

  setItem(key: string, value: string): void {
    try {
      if (typeof window !== "undefined") {
        window.localStorage.setItem(key, value);
        return;
      }
    } catch {
      // Storage failed or quota exceeded
    }
    memoryFallback.set(key, value);
  },

  removeItem(key: string): void {
    try {
      if (typeof window !== "undefined") {
        window.localStorage.removeItem(key);
        return;
      }
    } catch {
      // Ignore exception
    }
    memoryFallback.delete(key);
  },
};
