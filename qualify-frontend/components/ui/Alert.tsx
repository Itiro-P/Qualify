import React from "react";

interface AlertProps {
  children: React.ReactNode;
  variant: "error" | "success" | "warning";
}

const variantClasses = {
  error: "bg-red-500/10 border-red-500/30 text-red-400",
  success: "bg-accent/10 border-accent/30 text-accent",
  warning: "bg-yellow-500/10 border-yellow-500/30 text-yellow-400",
};

export function Alert({ children, variant }: AlertProps) {
  return (
    <div
      className={`mb-6 rounded-lg border px-4 py-3 text-sm ${variantClasses[variant]}`}
    >
      {children}
    </div>
  );
}
