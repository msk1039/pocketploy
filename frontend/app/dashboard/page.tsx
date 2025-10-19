"use client";

import { useAuth } from "@/contexts/AuthContext";
import ProtectedRoute from "@/components/ProtectedRoute";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { LogOut, User, Mail, Calendar, CheckCircle } from "lucide-react";
import { toast } from "sonner";

function DashboardContent() {
  const { user, logout } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
      toast.success("Logged out successfully");
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
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <h1 className="text-2xl font-bold text-gray-900">pocketploy</h1>
            </div>
            <Button
              variant="outline"
              onClick={handleLogout}
              className="flex items-center space-x-2"
            >
              <LogOut className="h-4 w-4" />
              <span>Sign out</span>
            </Button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="space-y-8">
          {/* Welcome Section */}
          <div className="flex items-center space-x-4">
            <Avatar className="h-16 w-16">
              <AvatarFallback className="bg-blue-600 text-white text-xl">
                {user?.username ? getInitials(user.username) : "U"}
              </AvatarFallback>
            </Avatar>
            <div>
              <h2 className="text-3xl font-bold text-gray-900">
                Welcome back, {user?.username}!
              </h2>
              <p className="text-gray-600 mt-1">
                Your pocketploy dashboard
              </p>
            </div>
          </div>

          {/* User Information Card */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <User className="h-5 w-5" />
                <span>Account Information</span>
              </CardTitle>
              <CardDescription>
                Your account details and authentication status
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-1">
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <User className="h-4 w-4" />
                    <span className="font-medium">Username</span>
                  </div>
                  <p className="text-lg font-semibold">{user?.username}</p>
                </div>

                <div className="space-y-1">
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <Mail className="h-4 w-4" />
                    <span className="font-medium">Email</span>
                  </div>
                  <p className="text-lg font-semibold">{user?.email}</p>
                </div>

                <div className="space-y-1">
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <Calendar className="h-4 w-4" />
                    <span className="font-medium">Account Created</span>
                  </div>
                  <p className="text-sm">{formatDate(user?.created_at)}</p>
                </div>

                <div className="space-y-1">
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <Calendar className="h-4 w-4" />
                    <span className="font-medium">Last Login</span>
                  </div>
                  <p className="text-sm">{formatDate(user?.last_login_at)}</p>
                </div>

                <div className="space-y-1">
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <CheckCircle className="h-4 w-4" />
                    <span className="font-medium">Account Status</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span
                      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                        user?.is_active
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {user?.is_active ? "Active" : "Inactive"}
                    </span>
                  </div>
                </div>

                <div className="space-y-1">
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <span className="font-medium">User ID</span>
                  </div>
                  <p className="text-xs font-mono text-gray-600">{user?.id}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Phase 2 Placeholder */}
          <Card>
            <CardHeader>
              <CardTitle>PocketBase Instances</CardTitle>
              <CardDescription>
                Manage your PocketBase database instances
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="rounded-lg border-2 border-dashed border-gray-300 p-12 text-center">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
                  <svg
                    className="h-6 w-6 text-gray-600"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                    />
                  </svg>
                </div>
                <h3 className="mt-4 text-lg font-medium text-gray-900">
                  No instances yet
                </h3>
                <p className="mt-2 text-sm text-gray-500">
                  Instance creation will be available in Phase 2
                </p>
                <p className="mt-1 text-xs text-gray-400">
                  Coming soon: Create and manage isolated PocketBase instances
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Authentication Status Card */}
          <Card className="bg-green-50 border-green-200">
            <CardContent className="pt-6">
              <div className="flex items-center space-x-3">
                <CheckCircle className="h-6 w-6 text-green-600" />
                <div>
                  <p className="font-semibold text-green-900">
                    Authentication Successful!
                  </p>
                  <p className="text-sm text-green-700">
                    Your JWT authentication is working correctly. Tokens are being
                    automatically refreshed.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <DashboardContent />
    </ProtectedRoute>
  );
}
