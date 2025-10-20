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
} from "@/components/coss-ui/dialog";
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
  const [adminEmail, setAdminEmail] = useState("");
  const [adminPassword, setAdminPassword] = useState("");
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

    // Validate email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(adminEmail)) {
      setError("Please enter a valid email address");
      return;
    }

    // Validate password
    if (adminPassword.length < 10) {
      setError("Admin password must be at least 10 characters");
      return;
    }

    setLoading(true);

    try {
      const response = await createInstance({ 
        name, 
        admin_email: adminEmail, 
        admin_password: adminPassword 
      });
      toast.success("Instance created successfully!", {
        description: `Access at: ${response.url}`,
      });
      setName("");
      setAdminEmail("");
      setAdminPassword("");
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
      <DialogTrigger render={<Button className="flex items-center gap-2" />}>
       
          <Plus className="h-4 w-4" />
          Create Instance

      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create PocketBase Instance</DialogTitle>
            <DialogDescription>
              Create a new isolated PocketBase database instance with admin credentials. You can create up to 5 instances.
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
            
            <div className="grid gap-2">
              <Label htmlFor="adminEmail">Admin Email</Label>
              <Input
                id="adminEmail"
                type="email"
                placeholder="admin@example.com"
                value={adminEmail}
                onChange={(e) => setAdminEmail(e.target.value)}
                disabled={loading}
                required
              />
              <p className="text-sm text-gray-500">
                Email for the PocketBase admin account
              </p>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="adminPassword">Admin Password</Label>
              <Input
                id="adminPassword"
                type="password"
                placeholder="Enter a secure password"
                value={adminPassword}
                onChange={(e) => setAdminPassword(e.target.value)}
                disabled={loading}
                required
                minLength={10}
              />
              <p className="text-sm text-gray-500">
                Password must be at least 10 characters
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
