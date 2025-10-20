"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Archive, Clock, Database } from "lucide-react";

export default function ArchivedPage() {
  return (
    <div className="space-y-6">
      <h1>COMMING SOON</h1>
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Archived Instances</h2>
        <p className="text-muted-foreground mt-1">
          View and manage your deleted instances
        </p>
      </div>

      <Card className="border-dashed border-2">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="p-3 bg-gray-100 rounded-lg">
              <Archive className="h-6 w-6 text-gray-600" />
            </div>
            <div>
              <CardTitle>No Archived Instances</CardTitle>
              <CardDescription>
                Your deleted instances will appear here
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <Clock className="h-5 w-5 text-blue-600 mt-0.5" />
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-blue-900">Data Retention Policy</h4>
                <p className="text-sm text-blue-700">
                  When you delete an instance, we keep your data for 30 days. This allows you to restore 
                  your instance if needed. After 30 days, the data is permanently deleted.
                </p>
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <h4 className="font-medium text-sm">What gets archived:</h4>
            <ul className="space-y-1 text-sm text-muted-foreground">
              <li className="flex items-center gap-2">
                <Database className="h-4 w-4" />
                <span>Instance metadata (name, subdomain, status)</span>
              </li>
              <li className="flex items-center gap-2">
                <Database className="h-4 w-4" />
                <span>Database files and collections</span>
              </li>
              <li className="flex items-center gap-2">
                <Database className="h-4 w-4" />
                <span>Uploaded files and assets</span>
              </li>
              <li className="flex items-center gap-2">
                <Clock className="h-4 w-4" />
                <span>Deletion timestamp and user information</span>
              </li>
            </ul>
          </div>

          <div className="pt-4 border-t">
            <p className="text-xs text-muted-foreground">
              Note: The restore feature is coming soon. For now, archived instances are kept for audit purposes 
              and potential future recovery.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
