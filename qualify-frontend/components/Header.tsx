"use client";

import { useState, useEffect, useRef } from "react";
import {
  Terminal,
  Search,
  User,
  LogOut,
  Settings,
  ChevronDown,
  ChevronsUp,
} from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { getSessionUser, clearSession } from "@/libs/session";
import { analystService } from "@/libs/services";
import type { User as SessionUser } from "@/types/services/user";

export function Header() {
  const router = useRouter();
  const [user, setUser] = useState<SessionUser | null>(null);
  const [hasAnalyst, setHasAnalyst] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    getSessionUser().then((sessionUser) => {
      setUser(sessionUser);
      if (sessionUser) {
        analystService
          .getByUserId(sessionUser.id)
          .then((analyst) => setHasAnalyst(analyst != null));
      } else {
        setHasAnalyst(false);
      }
    });
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

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    const term = searchTerm.trim();
    router.push(
      term
        ? `/customer/searchAnalyst?skill=${encodeURIComponent(term)}`
        : "/customer/searchAnalyst",
    );
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
            <Link
              className="text-sm font-medium hover:text-accent transition-colors"
              href="/#tiers"
            >
              Tiers de testes
            </Link>
            <Link
              className="text-sm font-medium hover:text-accent transition-colors"
              href="/#devices"
            >
              Dispositivos
            </Link>
            <Link
              className="text-sm font-medium hover:text-accent transition-colors"
              href="/#guarantee"
            >
              Garantias
            </Link>
          </nav>
        </div>
        <div className="flex items-center gap-6">
          <form
            onSubmit={handleSearch}
            className="hidden md:flex items-center bg-white/5 rounded px-3 py-1.5 border border-white/10 focus-within:border-accent/50 transition-colors"
          >
            <button
              type="submit"
              aria-label="Buscar competências"
              className="text-neutral-slate hover:text-accent transition-colors cursor-pointer"
            >
              <Search className="w-5 h-5" />
            </button>
            <input
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="bg-transparent border-none focus:ring-0 focus:outline-none text-sm w-50 placeholder:text-neutral-slate ml-2"
              placeholder="Buscar competências"
              type="text"
            />
          </form>

          {user ? (
            <div className="relative" ref={menuRef}>
              <button
                onClick={() => setMenuOpen(!menuOpen)}
                className="flex items-center gap-2 bg-white/5 border border-white/10 hover:border-accent/50 px-4 py-2 rounded-lg transition-all cursor-pointer"
              >
                <User className="w-4 h-4 text-accent" />
                <span className="text-sm font-medium truncate max-w-30">
                  {user.name.split(" ")[0]}
                </span>
                <ChevronDown
                  className={`w-3 h-3 text-neutral-slate transition-transform ${menuOpen ? "rotate-180" : ""}`}
                />
              </button>

              {menuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-bg-dark border border-white/10 rounded-lg shadow-lg overflow-hidden">
                  <div className="px-4 py-3 border-b border-white/10">
                    <p className="text-sm font-medium text-white truncate">
                      {user.name}
                    </p>
                    <p className="text-xs text-neutral-slate truncate">
                      {user.email}
                    </p>
                  </div>
                  <Link
                    href="/user/edit"
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2 px-4 py-2.5 text-sm text-white/80 hover:bg-white/5 hover:text-accent transition-colors"
                  >
                    <Settings className="w-4 h-4" />
                    Minha conta
                  </Link>
                  <Link
                    href="/analyst/profile"
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2 px-4 py-2.5 text-sm text-white/80 hover:bg-white/5 hover:text-accent transition-colors"
                  >
                    <User className="w-4 h-4" />
                    Perfil
                  </Link>
                  {!hasAnalyst && (
                    <Link
                      href="/analyst/register"
                      onClick={() => setMenuOpen(false)}
                      className="flex items-center gap-2 px-4 py-2.5 text-sm text-white/80 hover:bg-white/5 hover:text-accent transition-colors"
                    >
                      <ChevronsUp className="w-4 h-4" />
                      Cadastrar como analista
                    </Link>
                  )}
                  {hasAnalyst && (
                    <Link
                      href="/analyst/edit"
                      onClick={() => setMenuOpen(false)}
                      className="flex items-center gap-2 px-4 py-2.5 text-sm text-white/80 hover:bg-white/5 hover:text-accent transition-colors"
                    >
                      <Settings className="w-4 h-4" />
                      Editar analista
                    </Link>
                  )}
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
