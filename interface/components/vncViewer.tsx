"use client";

import React, { useEffect, useRef, useState } from 'react';

interface VNCViewerProps {
  VMID: string;
  className?: string;
}

export const VNCViewer: React.FC<VNCViewerProps> = ({ VMID, className = '' }) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const rfbRef = useRef<any>(null);
  const [status, setStatus] = useState('Connecting...');

  useEffect(() => {
    if (!VMID || !containerRef.current) return;

    const connect = async () => {
      try {
        const token = localStorage.getItem('authToken');
        if (!token) {
          setStatus('No authentication token');
          return;
        }

        const wsUrl = `ws://localhost:8080/vnc/vms/${VMID}/stream?token=${encodeURIComponent(token)}`;

        // Dynamic import for client-side only
        // @ts-ignore - novnc library doesn't have proper TypeScript types
        const { default: RFB } = await import('@novnc/novnc/lib/rfb');

        const container = containerRef.current;
        if (!container) return;

        container.innerHTML = '';

        const rfb = new RFB(container, wsUrl);
        rfbRef.current = rfb;

        // Configure RFB for proper scaling
        rfb.scaleViewport = true;
        rfb.resizeSession = true;

        rfb.addEventListener('connect', () => {
          setStatus('Connected');
          // Ensure proper scaling after connection
          rfb.scaleViewport = true;
        });

        rfb.addEventListener('disconnect', (e: any) => {
          setStatus(`Disconnected: ${e?.detail?.reason || 'Unknown'}`);
        });

        setStatus('Handshaking...');

      } catch (err) {
        setStatus(`Error: ${err instanceof Error ? err.message : 'Unknown error'}`);
      }
    };

    const timer = setTimeout(connect, 100);

    return () => {
      clearTimeout(timer);
      if (rfbRef.current) {
        try {
          rfbRef.current.disconnect();
        } catch (disconnectError) {
          console.warn('Error during RFB disconnect:', disconnectError);
        }
        rfbRef.current = null;
      }
    };
  }, [VMID]);

  return (
    <div className={`relative w-full h-full ${className}`}>
      {status !== 'Connected' && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-50">
          <div className="text-center p-4 bg-white rounded-lg shadow">
            {status}
          </div>
        </div>
      )}
      <div
        ref={containerRef}
        className="w-full h-full bg-black"
        style={{ 
          minHeight: '400px',
          maxHeight: '100%',
          overflow: 'hidden'
        }}
      />
    </div>
  );
}; 