"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, Plus, Info } from "lucide-react";
import { createInstance } from "@/lib/api";
import { toast } from "sonner";

export default function CreateInstancePage() {
  const router = useRouter();
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
      
      // Redirect to instances page after creation
      router.push("/dashboard/instances");
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
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Create New Instance</h2>
        <p className="text-muted-foreground mt-1">
          Set up a new PocketBase instance with admin credentials
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {/* Form Card */}
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>Instance Configuration</CardTitle>
            <CardDescription>
              Configure your new PocketBase instance. All fields are required.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-2">
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
                <p className="text-sm text-muted-foreground">
                  Letters, numbers, spaces, hyphens, and underscores allowed (3-100 characters)
                </p>
              </div>

              <div className="space-y-2">
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
                <p className="text-sm text-muted-foreground">
                  Email for the PocketBase admin account
                </p>
              </div>

              <div className="space-y-2">
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
                <p className="text-sm text-muted-foreground">
                  Password must be at least 10 characters
                </p>
              </div>

              {error && (
                <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md border border-red-200">
                  {error}
                </div>
              )}

              <div className="flex gap-3 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => router.push("/dashboard/instances")}
                  disabled={loading}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={loading} className="flex-1">
                  {loading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Creating...
                    </>
                  ) : (
                    <>
                      <Plus className="mr-2 h-4 w-4" />
                      Create Instance
                    </>
                  )}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* Info Card */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Info className="h-4 w-4" />
              What happens next?
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-3 text-sm">
              <div className="space-y-1">
                <h4 className="font-medium">1. Instance Creation</h4>
                <p className="text-muted-foreground">
                  We'll create a Docker container with PocketBase
                </p>
              </div>

              <div className="space-y-1">
                <h4 className="font-medium">2. Admin Setup</h4>
                <p className="text-muted-foreground">
                  Your admin credentials will be configured automatically
                </p>
              </div>

              <div className="space-y-1">
                <h4 className="font-medium">3. Subdomain Assignment</h4>
                <p className="text-muted-foreground">
                  A unique subdomain will be generated for access
                </p>
              </div>

              <div className="space-y-1">
                <h4 className="font-medium">4. Ready to Use</h4>
                <p className="text-muted-foreground">
                  Access your admin panel and start building
                </p>
              </div>
            </div>

            <div className="pt-4 border-t space-y-2">
              <h4 className="font-medium text-sm">Quick Tips:</h4>
              <ul className="text-xs text-muted-foreground space-y-1">
                <li>• Use a strong, unique password</li>
                <li>• Instance names can be changed later</li>
                <li>• You can create up to 5 instances</li>
                <li>• Data persists across restarts</li>
              </ul>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
