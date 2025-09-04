"use client";

import { useCallback, useEffect, useRef, useState } from "react";

export default function Home() {
  const [VMID, setVMID] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const rfbRef = useRef<any>(null);
  const backendBase = useRef<string>("");

  useEffect(() => {
    const host = window.location.hostname;
    const proto = window.location.protocol === "https:" ? "https" : "http";
    backendBase.current = `${proto}://${host}:8080`;
  }, []);

  const connect = useCallback(async (id: string) => {
    if (!containerRef.current) return;
    const wsProto = location.protocol === 'https:' ? 'wss' : 'ws';
    const directUrl = `${wsProto}://${location.hostname}:8080/vm/${id}/stream`;
    const proxiedUrl = `${wsProto}://${location.host}/vm/${id}/stream`;
    let wsUrl = directUrl;

    // @ts-ignore - novnc has no types here
    const { default: RFB } = await import("@novnc/novnc/lib/rfb");
    containerRef.current.innerHTML = "";
    let rfb: any;
    try {
      rfb = new RFB(containerRef.current, wsUrl);
    } catch {
      wsUrl = proxiedUrl;
      rfb = new RFB(containerRef.current, wsUrl);
    }
    rfb.scaleViewport = true;
    rfb.resizeSession = true;
    rfbRef.current = rfb;
  }, []);

  const startvm = useCallback(async () => {
    if (VMID) return; // already started
    try {
      const res = await fetch(`${backendBase.current}/vm/start`, { method: "POST" });
      if (!res.ok) return;
      const data = await res.json();
      const id = data.VMID as string;
      setVMID(id);

      const start = Date.now();
      while (Date.now() - start < 20000) {
        try {
          await connect(id);
          break;
        } catch {
          await new Promise(r => setTimeout(r, 800));
        }
      }
    } catch {}
  }, [connect, VMID]);

  useEffect(() => {
    const handleUnload = () => {
      if (VMID) {
        navigator.sendBeacon(`${backendBase.current}/vm/${VMID}/stop`);
      }
    };
    window.addEventListener("beforeunload", handleUnload);
    return () => {
      window.removeEventListener("beforeunload", handleUnload);
      if (rfbRef.current) {
        try { rfbRef.current.disconnect(); } catch {}
      }
      if (VMID) {
        fetch(`${backendBase.current}/vm/${VMID}/stop`, { method: "POST", keepalive: true }).catch(() => {});
      }
    };
  }, [VMID]);

  return (
    <div style={{
      width: "100vw",
      height: "100vh",
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      background: "#111"
    }}>
      {!VMID && (
        <button
          onClick={startvm}
          style={{
            padding: "16px 32px",
            fontSize: 24,
            borderRadius: 8,
            border: "none",
            cursor: "pointer",
          }}
        >
          Start
        </button>
      )}
      <div
        ref={containerRef}
        style={{ width: "100%", height: "100%", display: VMID ? "block" : "none" }}
      />
    </div>
  );
}
