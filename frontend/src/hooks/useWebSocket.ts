import { useEffect } from "react";
import type { SocketCallback } from "@/services/ws";
import { wsManager } from "@/services/ws";

export function useWebSocket(cb: SocketCallback, url?: string) {
  useEffect(() => {
    const id = wsManager.add(cb, url);
    return () => wsManager.remove(id);
  }, [cb, url]);
}
