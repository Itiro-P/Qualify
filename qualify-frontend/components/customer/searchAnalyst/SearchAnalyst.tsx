"use client";

import { SearchForm, SearchResult } from "@/components/customer/searchAnalyst";
import { IFormResponse } from "@/types/customer/formResponse";
import { useSearchParams } from "next/navigation";
import { useMemo, useState } from "react";

function emptyFilters(): IFormResponse {
  return {
    min_hourly_rate: 0,
    max_hourly_rate: 0,
    rating: 0,
    skills: [],
    localization: {
      country: "",
      state: "",
      city: "",
    },
  };
}

export function SearchAnalyst() {
  const searchParams = useSearchParams();
  const initialSkill = searchParams.get("skill")?.trim() ?? "";
  const initialCountry = searchParams.get("country")?.trim() ?? "";

  const initialFilters = useMemo<IFormResponse>(() => {
    const base = emptyFilters();
    if (initialSkill) {
      base.skills = [initialSkill];
    }
    if (initialCountry) {
      base.localization.country = initialCountry;
    }
    return base;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const [formResponse, setFormResponse] =
    useState<IFormResponse>(initialFilters);
  // Busca todos os analistas por padrão; os filtros da URL apenas refinam.
  const [appliedFilters, setAppliedFilters] =
    useState<IFormResponse | null>(initialFilters);

  function handleSearch() {
    setAppliedFilters({
      ...formResponse,
      localization: { ...formResponse.localization },
      skills: [...formResponse.skills],
    });
  }

  return (
    <section className="px-6 md:px-20 py-14">
      <div className="max-w-7xl mx-auto">
        <div className="mb-10">
          <h1 className="text-3xl font-bold mb-4">Encontrar analistas</h1>
          <div className="h-1 w-20 bg-accent mb-6" />
          <p className="text-neutral-slate max-w-2xl">
            Filtre por competências, avaliação, faixa de preço e localização
            para encontrar o analista de testes ideal para o seu projeto.
          </p>
        </div>

        <div className="grid lg:grid-cols-[320px_1fr] gap-8 items-start">
          <SearchForm
            formResponse={formResponse}
            setFormResponse={setFormResponse}
            onSearch={handleSearch}
          />
          <SearchResult filters={appliedFilters} />
        </div>
      </div>
    </section>
  );
}
