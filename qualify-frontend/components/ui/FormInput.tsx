import React from "react";

interface FormInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
  hint?: string;
  required?: boolean;
}

export function FormInput({
  label,
  error,
  hint,
  required,
  id,
  ...inputProps
}: FormInputProps) {
  const inputId = id ?? inputProps.name;

  return (
    <div>
      <label
        htmlFor={inputId}
        className="block text-sm font-medium text-white/80 mb-1"
      >
        {label}
        {required && " *"}
      </label>
      <input
        id={inputId}
        {...inputProps}
        className={`w-full rounded-lg border bg-white/5 px-4 py-2.5 text-sm text-white placeholder:text-neutral-slate focus:outline-none focus:ring-1 transition-colors ${
          error
            ? "border-red-500/50 focus:border-red-500 focus:ring-red-500"
            : "border-white/10 focus:border-accent focus:ring-accent"
        }`}
      />
      {error && <p className="text-red-400 text-xs mt-1">{error}</p>}
      {hint && !error && (
        <p className="text-neutral-slate text-xs mt-1">{hint}</p>
      )}
    </div>
  );
}
