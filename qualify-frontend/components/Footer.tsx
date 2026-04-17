import { Terminal, Share2, Mail } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t border-white/10 bg-bg-dark px-6 md:px-20 py-16">
      <div className="max-w-7xl mx-auto grid md:grid-cols-2 gap-12">
        <div>
          <div className="flex items-center gap-2 mb-6">
            <Terminal className="w-6 h-6 text-accent" />
            <h2 className="text-xl font-bold">Qualify</h2>
          </div>
          <p className="text-sm text-neutral-slate leading-relaxed">
            Conectando empresas a testadores de software qualificados globalmente
          </p>
        </div>
      </div>
    </footer>
  );
}
