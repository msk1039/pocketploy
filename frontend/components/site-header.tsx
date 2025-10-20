"use client"

import { usePathname } from "next/navigation"
import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { RefreshCw } from "lucide-react"
import { Button } from "@/components/ui/button"

export function SiteHeader() {
  const pathname = usePathname()

  const getPageTitle = () => {
    if (pathname === "/dashboard") return "Overview"
    if (pathname === "/dashboard/instances") return "Instances"
    if (pathname === "/dashboard/archived") return "Archived Instances"
    if (pathname === "/dashboard/create") return "Create Instance"
    if (pathname === "/dashboard/profile") return "Profile"
    if (pathname?.startsWith("/dashboard/instances/")) return "Instance Details"
    return "Dashboard"
  }

  const handleRefresh = () => {
    window.location.reload()
  }

  return (
    <header className="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mx-2 data-[orientation=vertical]:h-4"
        />
        <h1 className="text-base font-medium">{getPageTitle()}</h1>
        <div className="ml-auto flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={handleRefresh}
            className="flex items-center gap-2"
          >
            <RefreshCw className="h-4 w-4" />
            <span className="hidden sm:inline">Refresh</span>
          </Button>
        </div>
      </div>
    </header>
  )
}
