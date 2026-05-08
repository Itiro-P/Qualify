import React from "react";

type ButtonVariant = "primary" | "danger" | "outline";

interface FormButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  loading?: boolean;
  loadingText?: string;
  fullWidth?: boolean;
}

const variantClasses: Record<ButtonVariant, string> = {
  primary: "bg-primary hover:bg-primary-light text-white",
  danger: "bg-red-600 hover:bg-red-700 text-white",
  outline: "border border-white/10 text-white hover:bg-white/5",
};

export function FormButton({
  variant = "primary",
  loading = false,
  loadingText,
  fullWidth = true,
  children,
  disabled,
  className = "",
  ...buttonProps
}: FormButtonProps) {
  return (
    <button
      disabled={disabled || loading}
      className={`font-bold py-3 rounded-lg transition-all cursor-pointer disabled:opacity-50 ${
        variantClasses[variant]
      } ${fullWidth ? "w-full" : ""} ${className}`}
      {...buttonProps}
    >
      {loading ? (loadingText ?? "Carregando...") : children}
    </button>
  );
}
