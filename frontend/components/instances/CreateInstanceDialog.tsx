"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, Loader2 } from "lucide-react";
import { createInstance } from "@/lib/api";
import { toast } from "sonner";

interface CreateInstanceDialogProps {
  onInstanceCreated: () => void;
}

export function CreateInstanceDialog({ onInstanceCreated }: CreateInstanceDialogProps) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Validate name
    if (name.length < 3 || name.length > 100) {
      setError("Instance name must be between 3 and 100 characters");
      return;
    }

    setLoading(true);

    try {
      const response = await createInstance({ name });
      toast.success("Instance created successfully!", {
        description: `Access at: ${response.url}`,
      });
      setName("");
      setOpen(false);
      onInstanceCreated();
    } catch (error: any) {
      console.error("Failed to create instance:", error);
      const errorMessage = error.message || "Failed to create instance";
      setError(errorMessage);
      toast.error("Failed to create instance", {
        description: errorMessage,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="flex items-center gap-2">
          <Plus className="h-4 w-4" />
          Create Instance
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create PocketBase Instance</DialogTitle>
            <DialogDescription>
              Create a new isolated PocketBase database instance. You can create up to 5 instances.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Instance Name</Label>
              <Input
                id="name"
                placeholder="My Project Database"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={loading}
                required
                minLength={3}
                maxLength={100}
              />
              <p className="text-sm text-gray-500">
                Letters, numbers, spaces, hyphens, and underscores allowed (3-100 characters)
              </p>
            </div>
            {error && (
              <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md">
                {error}
              </div>
            )}
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Creating...
                </>
              ) : (
                "Create Instance"
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
