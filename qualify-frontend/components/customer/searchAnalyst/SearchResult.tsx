"use client";

import { Alert } from "@/components/ui";
import { Loading } from "@/components/ui/Loading";
import { Analyst } from "@/types/services/analyst";
import { IFormResponse } from "@/types/customer/formResponse";
import { Dispatch, SetStateAction, useState } from "react";

export function SearchResult({
  formResponse,
  setFormResponse,
}: {
  formResponse: IFormResponse | null;
  setFormResponse: Dispatch<SetStateAction<IFormResponse | null>>;
}) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [analysts, setAnalysts] = useState<Analyst[]>([]);

  useState(() => {
    if (!formResponse) return;

    setLoading(true);
    setError("");
    setAnalysts([]);

    // Pegar os analistas recomendados com base na resposta do formulário
  });

  return (
    <div>
      {error && <Alert variant="error">{error}</Alert>}
      {loading ? <Loading /> : <div></div>}
    </div>
  );
}
