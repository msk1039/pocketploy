"use client";

import { useState } from "react";
import { Instance } from "@/types/instance";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Play, Square, RotateCw, Trash2, Loader2 } from "lucide-react";
import { startInstance, stopInstance, restartInstance } from "@/lib/api";
import { toast } from "sonner";
import { DeleteInstanceDialog } from "./DeleteInstanceDialog";
import { useRouter } from "next/navigation";

interface InstanceControlsProps {
  instance: Instance;
  onInstanceUpdated: () => void;
}

export function InstanceControls({ instance, onInstanceUpdated }: InstanceControlsProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleStart = async () => {
    setLoading(true);
    try {
      await startInstance(instance.id);
      toast.success("Instance started successfully");
      onInstanceUpdated();
    } catch (error: any) {
      toast.error("Failed to start instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleStop = async () => {
    setLoading(true);
    try {
      await stopInstance(instance.id);
      toast.success("Instance stopped successfully");
      onInstanceUpdated();
    } catch (error: any) {
      toast.error("Failed to stop instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleRestart = async () => {
    setLoading(true);
    try {
      await restartInstance(instance.id);
      toast.success("Instance restarted successfully");
      onInstanceUpdated();
    } catch (error: any) {
      toast.error("Failed to restart instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = () => {
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirmed = () => {
    toast.success("Instance deleted successfully");
    router.push("/dashboard");
  };

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Instance Controls</CardTitle>
          <CardDescription>Manage your instance lifecycle</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <Button
            onClick={handleStart}
            disabled={loading || instance.status === "running" || instance.status === "creating"}
            className="w-full flex items-center justify-center gap-2"
            variant="default"
          >
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Play className="h-4 w-4" />}
            Start Instance
          </Button>

          <Button
            onClick={handleStop}
            disabled={loading || instance.status === "stopped" || instance.status === "creating"}
            className="w-full flex items-center justify-center gap-2"
            variant="outline"
          >
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Square className="h-4 w-4" />}
            Stop Instance
          </Button>

          <Button
            onClick={handleRestart}
            disabled={loading || instance.status === "stopped" || instance.status === "creating"}
            className="w-full flex items-center justify-center gap-2"
            variant="outline"
          >
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <RotateCw className="h-4 w-4" />}
            Restart Instance
          </Button>

          <div className="pt-4 border-t">
            <Button
              onClick={handleDelete}
              disabled={loading}
              className="w-full flex items-center justify-center gap-2 bg-red-600 hover:bg-red-700 text-white"
              variant="destructive"
            >
              <Trash2 className="h-4 w-4" />
              Delete Instance
            </Button>
          </div>
        </CardContent>
      </Card>

      <DeleteInstanceDialog
        instanceId={instance.id}
        instanceName={instance.name}
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onDeleted={handleDeleteConfirmed}
      />
    </>
  );
}
