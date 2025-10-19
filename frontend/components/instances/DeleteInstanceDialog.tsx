"use client";

import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Loader2, AlertTriangle } from "lucide-react";
import { deleteInstance } from "@/lib/api";
import { toast } from "sonner";

interface DeleteInstanceDialogProps {
  instanceId: string | null;
  instanceName: string | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDeleted: () => void;
}

export function DeleteInstanceDialog({
  instanceId,
  instanceName,
  open,
  onOpenChange,
  onDeleted,
}: DeleteInstanceDialogProps) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    if (!instanceId) return;

    setLoading(true);

    try {
      await deleteInstance(instanceId);
      toast.success("Instance deleted successfully");
      onOpenChange(false);
      onDeleted();
    } catch (error: any) {
      console.error("Failed to delete instance:", error);
      toast.error("Failed to delete instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="p-2 bg-red-100 rounded-lg">
              <AlertTriangle className="h-5 w-5 text-red-600" />
            </div>
            <DialogTitle>Delete Instance</DialogTitle>
          </div>
          <DialogDescription className="pt-3">
            Are you sure you want to delete <strong>{instanceName}</strong>?
          </DialogDescription>
        </DialogHeader>
        <div className="py-4">
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
            <p className="text-sm text-amber-900">
              <strong>Warning:</strong> This will stop and remove the Docker container. 
              Your data will be preserved in the storage directory but the instance will no longer be accessible.
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={loading}
          >
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Deleting...
              </>
            ) : (
              "Delete Instance"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
