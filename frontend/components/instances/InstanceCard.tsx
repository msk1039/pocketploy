"use client";

import { Instance } from "@/types/instance";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ExternalLink, MoreVertical, Trash2, Copy, Database, Loader2, Eye } from "lucide-react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import {
  Frame,
  FrameDescription,
  FrameFooter,
  FrameHeader,
  FramePanel,
  FrameTitle,
} from "@/components/coss-ui/frame"
interface InstanceCardProps {
  instance: Instance;
  onDelete: (id: string) => void;
}

export function InstanceCard({ instance, onDelete }: InstanceCardProps) {
  const router = useRouter();

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

  const isCreating = instance.status === "creating";

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleCopyUrl = () => {
    const url = `http://${instance.subdomain}/_`;
    navigator.clipboard.writeText(url);
    toast.success("URL copied to clipboard!");
  };

  const handleOpenInstance = () => {
    const url = `http://${instance.subdomain}/_/`;
    window.open(url, "_blank");
  };

  const handleViewDetails = () => {
    router.push(`/dashboard/instances/${instance.id}`);
  };

  return (
    <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={handleViewDetails}>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3 flex-1">
            <div className="p-2 bg-blue-50 rounded-lg">
              <Database className="h-5 w-5 text-blue-600" />
            </div>
            <div className="flex-1 min-w-0">
              <CardTitle className="text-lg truncate">{instance.name}</CardTitle>
              <CardDescription className="mt-1">
                <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                  {instance.slug}
                </code>
              </CardDescription>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Badge className={getStatusColor(instance.status)} variant="outline">
              {isCreating && (
                <Loader2 className="h-3 w-3 mr-1 animate-spin" />
              )}
              {instance.status}
            </Badge>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button 
                  variant="ghost" 
                  size="sm" 
                  className="h-8 w-8 p-0"
                  onClick={(e) => e.stopPropagation()}
                >
                  <MoreVertical className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" onClick={(e) => e.stopPropagation()}>
                <DropdownMenuItem onClick={handleViewDetails}>
                  <Eye className="mr-2 h-4 w-4" />
                  View Details
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleOpenInstance} disabled={instance.status !== "running"}>
                  <ExternalLink className="mr-2 h-4 w-4" />
                  Open Instance
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleCopyUrl}>
                  <Copy className="mr-2 h-4 w-4" />
                  Copy URL
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => onDelete(instance.id)}
                  className="text-red-600 focus:text-red-600"
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete Instance
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">URL:</span>
            <a
              href={`http://${instance.subdomain}/_/`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:underline flex items-center gap-1 truncate max-w-xs"
            >
              {instance.subdomain}/_/
              <ExternalLink className="h-3 w-3 flex-shrink-0" />
            </a>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">Created:</span>
            <span className="text-gray-700">{formatDate(instance.created_at)}</span>
          </div>
          {instance.last_accessed_at && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-500">Last Accessed:</span>
              <span className="text-gray-700">{formatDate(instance.last_accessed_at)}</span>
            </div>
          )}
          {instance.container_id && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-500">Container:</span>
              <code className="text-xs bg-gray-100 px-2 py-1 rounded truncate max-w-xs">
                {instance.container_id.substring(0, 12)}
              </code>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
