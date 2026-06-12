"use client";

import { Alert } from "@/components/ui";
import { Loading } from "@/components/ui/Loading";
import { Analyst } from "@/types/services/analyst";
import { analystService } from "@/libs/services/analystService";
import { IFormResponse } from "@/types/customer/formResponse";
import { useState } from "react";

function AnalystCard({ analyst }: { analyst: Analyst }) {
  return (
    <div className="border p-4 rounded-md mb-4">
      <h3 className="text-lg font-semibold">{analyst.name}</h3>
      <p>Email: {analyst.email}</p>
      <p>Telefone: {analyst.phone}</p>
      <p>
        Localização: {analyst.city}, {analyst.country_name}
      </p>
      <p>Fuso horário: {analyst.timezone}</p>
      <p>Valor por hora: ${analyst.hourly_rate}</p>
      <p>
        Avaliação média: {analyst.mean_rating} ({analyst.total_reviews}{" "}
        avaliações)
      </p>
    </div>
  );
}

export function SearchResult({
  formResponse,
}: {
  formResponse: IFormResponse | null;
}) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [analysts, setAnalysts] = useState<Analyst[]>([]);

  useState(async () => {
    if (!formResponse) return;

    setLoading(true);
    setError("");
    setAnalysts([]);

    // Pegar os analistas recomendados com base na resposta do formulário
    const analystOnDataBase = await analystService
      .list({
        country: formResponse.localization.country,
        country_state: formResponse.localization.state,
        city: formResponse.localization.city,
        min_hourly_rate: formResponse.min_hourly_rate,
        max_hourly_rate: formResponse.max_hourly_rate,
        min_rating: formResponse.rating,
        skills: formResponse.skills ? formResponse.skills.join(",") : undefined,
      })
      .then(
        (resp) => {
          return resp.analysts;
        },
        () => {
          return null;
        },
      );

    setAnalysts(analystOnDataBase || []);

    setLoading(false);
  });

  return (
    <div>
      {error && <Alert variant="error">{error}</Alert>}
      {loading ? (
        <Loading />
      ) : (
        <div className="flex flex-row justify-start flex-wrap">
          {analysts.map((analyst) => (
            <AnalystCard key={analyst.id} analyst={analyst} />
          ))}
        </div>
      )}
    </div>
  );
}
