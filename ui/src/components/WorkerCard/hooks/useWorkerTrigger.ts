import { useState, useRef, useEffect, useCallback } from "react";
import { WorkerStatus, type Worker } from "../../../models/Worker";
import { triggerWorker } from "../../../services/workerApi";

export interface TriggerNotice {
  type: "success" | "conflict" | "error";
  text: string;
}

interface UseWorkerTriggerProps {
  worker: Worker;
  onCustomTrigger?: (workerId: string) => Promise<void> | void;
  activeConfirmationWorkerId?: string | null;
  onRequestConfirmation?: (workerId: string) => void;
  onClearConfirmation?: () => void;
}

export function useWorkerTrigger({
  worker,
  onCustomTrigger,
  activeConfirmationWorkerId,
  onRequestConfirmation,
  onClearConfirmation,
}: UseWorkerTriggerProps) {
  const [isExecuting, setIsExecuting] = useState(false);
  const [notice, setNotice] = useState<TriggerNotice | null>(null);
  const [localConfirming, setLocalConfirming] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const isOperational = worker.status === WorkerStatus.OPERATIONAL;
  const isConfirming = activeConfirmationWorkerId !== undefined
    ? activeConfirmationWorkerId === worker.id
    : localConfirming;

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  useEffect(() => {
    return () => {
      clearTimer();
    };
  }, [clearTimer]);

  const cancelConfirmation = useCallback(() => {
    clearTimer();
    setLocalConfirming(false);
    onClearConfirmation?.();
  }, [clearTimer, onClearConfirmation]);

  const executeTrigger = useCallback(async () => {
    cancelConfirmation();
    if (!isOperational || isExecuting) return;

    setIsExecuting(true);
    setNotice(null);

    try {
      if (onCustomTrigger) {
        await onCustomTrigger(worker.id);
      } else {
        await triggerWorker(worker.id);
      }
      setNotice({ type: "success", text: "Cycle queued!" });
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Failed to trigger worker";
      const isConflict = message.toLowerCase().includes("ocupado") || message.toLowerCase().includes("busy");
      setNotice({ type: isConflict ? "conflict" : "error", text: message });
    } finally {
      setIsExecuting(false);
      timerRef.current = setTimeout(() => {
        setNotice(null);
      }, 3000);
    }
  }, [cancelConfirmation, isOperational, isExecuting, onCustomTrigger, worker.id]);

  const handleTriggerClick = useCallback(() => {
    if (isExecuting || !isOperational) return;

    if (isConfirming) {
      void executeTrigger();
    } else {
      setLocalConfirming(true);
      onRequestConfirmation?.(worker.id);
      clearTimer();
      timerRef.current = setTimeout(() => {
        cancelConfirmation();
      }, 3000);
    }
  }, [isExecuting, isOperational, isConfirming, executeTrigger, onRequestConfirmation, worker.id, clearTimer, cancelConfirmation]);

  const tooltipText = isExecuting
    ? "Sending request..."
    : isConfirming
    ? "Click again to confirm immediate trigger (3s)"
    : !isOperational
    ? `Trigger unavailable: worker in "${worker.status}" state`
    : `Trigger worker ${worker.name} immediately`;

  return {
    isOperational,
    isExecuting,
    isConfirming,
    notice,
    handleTriggerClick,
    tooltipText,
  };
}
