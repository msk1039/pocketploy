"use client";

import * as React from "react";
import * as SwitchPrimitive from "@radix-ui/react-switch";
import { cn } from "@/lib/utils";
import { Loader2 } from "lucide-react";

interface InstanceSwitchProps extends React.ComponentPropsWithoutRef<typeof SwitchPrimitive.Root> {
  loading?: boolean;
}

export function InstanceSwitch({ className, loading, disabled, ...props }: InstanceSwitchProps) {
  return (
    <SwitchPrimitive.Root
      data-slot="switch"
      disabled={disabled || loading}
      className={cn(
        "peer inline-flex h-8 w-16 shrink-0 items-center rounded-full border-2 shadow-md transition-all outline-none",
        "focus-visible:ring-4 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        // When checked (on) - green background
        "data-[state=checked]:bg-green-500 data-[state=checked]:border-green-600",
        "data-[state=checked]:hover:bg-green-600",
        // When unchecked (off) - grey background
        "data-[state=unchecked]:bg-gray-300 data-[state=unchecked]:border-gray-400",
        "data-[state=unchecked]:hover:bg-gray-400",
        // Focus states
        "focus-visible:data-[state=checked]:ring-green-200",
        "focus-visible:data-[state=unchecked]:ring-gray-200",
        className
      )}
      {...props}
    >
      <SwitchPrimitive.Thumb
        data-slot="switch-thumb"
        className={cn(
          "pointer-events-none flex items-center justify-center bg-white size-7 rounded-full shadow-lg ring-0 transition-transform",
          "data-[state=checked]:translate-x-8 data-[state=unchecked]:translate-x-0"
        )}
      >
        {loading && <Loader2 className="h-3.5 w-3.5 animate-spin text-gray-600" />}
      </SwitchPrimitive.Thumb>
    </SwitchPrimitive.Root>
  );
}
