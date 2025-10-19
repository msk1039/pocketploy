"use client";

import { Instance } from "@/types/instance";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Database, Globe, Folder, Calendar, Clock, Loader2, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

interface InstanceMetadataProps {
  instance: Instance;
}

export function InstanceMetadata({ instance }: InstanceMetadataProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "running":
        return "bg-green-100 text-green-800 border-green-200";
      case "creating":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "stopped":
        return "bg-gray-100 text-gray-800 border-gray-200";
      case "failed":
        return "bg-red-100 text-red-800 border-red-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const isCreating = instance.status === "creating";

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleCopyUrl = () => {
    const url = `http://${instance.subdomain}`;
    navigator.clipboard.writeText(url);
    toast.success("URL copied to clipboard!");
  };

  const handleOpenInstance = () => {
    const url = `http://${instance.subdomain}`;
    window.open(url, "_blank");
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-xl">Instance Metadata</CardTitle>
          <Badge className={getStatusColor(instance.status)} variant="outline">
            {isCreating && <Loader2 className="h-3 w-3 mr-1 animate-spin" />}
            {instance.status}
          </Badge>
        </div>
        <CardDescription>View instance configuration and status</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* URL Section */}
        <div className="space-y-2">
          <div className="flex items-center text-sm font-medium text-gray-700">
            <Globe className="h-4 w-4 mr-2" />
            Instance URL
          </div>
          <div className="flex items-center gap-2">
            <code className="flex-1 bg-gray-100 px-3 py-2 rounded text-sm break-all">
              http://{instance.subdomain}
            </code>
            <Button
              variant="outline"
              size="sm"
              onClick={handleCopyUrl}
              className="flex-shrink-0"
            >
              Copy
            </Button>
            <Button
              variant="default"
              size="sm"
              onClick={handleOpenInstance}
              disabled={instance.status !== "running"}
              className="flex-shrink-0"
            >
              <ExternalLink className="h-4 w-4 mr-1" />
              Open
            </Button>
          </div>
        </div>

        {/* Instance Details */}
        <div className="grid grid-cols-2 gap-4 pt-4 border-t">
          <div className="space-y-1">
            <div className="flex items-center text-sm font-medium text-gray-700">
              <Database className="h-4 w-4 mr-2" />
              Slug
            </div>
            <code className="text-sm bg-gray-100 px-2 py-1 rounded">{instance.slug}</code>
          </div>

          <div className="space-y-1">
            <div className="flex items-center text-sm font-medium text-gray-700">
              <Folder className="h-4 w-4 mr-2" />
              Container ID
            </div>
            <code className="text-xs bg-gray-100 px-2 py-1 rounded break-all">
              {instance.container_id ? instance.container_id.substring(0, 12) : "N/A"}
            </code>
          </div>

          <div className="space-y-1">
            <div className="flex items-center text-sm font-medium text-gray-700">
              <Calendar className="h-4 w-4 mr-2" />
              Created
            </div>
            <p className="text-sm text-gray-600">{formatDate(instance.created_at)}</p>
          </div>

          <div className="space-y-1">
            <div className="flex items-center text-sm font-medium text-gray-700">
              <Clock className="h-4 w-4 mr-2" />
              Last Updated
            </div>
            <p className="text-sm text-gray-600">{formatDate(instance.updated_at)}</p>
          </div>
        </div>

        {/* Data Path */}
        <div className="space-y-2 pt-4 border-t">
          <div className="flex items-center text-sm font-medium text-gray-700">
            <Folder className="h-4 w-4 mr-2" />
            Data Path
          </div>
          <code className="block bg-gray-100 px-3 py-2 rounded text-xs break-all">
            {instance.data_path}
          </code>
        </div>
      </CardContent>
    </Card>
  );
}
