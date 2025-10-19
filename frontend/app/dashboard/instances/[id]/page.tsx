"use client";

import { useParams } from "next/navigation";
import { InstanceDetails } from "@/components/instances/InstanceDetails";

export default function InstancePage() {
  const params = useParams();
  const instanceId = params.id as string;

  return (
    <div className="space-y-6">
      <InstanceDetails instanceId={instanceId} />
    </div>
  );
}
