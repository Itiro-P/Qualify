"use client";

import { Alert } from "@/components/ui";
import { Loading } from "@/components/ui/Loading";
import { Analyst } from "@/types/services/analyst";
import { analystService } from "@/libs/services/analystService";
import { IFormResponse } from "@/types/customer/formResponse";
import { MapPin, Star, Clock, Mail } from "lucide-react";
import Link from "next/link";
import { useEffect, useState } from "react";

function AnalystCard({ analyst }: { analyst: Analyst }) {
  return (
    <Link
      href={`/analyst/profile?id=${analyst.id}`}
      className="group block rounded-xl border border-white/10 bg-white/5 p-6 hover:border-accent/50 transition-all"
    >
      <div className="flex items-start justify-between gap-4">
        <div>
          <h3 className="text-lg font-bold text-white group-hover:text-accent transition-colors">
            {analyst.name}
          </h3>
          <p className="flex items-center gap-1.5 text-sm text-neutral-slate mt-1">
            <MapPin className="w-4 h-4" />
            {[analyst.city, analyst.country_state, analyst.country_name]
              .filter(Boolean)
              .join(", ") || "Localização não informada"}
          </p>
        </div>
        <div className="flex items-center gap-1 rounded-full bg-accent/10 px-3 py-1 text-sm font-semibold text-accent shrink-0">
          <Star className="w-4 h-4 fill-accent" />
          {analyst.mean_rating ? analyst.mean_rating.toFixed(1) : "—"}
        </div>
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-x-6 gap-y-2 text-sm text-slate-300">
        <span className="flex items-center gap-1.5">
          <Clock className="w-4 h-4 text-accent" />
          R$ {analyst.hourly_rate}/h
        </span>
        <span className="flex items-center gap-1.5">
          <Star className="w-4 h-4 text-accent" />
          {analyst.total_reviews} avaliações
        </span>
        {analyst.email && (
          <span className="flex items-center gap-1.5 text-neutral-slate">
            <Mail className="w-4 h-4" />
            {analyst.email}
          </span>
        )}
      </div>
    </Link>
  );
}

export function SearchResult({ filters }: { filters: IFormResponse | null }) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [analysts, setAnalysts] = useState<Analyst[]>([]);
  const [searched, setSearched] = useState(false);

  useEffect(() => {
    if (!filters) return;

    let cancelled = false;

    async function run() {
      setLoading(true);
      setError("");

      try {
        const resp = await analystService.list({
          country: filters!.localization.country || undefined,
          country_state: filters!.localization.state || undefined,
          city: filters!.localization.city || undefined,
          min_hourly_rate: filters!.min_hourly_rate || undefined,
          max_hourly_rate: filters!.max_hourly_rate || undefined,
          min_rating: filters!.rating || undefined,
          skills: filters!.skills.length
            ? filters!.skills.join(",")
            : undefined,
        });

        if (!cancelled) {
          setAnalysts(resp?.analysts ?? []);
        }
      } catch {
        if (!cancelled) {
          setError("Erro ao buscar analistas. Tente novamente.");
          setAnalysts([]);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
          setSearched(true);
        }
      }
    }

    run();

    return () => {
      cancelled = true;
    };
  }, [filters]);

  if (!filters && !searched) {
    return (
      <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
        <p className="text-neutral-slate">
          Use os filtros ao lado e clique em{" "}
          <span className="text-accent font-medium">Buscar analistas</span> para
          começar.
        </p>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="flex flex-col items-center gap-4">
          <div className="relative h-12 w-12">
            <div className="absolute inset-0 animate-ping rounded-full bg-accent opacity-30" />
            <div className="absolute inset-1 animate-spin rounded-full border-4 border-accent border-t-transparent" />
          </div>
          <p className="text-sm text-neutral-slate">Buscando analistas...</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      {error && <Alert variant="error">{error}</Alert>}

      {!error && analysts.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
          <p className="text-neutral-slate">
            Nenhum analista encontrado com esses filtros.
          </p>
        </div>
      ) : (
        <>
          {analysts.length > 0 && (
            <p className="text-sm text-neutral-slate mb-4">
              {analysts.length}{" "}
              {analysts.length === 1
                ? "analista encontrado"
                : "analistas encontrados"}
            </p>
          )}
          <div className="grid sm:grid-cols-2 gap-4">
            {analysts.map((analyst) => (
              <AnalystCard key={analyst.id} analyst={analyst} />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
