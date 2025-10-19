"use client";

import { useState, useEffect, useRef } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { RefreshCw, Download, Terminal } from "lucide-react";
import { getInstanceLogs } from "@/lib/api";
import { toast } from "sonner";

interface InstanceLogsProps {
  instanceId: string;
}

export function InstanceLogs({ instanceId }: InstanceLogsProps) {
  const [logs, setLogs] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [tail, setTail] = useState("100");
  const [autoRefresh, setAutoRefresh] = useState(false);
  const logsEndRef = useRef<HTMLDivElement>(null);

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const response = await getInstanceLogs(instanceId, tail);
      setLogs(response.logs);
      // Auto-scroll to bottom
      setTimeout(() => {
        logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
      }, 100);
    } catch (error: any) {
      console.error("Failed to fetch logs:", error);
      toast.error("Failed to load logs", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, [instanceId, tail]);

  // Auto-refresh logs every 5 seconds if enabled
  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(() => {
      fetchLogs();
    }, 5000);

    return () => clearInterval(interval);
  }, [autoRefresh, instanceId, tail]);

  const handleRefresh = () => {
    fetchLogs();
    toast.info("Refreshing logs...");
  };

  const handleDownload = () => {
    const blob = new Blob([logs], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `instance-${instanceId}-logs.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    toast.success("Logs downloaded successfully");
  };

  const handleTailChange = (newTail: string) => {
    setTail(newTail);
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-xl flex items-center gap-2">
              <Terminal className="h-5 w-5" />
              Container Logs
            </CardTitle>
            <CardDescription>Real-time logs from your PocketBase instance</CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <select
              value={tail}
              onChange={(e) => handleTailChange(e.target.value)}
              className="px-3 py-1 text-sm border rounded-md bg-white"
            >
              <option value="50">Last 50 lines</option>
              <option value="100">Last 100 lines</option>
              <option value="500">Last 500 lines</option>
              <option value="all">All logs</option>
            </select>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setAutoRefresh(!autoRefresh)}
              className={autoRefresh ? "bg-blue-50 border-blue-300" : ""}
            >
              {autoRefresh ? "Auto-refresh ON" : "Auto-refresh OFF"}
            </Button>
            <Button variant="outline" size="sm" onClick={handleRefresh} disabled={loading}>
              <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
            </Button>
            <Button variant="outline" size="sm" onClick={handleDownload} disabled={!logs}>
              <Download className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="bg-gray-900 rounded-lg p-4 h-[500px] overflow-y-auto font-mono text-sm">
          {loading && !logs ? (
            <div className="flex items-center justify-center h-full text-gray-500">
              <RefreshCw className="h-6 w-6 animate-spin mr-2" />
              Loading logs...
            </div>
          ) : !logs ? (
            <div className="flex items-center justify-center h-full text-gray-500">
              No logs available
            </div>
          ) : (
            <pre className="text-green-400 whitespace-pre-wrap break-words">
              {logs}
              <div ref={logsEndRef} />
            </pre>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
