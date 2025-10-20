"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { type Icon } from "@tabler/icons-react"

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
} from "@/components/ui/sidebar"

export function NavClouds({
  items,
}: {
  items: {
    title: string
    url: string
    icon?: Icon
    isActive?: boolean
    items?: {
      title: string
      url: string
    }[]
  }[]
}) {
  const pathname = usePathname()

  return (
    <SidebarGroup>
      <SidebarGroupLabel className="text-xs font-semibold">Quick Actions</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => (
          item.items?.map((subItem) => {
            const isActive = pathname === subItem.url
            return (
              <SidebarMenuItem key={subItem.title}>
                <SidebarMenuButton 
                  asChild 
                  isActive={isActive}
                  size="lg"
                  className="text-base font-sans font-medium"
                >
                  <Link href={subItem.url}>
                    {item.icon && <item.icon className="!size-5" />}
                    <span>{subItem.title}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          })
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}
