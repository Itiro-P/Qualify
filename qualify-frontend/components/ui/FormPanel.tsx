import React from "react";

interface FormPanelProps {
  title: string;
  description?: string;
  children: React.ReactNode;
  maxWidth?: string;
}

export function FormPanel({
  title,
  description,
  children,
  maxWidth = "max-w-lg",
}: FormPanelProps) {
  return (
    <section className="min-h-screen flex items-center justify-center px-4 py-12">
      <div className={`w-full ${maxWidth} glass-panel rounded-2xl p-8`}>
        <h1 className="text-2xl font-bold text-white mb-2">{title}</h1>
        {description && (
          <p className="text-sm text-neutral-slate mb-8">{description}</p>
        )}
        {children}
      </div>
    </section>
  );
}
