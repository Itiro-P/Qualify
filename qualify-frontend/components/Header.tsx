"use client";

import { useState, useEffect, useRef } from "react";
import { Terminal, Search, User, LogOut, Settings, ChevronDown } from "lucide-react";
import Link from "next/link";
import { getSessionUser, clearSession } from "@/libs/session";
import type { SessionUser } from "@/libs/session";

export function Header() {
  const [user, setUser] = useState<SessionUser | null>(null);
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setUser(getSessionUser());
  }, []);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  function handleLogout() {
    clearSession();
    setUser(null);
    setMenuOpen(false);
    window.location.href = "/";
  }

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

          {user ? (
            <div className="relative" ref={menuRef}>
              <button
                onClick={() => setMenuOpen(!menuOpen)}
                className="flex items-center gap-2 bg-white/5 border border-white/10 hover:border-accent/50 px-4 py-2 rounded-lg transition-all cursor-pointer"
              >
                <User className="w-4 h-4 text-accent" />
                <span className="text-sm font-medium truncate max-w-[120px]">
                  {user.name.split(" ")[0]}
                </span>
                <ChevronDown className={`w-3 h-3 text-neutral-slate transition-transform ${menuOpen ? "rotate-180" : ""}`} />
              </button>

              {menuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-bg-dark border border-white/10 rounded-lg shadow-lg overflow-hidden">
                  <div className="px-4 py-3 border-b border-white/10">
                    <p className="text-sm font-medium text-white truncate">{user.name}</p>
                    <p className="text-xs text-neutral-slate truncate">{user.email}</p>
                  </div>
                  <Link
                    href="/user/edit"
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2 px-4 py-2.5 text-sm text-white/80 hover:bg-white/5 hover:text-accent transition-colors"
                  >
                    <Settings className="w-4 h-4" />
                    Minha conta
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="flex items-center gap-2 w-full px-4 py-2.5 text-sm text-red-400 hover:bg-red-500/10 transition-colors cursor-pointer"
                  >
                    <LogOut className="w-4 h-4" />
                    Sair
                  </button>
                </div>
              )}
            </div>
          ) : (
            <div className="flex items-center gap-3">
              <Link
                href="/user/login"
                className="text-sm font-medium text-white/80 hover:text-accent transition-colors"
              >
                Entrar
              </Link>
              <Link
                href="/user/register"
                className="bg-primary hover:bg-primary-light text-white px-5 py-2 text-sm font-bold transition-all rounded cursor-pointer"
              >
                Cadastrar
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
