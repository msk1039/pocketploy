"use client"

import * as React from "react"
import Link from "next/link"
import {
  IconCirclePlusFilled,
  IconDashboard,
  IconDatabase,
  IconFolder,
  IconInnerShadowTop,
  IconSettings,
  IconUserCircle,
} from "@tabler/icons-react"

import { NavMain } from "@/components/nav-main"
import { NavSecondary } from "@/components/nav-secondary"
import { NavUser } from "@/components/nav-user"
import { NavClouds } from "@/components/nav-clouds"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { useAuth } from "@/contexts/AuthContext"

const data = {
  navMain: [
    {
      title: "Overview",
      url: "/dashboard",
      icon: IconDashboard,
    },
    {
      title: "Instances",
      url: "/dashboard/instances",
      icon: IconDatabase,
    },
    {
      title: "Archived",
      url: "/dashboard/archived",
      icon: IconFolder,
    },
  ],
  navCreate: [
    {
      title: "Create",
      icon: IconCirclePlusFilled,
      isActive: true,
      url: "#",
      items: [
        {
          title: "Create Instance",
          url: "/dashboard/create",
        },
      ],
    },
  ],
  navSecondary: [
    {
      title: "Profile",
      url: "/dashboard/profile",
      icon: IconUserCircle,
    },
    {
      title: "Settings",
      url: "#",
      icon: IconSettings,
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth()

  const userData = user ? {
    name: user.username,
    email: user.email,
    avatar: "/avatars/default.jpg",
  } : {
    name: "Guest",
    email: "guest@example.com",
    avatar: "/avatars/default.jpg",
  }

  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              size="lg"
              className="data-[slot=sidebar-menu-button]:!p-2 font-sans"
            >
              <Link href="/dashboard">
                <IconDatabase className="!size-6" />
                <span className="text-lg font-semibold">PocketPloy</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavClouds items={data.navCreate} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={userData} />
      </SidebarFooter>
    </Sidebar>
  )
}
