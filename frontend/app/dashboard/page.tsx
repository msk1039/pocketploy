"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Database, Plus, Archive, User, TrendingUp } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { listInstances } from "@/lib/api";

export default function DashboardPage() {
  const { user } = useAuth();
  const router = useRouter();
  const [instanceCount, setInstanceCount] = useState<number>(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchInstanceCount = async () => {
      try {
        const response = await listInstances();
        setInstanceCount(response.instances?.length || 0);
      } catch (error) {
        console.error("Failed to fetch instances:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchInstanceCount();
  }, []);

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return "Good morning";
    if (hour < 18) return "Good afternoon";
    return "Good evening";
  };

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          {getGreeting()}, {user?.username}!
        </h1>
        <p className="text-muted-foreground mt-2">
          Welcome to your PocketPloy dashboard. Manage your PocketBase instances with ease.
        </p>
      </div>

      {/* Stats Cards */}
      {/* <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Instances
            </CardTitle>
            <Database className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {loading ? "..." : instanceCount}
            </div>
            <p className="text-xs text-muted-foreground">
              Active PocketBase instances
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Available Slots
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {loading ? "..." : 5 - instanceCount}
            </div>
            <p className="text-xs text-muted-foreground">
              Out of 5 maximum instances
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Quick Actions
            </CardTitle>
            <Plus className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <Button
              onClick={() => router.push("/dashboard/create")}
              className="w-full mt-2"
              size="sm"
            >
              Create Instance
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Account
            </CardTitle>
            <User className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <Button
              onClick={() => router.push("/dashboard/profile")}
              variant="outline"
              className="w-full mt-2"
              size="sm"
            >
              View Profile
            </Button>
          </CardContent>
        </Card>
      </div> */}

      {/* Quick Links */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => router.push("/dashboard/instances")}>
          <CardHeader>
            <Database className="h-8 w-8 mb-2 text-blue-600" />
            <CardTitle>View All Instances</CardTitle>
            <CardDescription>
              Manage and monitor your PocketBase instances
            </CardDescription>
          </CardHeader>
        </Card>

        <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => router.push("/dashboard/archived")}>
          <CardHeader>
            <Archive className="h-8 w-8 mb-2 text-gray-600" />
            <CardTitle>Archived Instances</CardTitle>
            <CardDescription>
              View your deleted instances and restore data
            </CardDescription>
          </CardHeader>
        </Card>
      </div>

      {/* Getting Started */}
      {/* {instanceCount === 0 && !loading && (
        <Card className="border-dashed border-2">
          <CardHeader>
            <CardTitle>Get Started</CardTitle>
            <CardDescription>
              You haven't created any instances yet. Let's get you started!
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <h4 className="font-medium">Steps to create your first instance:</h4>
              <ol className="list-decimal list-inside space-y-1 text-sm text-muted-foreground">
                <li>Click on "Create Instance" or navigate to the Create page</li>
                <li>Enter a name for your instance</li>
                <li>Set up admin credentials (email and password)</li>
                <li>Click create and wait for the instance to be ready</li>
                <li>Access your PocketBase admin panel via the generated URL</li>
              </ol>
            </div>
            <Button onClick={() => router.push("/dashboard/create")} className="w-full">
              <Plus className="h-4 w-4 mr-2" />
              Create Your First Instance
            </Button>
          </CardContent>
        </Card>
      )} */}
    </div>
  );
}
