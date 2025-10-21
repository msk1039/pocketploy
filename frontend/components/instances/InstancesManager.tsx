"use client";

import { useState, useEffect } from "react";
import { Instance } from "@/types/instance";
import { listInstances } from "@/lib/api";
import { CreateInstanceDialog } from "./CreateInstanceDialog";
import { InstancesListFrame } from "./InstancesListFrame";
import { DeleteInstanceDialog } from "./DeleteInstanceDialog";
import { Button } from "@/components/ui/button";
import { RefreshCw } from "lucide-react";
import { toast } from "sonner";
import {
  Frame,
  FrameDescription,
  FrameHeader,
  FramePanel,
  FrameTitle,
} from "@/components/coss-ui/frame";

export function InstancesManager() {
  const [instances, setInstances] = useState<Instance[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [instanceToDelete, setInstanceToDelete] = useState<{
    id: string;
    name: string;
  } | null>(null);

  const fetchInstances = async () => {
    setLoading(true);
    try {
      const response = await listInstances();
      setInstances(response.instances || []);
    } catch (error: any) {
      console.error("Failed to fetch instances:", error);
      toast.error("Failed to load instances", {
        description: error.message || "An error occurred",
      });
      setInstances([]); // Set to empty array on error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInstances();
  }, []);

  // Poll for instance status updates when any instance is in "creating" state
  useEffect(() => {
    const hasCreatingInstances = instances.some(
      (instance) => instance.status === "creating"
    );

    if (!hasCreatingInstances) {
      return;
    }

    // Poll every 2 seconds
    const intervalId = setInterval(() => {
      fetchInstances();
    }, 2000);

    return () => clearInterval(intervalId);
  }, [instances]);

  const handleDelete = (id: string) => {
    const instance = instances.find((i) => i.id === id);
    if (instance) {
      setInstanceToDelete({ id: instance.id, name: instance.name });
      setDeleteDialogOpen(true);
    }
  };

  const handleDeleteConfirm = () => {
    fetchInstances();
    setInstanceToDelete(null);
  };

  const handleRefresh = () => {
    toast.info("Refreshing instances...");
    fetchInstances();
  };

  return (
    <div className="space-y-6">
      {/* Header */}

      {/* Instances List */}
      <Frame>
        <FrameHeader>
          <div className="flex items-center justify-between">
            <div>
              <FrameTitle>Your PocketBase Instances</FrameTitle>
              <FrameDescription>
                Manage your database instances ({instances.length}/5)
              </FrameDescription>
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={handleRefresh}
                disabled={loading}
                className="flex items-center gap-2"
              >
                <RefreshCw
                  className={`h-4 w-4 ${loading ? "animate-spin" : ""}`}
                />
                Refresh
              </Button>
              <CreateInstanceDialog onInstanceCreated={fetchInstances} />
            </div>
          </div>
        </FrameHeader>
        <InstancesListFrame
          instances={instances}
          loading={loading}
          onDelete={handleDelete}
          onInstanceUpdated={fetchInstances}
        />
      </Frame>

      {/* Delete Confirmation Dialog */}
      <DeleteInstanceDialog
        instanceId={instanceToDelete?.id || null}
        instanceName={instanceToDelete?.name || null}
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onDeleted={handleDeleteConfirm}
      />
    </div>
  );
}
