"use client";

import { Terminal, Search } from "lucide-react";
import Link from "next/link";

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-white/10 bg-bg-dark/80 backdrop-blur-md px-6 md:px-20 py-4">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <div className="flex items-center gap-10">
          <Link href="/" className="flex items-center gap-3">
            <Terminal className="text-accent w-8 h-8" />
            <h2 className="text-xl font-bold tracking-tight">Qualify</h2>
          </Link>
          <nav className="hidden lg:flex items-center gap-8">
            <a
              className="text-sm font-medium hover:text-accent transition-colors"
              href="#tiers"
            >
              Tiers de testes
            </a>
            <a
              className="text-sm font-medium hover:text-accent transition-colors"
              href="#devices"
            >
              Dispositivos
            </a>
            <a
              className="text-sm font-medium hover:text-accent transition-colors"
              href="#guarantee"
            >
              Garantias
            </a>
          </nav>
        </div>
        <div className="flex items-center gap-6">
          <div className="hidden md:flex items-center bg-white/5 rounded px-3 py-1.5 border border-white/10">
            <Search className="text-neutral-slate w-5 h-5" />
            <input
              className="bg-transparent border-none focus:ring-0 focus:outline-none text-sm w-50 placeholder:text-neutral-slate ml-2"
              placeholder="Buscar competências"
              type="text"
            />
          </div>
          <button className="bg-primary hover:bg-primary-light text-white px-5 py-2 text-sm font-bold transition-all rounded cursor-pointer">
            Cadastrar
          </button>
        </div>
      </div>
    </header>
  );
}
