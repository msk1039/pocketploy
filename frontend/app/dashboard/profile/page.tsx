"use client";

import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { User, Mail, Calendar, CheckCircle, LogOut, Shield } from "lucide-react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

export default function ProfilePage() {
  const { user, logout } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await logout();
      toast.success("Logged out successfully");
      router.push("/login");
    } catch (error) {
      toast.error("Logout failed");
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "N/A";
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getInitials = (username: string) => {
    return username.substring(0, 2).toUpperCase();
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Profile</h2>
        <p className="text-muted-foreground mt-1">
          Manage your account settings and view your information
        </p>
      </div>

      {/* Profile Header */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center gap-4">
            <Avatar className="h-20 w-20">
              <AvatarFallback className="bg-blue-600 text-white text-2xl">
                {user?.username ? getInitials(user.username) : "U"}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <h3 className="text-2xl font-bold">{user?.username}</h3>
              <p className="text-muted-foreground">{user?.email}</p>
              <div className="mt-2">
                <Badge variant={user?.is_active ? "default" : "destructive"}>
                  {user?.is_active ? "Active" : "Inactive"}
                </Badge>
              </div>
            </div>
            <Button variant="destructive" onClick={handleLogout}>
              <LogOut className="h-4 w-4 mr-2" />
              Sign Out
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Account Information */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5" />
              Account Details
            </CardTitle>
            <CardDescription>
              Your personal information and account status
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <User className="h-4 w-4" />
                <span className="font-medium">Username</span>
              </div>
              <p className="text-lg font-semibold">{user?.username}</p>
            </div>

            <div className="space-y-1">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Mail className="h-4 w-4" />
                <span className="font-medium">Email Address</span>
              </div>
              <p className="text-lg font-semibold">{user?.email}</p>
            </div>

            <div className="space-y-1">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle className="h-4 w-4" />
                <span className="font-medium">Account Status</span>
              </div>
              <Badge variant={user?.is_active ? "default" : "destructive"}>
                {user?.is_active ? "Active" : "Inactive"}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              Activity Timeline
            </CardTitle>
            <CardDescription>
              Your account activity and important dates
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Calendar className="h-4 w-4" />
                <span className="font-medium">Account Created</span>
              </div>
              <p className="text-sm">{formatDate(user?.created_at)}</p>
            </div>

            <div className="space-y-1">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Calendar className="h-4 w-4" />
                <span className="font-medium">Last Updated</span>
              </div>
              <p className="text-sm">{formatDate(user?.updated_at)}</p>
            </div>

            <div className="space-y-1">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Calendar className="h-4 w-4" />
                <span className="font-medium">Last Login</span>
              </div>
              <p className="text-sm">{formatDate(user?.last_login_at)}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Account ID */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Advanced Information
          </CardTitle>
          <CardDescription>
            Technical details about your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-1">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <span className="font-medium">User ID</span>
            </div>
            <div className="flex items-center gap-2">
              <code className="text-xs font-mono bg-gray-100 px-2 py-1 rounded">
                {user?.id}
              </code>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  navigator.clipboard.writeText(user?.id || "");
                  toast.success("User ID copied to clipboard");
                }}
              >
                Copy
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Actions */}
      <Card className="border-dashed border-2">
        <CardHeader>
          <CardTitle className="text-base">Account Actions</CardTitle>
          <CardDescription>
            Manage your account settings (Coming soon)
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <Button variant="outline" className="w-full" disabled>
            Change Password
          </Button>
          <Button variant="outline" className="w-full" disabled>
            Update Email
          </Button>
          <Button variant="outline" className="w-full" disabled>
            Enable Two-Factor Authentication
          </Button>
          <Button variant="destructive" className="w-full" disabled>
            Delete Account
          </Button>
          <p className="text-xs text-muted-foreground text-center pt-2">
            These features will be available in a future update
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
