"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Instance } from "@/types/instance";
import { getInstance, startInstance, stopInstance, deleteInstance } from "@/lib/api";
import { InstanceLogs } from "./InstanceLogs";
import { InstanceSwitch } from "./instanceSwitch";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { ArrowLeft, Copy, ExternalLink, Loader2, AlertTriangle, Globe } from "lucide-react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface InstanceDetailsProps {
  instanceId: string;
}

export function InstanceDetails({ instanceId }: InstanceDetailsProps) {
  const router = useRouter();
  const [instance, setInstance] = useState<Instance | null>(null);
  const [loading, setLoading] = useState(true);
  const [switchLoading, setSwitchLoading] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deletingInstance, setDeletingInstance] = useState(false);

  const fetchInstance = async () => {
    setLoading(true);
    try {
      const response = await getInstance(instanceId);
      setInstance(response.instance);
    } catch (error: any) {
      console.error("Failed to fetch instance:", error);
      toast.error("Failed to load instance", {
        description: error.message || "An error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInstance();
  }, [instanceId]);

  const handleBack = () => {
    router.push("/dashboard/instances");
  };

  const handleSwitchChange = async (checked: boolean) => {
    if (!instance) return;
    
    setSwitchLoading(true);
    
    const action = checked ? "start" : "stop";
    const actionText = checked ? "Starting" : "Stopping";
    
    toast.promise(
      async () => {
        if (checked) {
          await startInstance(instance.id);
        } else {
          await stopInstance(instance.id);
        }
        await fetchInstance();
      },
      {
        loading: `${actionText} instance...`,
        success: `Instance ${action}ed successfully!`,
        error: (error: any) => error.message || `Failed to ${action} instance`,
        finally: () => setSwitchLoading(false),
      }
    );
  };

  const handleCopyUrl = () => {
    if (!instance) return;
    const url = `http://${instance.subdomain}/_/`;
    navigator.clipboard.writeText(url);
    toast.success("URL copied to clipboard!");
  };

  const handleOpenInstance = () => {
    if (!instance) return;
    const url = `http://${instance.subdomain}/_/`;
    window.open(url, "_blank");
  };

  const handleDeleteClick = () => {
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!instance) return;
    
    setDeletingInstance(true);
    try {
      await deleteInstance(instance.id);
      toast.success("Instance deleted successfully");
      router.push("/dashboard/instances");
    } catch (error: any) {
      toast.error("Failed to delete instance", {
        description: error.message || "An error occurred",
      });
      setDeletingInstance(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

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

  if (loading) {
    return (
      <div className="space-y-6">
        {/* Back Button Skeleton */}
        <Skeleton className="h-10 w-24" />

        {/* Instance Name and Controls Card Skeleton */}
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <Skeleton className="h-8 w-48" />
                <Skeleton className="h-6 w-20" />
              </div>
              <div className="flex items-center gap-3">
                <Skeleton className="h-4 w-32" />
                <Skeleton className="h-5 w-11 rounded-full" />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Instance URL Card Skeleton */}
        <Card>
          <CardContent className="pt-6">
            <div className="space-y-2">
              <Skeleton className="h-4 w-24" />
              <div className="flex items-center gap-2">
                <Skeleton className="flex-1 h-10" />
                <Skeleton className="h-10 w-10" />
                <Skeleton className="h-10 w-10" />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Tabs Skeleton */}
        <div className="space-y-4">
          <Skeleton className="h-10 w-full" />
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  if (!instance) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Instance not found</h2>
        <p className="text-gray-600 mb-6">The instance you're looking for doesn't exist.</p>
        <Button onClick={handleBack}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Dashboard
        </Button>
      </div>
    );
  }

  const isCreating = instance.status === "creating";
  const isRunning = instance.status === "running";

  return (
    <div className="space-y-6">
      {/* Back Button */}
      <Button variant="outline" onClick={handleBack} className="flex items-center gap-2">
        <ArrowLeft className="h-4 w-4" />
        Back
      </Button>

      {/* Instance Name and Controls Card */}
      <Card>
        <CardContent className="">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <h1 className="text-2xl font-bold text-gray-900">{instance.name}</h1>
              <Badge className={getStatusColor(instance.status)} variant="outline">
                {isCreating && <Loader2 className="h-3 w-3 mr-1 animate-spin" />}
                {instance.status}
              </Badge>
            </div>
            <div className="flex items-center gap-3">
              <span className="text-sm text-gray-600">
                {isRunning ? "Instance is running" : "Instance is stopped"}
              </span>
              <InstanceSwitch
                checked={isRunning}
                onCheckedChange={handleSwitchChange}
                disabled={isCreating}
                loading={switchLoading}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Instance URL Card */}
      <Card className="border-2 border-dashed border-blue-200 bg-gradient-to-br from-blue-50/50 to-transparent">
        <CardContent className="">
          <div className="space-y-3">
            <div className="flex w-full justify-between">
            <div className="flex items-center gap-2 text-sm font-semibold text-blue-900">
              <Globe className="h-4 w-4" />
              Instance URL
            </div>
             {!isRunning && (
              <p className="text-xs text-amber-600 flex items-center gap-1.5 bg-amber-50 px-3 rounded-md border border-amber-200">
                <AlertTriangle className="h-3 w-3" />
                Instance must be running to access the URL
              </p>
            )}
            </div>
            <div className="flex items-center gap-2">
              <div className="flex-1 relative group">
                <div className="absolute inset-0 bg-gradient-to-r from-blue-500/10 to-purple-500/10 rounded-lg blur-sm group-hover:blur-md transition-all"></div>
                <code className="relative block bg-white px-4 py-3 rounded-lg text-sm break-all font-mono border-2 border-blue-100 shadow-sm hover:border-blue-300 transition-colors">
                  <span className="text-blue-600">http://</span>
                  <span className="text-gray-900 font-semibold">{instance.subdomain}</span>
                  <span className="text-blue-600">/_/</span>
                </code>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleCopyUrl}
                className="flex-shrink-0 h-11 border-2 hover:border-blue-300 hover:bg-blue-50"
                title="Copy URL"
              >
                <Copy className="h-4 w-4" />
              </Button>
              <Button
                size="sm"
                onClick={handleOpenInstance}
                disabled={!isRunning}
                className="flex-shrink-0 h-11 bg-blue-600 hover:bg-blue-700"
                title="Open in new tab"
              >
                <ExternalLink className="h-4 w-4 mr-1" />
                Open
              </Button>
            </div>
           
          </div>
        </CardContent>
      </Card>

      {/* Tabs Section */}
      <Tabs defaultValue="metadata" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
          <TabsTrigger value="logs">Logs</TabsTrigger>
          <TabsTrigger value="graphs">Graphs</TabsTrigger>
          <TabsTrigger value="destroy">Destroy</TabsTrigger>
        </TabsList>

        {/* Metadata Tab */}
        <TabsContent value="metadata">
          <Card>
            <CardContent className="pt-6">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[200px]">Property</TableHead>
                    <TableHead>Value</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow>
                    <TableCell className="font-medium">Instance ID</TableCell>
                    <TableCell>
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded">{instance.id}</code>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Name</TableCell>
                    <TableCell>{instance.name}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Slug</TableCell>
                    <TableCell>
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded">{instance.slug}</code>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Status</TableCell>
                    <TableCell>
                      <Badge className={getStatusColor(instance.status)} variant="outline">
                        {instance.status}
                      </Badge>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Subdomain</TableCell>
                    <TableCell>
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded">{instance.subdomain}</code>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Container ID</TableCell>
                    <TableCell>
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                        {instance.container_id ? instance.container_id.substring(0, 12) : "N/A"}
                      </code>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Data Path</TableCell>
                    <TableCell>
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded break-all">
                        {instance.data_path}
                      </code>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Created At</TableCell>
                    <TableCell>{formatDate(instance.created_at)}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Updated At</TableCell>
                    <TableCell>{formatDate(instance.updated_at)}</TableCell>
                  </TableRow>
                  {instance.last_accessed_at && (
                    <TableRow>
                      <TableCell className="font-medium">Last Accessed</TableCell>
                      <TableCell>{formatDate(instance.last_accessed_at)}</TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Logs Tab */}
        <TabsContent value="logs">
          <InstanceLogs instanceId={instanceId} />
        </TabsContent>

        {/* Graphs Tab */}
        <TabsContent value="graphs">
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="p-4 bg-gray-100 rounded-full mb-4">
                  <AlertTriangle className="h-8 w-8 text-gray-400" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  Coming Soon
                </h3>
                <p className="text-gray-600 max-w-md">
                  Analytics graphs and metrics visualization will be available soon.
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Destroy Tab */}
        <TabsContent value="destroy">
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-4">
                <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
                  <p className="text-sm text-amber-900">
                    <strong>Warning:</strong> This will permanently delete the instance. 
                    This action cannot be undone. The Docker container will be stopped and removed, 
                    but your data will be preserved in the storage directory.
                  </p>
                </div>
                <Button
                  variant="destructive"
                  onClick={handleDeleteClick}
                  className="w-full"
                  disabled={isCreating}
                >
                  Delete Instance
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-red-100 rounded-lg">
                <AlertTriangle className="h-5 w-5 text-red-600" />
              </div>
              <DialogTitle>Delete Instance</DialogTitle>
            </div>
            <DialogDescription className="pt-3">
              Are you sure you want to delete <strong>{instance.name}</strong>?
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
              onClick={() => setDeleteDialogOpen(false)}
              disabled={deletingInstance}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDeleteConfirm}
              disabled={deletingInstance}
            >
              {deletingInstance ? (
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
    </div>
  );
}
