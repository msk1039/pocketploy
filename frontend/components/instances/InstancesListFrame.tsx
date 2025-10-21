"use client";

import { Instance } from "@/types/instance";
import {
  Database,
  ExternalLink,
  MoreVertical,
  Play,
  Square,
  Trash2,
  Loader2,
  RotateCw,
  Copy,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Frame,
  FrameDescription,
  FrameHeader,
  FramePanel,
  FrameTitle,
} from "@/components/coss-ui/frame";
import {
  Menu,
  MenuItem,
  MenuPopup,
  MenuSeparator,
  MenuTrigger,
} from "@/components/coss-ui/menu";
import { useState } from "react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { startInstance, stopInstance, restartInstance } from "@/lib/api";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/coss-ui/dialog";

interface InstancesListFrameProps {
  instances: Instance[];
  loading: boolean;
  onDelete: (id: string) => void;
  onInstanceUpdated: () => void;
}

export function InstancesListFrame({
  instances,
  loading,
  onDelete,
  onInstanceUpdated,
}: InstancesListFrameProps) {
  const router = useRouter();
  const [loadingInstance, setLoadingInstance] = useState<string | null>(null);
  const [stopDialogOpen, setStopDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedInstance, setSelectedInstance] = useState<Instance | null>(
    null
  );

  const getStatusColor = (status: string) => {
    switch (status) {
      case "running":
        return "bg-green-100 text-green-800 border-green-200";
      case "creating":
      case "pending":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "stopped":
        return "bg-gray-100 text-gray-800 border-gray-200";
      case "failed":
        return "bg-red-100 text-red-800 border-red-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const handleStart = async (instance: Instance, e: React.MouseEvent) => {
    e.stopPropagation();
    setLoadingInstance(instance.id);
    try {
      await startInstance(instance.id);
      toast.success("Instance started successfully");
      onInstanceUpdated();
    } catch (error: any) {
      toast.error("Failed to start instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoadingInstance(null);
    }
  };

  const handleStopClick = (instance: Instance, e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedInstance(instance);
    setStopDialogOpen(true);
  };

  const handleStopConfirm = async () => {
    if (!selectedInstance) return;
    setLoadingInstance(selectedInstance.id);
    try {
      await stopInstance(selectedInstance.id);
      toast.success("Instance stopped successfully");
      onInstanceUpdated();
    } catch (error: any) {
      toast.error("Failed to stop instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoadingInstance(null);
      setStopDialogOpen(false);
      setSelectedInstance(null);
    }
  };

  const handleRestart = async (instance: Instance, e: React.MouseEvent) => {
    e.stopPropagation();
    setLoadingInstance(instance.id);
    try {
      await restartInstance(instance.id);
      toast.success("Instance restarted successfully");
      onInstanceUpdated();
    } catch (error: any) {
      toast.error("Failed to restart instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoadingInstance(null);
    }
  };

  const handleDeleteClick = (instance: Instance, e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedInstance(instance);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    if (!selectedInstance) return;
    onDelete(selectedInstance.id);
    setDeleteDialogOpen(false);
    setSelectedInstance(null);
  };

  const handleCopyUrl = (instance: Instance, e: React.MouseEvent) => {
    e.stopPropagation();
    const url = `http://${instance.subdomain}/_/`;
    navigator.clipboard.writeText(url);
    toast.success("URL copied to clipboard!");
  };

  const handleOpenInstance = (instance: Instance, e: React.MouseEvent) => {
    e.stopPropagation();
    const url = `http://${instance.subdomain}/_/`;
    window.open(url, "_blank");
  };

  const handleViewDetails = (instance: Instance) => {
    router.push(`/dashboard/instances/${instance.id}`);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (instances.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <div className="p-4 bg-gray-100 rounded-full mb-4">
          <Database className="h-8 w-8 text-gray-400" />
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">
          No instances yet
        </h3>
        <p className="text-gray-600 mb-6 max-w-md">
          Create your first PocketBase instance to get started. Each instance is
          isolated and accessible via a unique subdomain.
        </p>
      </div>
    );
  }

  return (
    <>
    
        <FramePanel>
          <div className="divide-y">
            {instances.map((instance) => {
              const isCreating = instance.status === "creating";
              const isLoading = loadingInstance === instance.id;

              return (
                <div
                  key={instance.id}
                  className="flex items-center justify-between p-4 hover:bg-muted/50 transition-colors cursor-pointer"
                  onClick={() => handleViewDetails(instance)}
                >
                  {/* Left section: Logo, Name, URL */}
                  <div className="flex items-center gap-4 flex-1 min-w-0">
                    <div className="p-2 bg-blue-50 rounded-lg flex-shrink-0">
                      <Database className="h-5 w-5 text-blue-600" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="font-semibold text-base truncate">
                          {instance.name}
                        </h3>
                        <Badge
                          className={getStatusColor(instance.status)}
                          variant="outline"
                        >
                          {isCreating && (
                            <Loader2 className="h-3 w-3 mr-1 animate-spin" />
                          )}
                          {instance.status}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <code className="text-xs bg-gray-100 px-2 py-0.5 rounded">
                          {instance.subdomain}/_/
                        </code>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-6 w-6 p-0"
                          onClick={(e) => handleCopyUrl(instance, e)}
                        >
                          <Copy className="h-3 w-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-6 w-6 p-0"
                          onClick={(e) => handleOpenInstance(instance, e)}
                          disabled={instance.status !== "running"}
                        >
                          <ExternalLink className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>
                  </div>

                  {/* Right section: Actions menu */}
                  <div
                    className="flex-shrink-0"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Menu>
                      <MenuTrigger
                        render={
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-8 w-8 p-0"
                            disabled={isLoading}
                          />
                        }
                      >
                        {isLoading ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <MoreVertical className="h-4 w-4" />
                        )}
                      </MenuTrigger>
                      <MenuPopup>
                        <MenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            handleStart(instance, e);
                          }}
                          disabled={
                            instance.status === "running" ||
                            instance.status === "creating" ||
                            isLoading
                          }
                        >
                          <Play className="opacity-72" />
                          Start Instance
                        </MenuItem>
                        <MenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            handleStopClick(instance, e);
                          }}
                          disabled={
                            instance.status === "stopped" ||
                            instance.status === "creating" ||
                            isLoading
                          }
                        >
                          <Square className="opacity-72" />
                          Stop Instance
                        </MenuItem>
                        <MenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            handleRestart(instance, e);
                          }}
                          disabled={
                            instance.status === "stopped" ||
                            instance.status === "creating" ||
                            isLoading
                          }
                        >
                          <RotateCw className="opacity-72" />
                          Restart Instance
                        </MenuItem>
                        <MenuSeparator />
                        <MenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            handleViewDetails(instance);
                          }}
                        >
                          <ExternalLink className="opacity-72" />
                          View Details
                        </MenuItem>
                        <MenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            handleCopyUrl(instance, e);
                          }}
                        >
                          <Copy className="opacity-72" />
                          Copy URL
                        </MenuItem>
                        <MenuSeparator />
                        <MenuItem
                          variant="destructive"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteClick(instance, e);
                          }}
                          disabled={isLoading}
                        >
                          <Trash2 className="opacity-72" />
                          Delete Instance
                        </MenuItem>
                      </MenuPopup>
                    </Menu>
                  </div>
                </div>
              );
            })}
          </div>
        </FramePanel>
      

      {/* Stop Confirmation Dialog */}
      <Dialog open={stopDialogOpen} onOpenChange={setStopDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Stop Instance</DialogTitle>
            <DialogDescription>
              Are you sure you want to stop{" "}
              <strong>{selectedInstance?.name}</strong>?
            </DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <p className="text-sm text-muted-foreground">
              The instance will be stopped and will no longer be accessible. You
              can start it again later.
            </p>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setStopDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button type="button" onClick={handleStopConfirm}>
              Stop Instance
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="text-destructive">
              Delete Instance
            </DialogTitle>
            <DialogDescription>
              Are you sure you want to delete{" "}
              <strong>{selectedInstance?.name}</strong>?
            </DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
              <p className="text-sm text-destructive">
                <strong>Warning:</strong> This will stop and remove the Docker
                container. Your data will be preserved in the storage directory
                but the instance will no longer be accessible.
              </p>
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDeleteConfirm}
            >
              Delete Instance
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      
    </>
  );
}
