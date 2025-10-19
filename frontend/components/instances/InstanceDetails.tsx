"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Instance } from "@/types/instance";
import { getInstance } from "@/lib/api";
import { InstanceMetadata } from "./InstanceMetadata";
import { InstanceControls } from "./InstanceControls";
import { InstanceLogs } from "./InstanceLogs";
import { Button } from "@/components/ui/button";
import { ArrowLeft, RefreshCw } from "lucide-react";
import { toast } from "sonner";

interface InstanceDetailsProps {
  instanceId: string;
}

export function InstanceDetails({ instanceId }: InstanceDetailsProps) {
  const router = useRouter();
  const [instance, setInstance] = useState<Instance | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshKey, setRefreshKey] = useState(0);

  const fetchInstance = async () => {
    setLoading(true);
    try {
      const response = await getInstance(instanceId);
      setInstance(response.instance);
    } catch (error: any) {
      console.error("Failed to fetch instance:", error);
      toast.error("Failed to load instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInstance();
  }, [instanceId, refreshKey]);

  const handleRefresh = () => {
    setRefreshKey((prev) => prev + 1);
    toast.info("Refreshing instance details...");
  };

  const handleBack = () => {
    router.push("/dashboard");
  };

  const handleInstanceUpdated = () => {
    fetchInstance();
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (!instance) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Instance not found</h2>
        <p className="text-gray-600 mb-6">The instance you're looking for doesn't exist.</p>
        <Button onClick={handleBack}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Dashboard
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="outline" onClick={handleBack} className="flex items-center gap-2">
            <ArrowLeft className="h-4 w-4" />
            Back
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">{instance.name}</h1>
            <p className="text-gray-600 mt-1">Instance Details</p>
          </div>
        </div>
        <Button variant="outline" onClick={handleRefresh} className="flex items-center gap-2">
          <RefreshCw className="h-4 w-4" />
          Refresh
        </Button>
      </div>

      {/* Metadata and Controls */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <InstanceMetadata instance={instance} />
        </div>
        <div>
          <InstanceControls instance={instance} onInstanceUpdated={handleInstanceUpdated} />
        </div>
      </div>

      {/* Logs */}
      <InstanceLogs instanceId={instanceId} />
    </div>
  );
}
