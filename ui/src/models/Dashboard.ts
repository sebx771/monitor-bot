export type ViewDensity = "grid" | "table";

export type StatusFilter = "all" | "operational" | "alert";

export type RefreshInterval = 0 | 15 | 30 | 60;

export type ToastType = "success" | "error" | "info";

export interface ToastMessage {
  id: string;
  type: ToastType;
  text: string;
  durationMs?: number;
}

export interface DashboardStats {
  total: number;
  operational: number;
  alertCount: number;
}
