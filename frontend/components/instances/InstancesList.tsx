"use client";

import { Instance } from "@/types/instance";
import { InstanceCard } from "./InstanceCard";
import { Loader2, Database } from "lucide-react";
import {
  Frame,
  FrameDescription,
  FrameFooter,
  FrameHeader,
  FramePanel,
  FrameTitle,
} from "@/components/coss-ui/frame"
interface InstancesListProps {
  instances: Instance[];
  loading: boolean;
  onDelete: (id: string) => void;
}

export function InstancesList({ instances, loading, onDelete }: InstancesListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (instances.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <div className="p-4 bg-gray-100 rounded-full mb-4">
          <Database className="h-8 w-8 text-gray-400" />
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">
          No instances yet
        </h3>
        <p className="text-gray-600 mb-6 max-w-md">
          Create your first PocketBase instance to get started. Each instance is isolated and accessible via a unique subdomain.
        </p>
      </div>
    );
  }

  return (
    <Frame>
      <FrameHeader>
    <FrameTitle>Title</FrameTitle>
    <FrameDescription>Description</FrameDescription>
  </FrameHeader>
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {instances.map((instance) => (
        <InstanceCard
          key={instance.id}
          instance={instance}
          onDelete={onDelete}
        />
      ))}
    </div>
    </Frame>
  );
}
