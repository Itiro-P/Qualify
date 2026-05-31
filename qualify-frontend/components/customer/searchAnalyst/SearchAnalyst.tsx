"use client";

import { Alert, FormInput } from "@/components/ui";
import { Loading } from "@/components/ui/Loading";
import { Analyst } from "@/types/services/analyst";
import { useState } from "react";

export function SearchAnalyst() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [analysts, setAnalysts] = useState<Analyst[]>([]);

  return loading ? (
    <Loading />
  ) : (
    <div>
      {error && <Alert variant="error">{error}</Alert>}
      <form>
        <h3>Categorias</h3>
        <FormInput label="Eletrônicos" type="checkbox" />
        <FormInput label="Informática" type="checkbox" />
        <FormInput label="Celulares" type="checkbox" />
        <FormInput label="Casa" type="checkbox" />
        <h3>Preço</h3>
        <FormInput label="Preço" type="range" min="0" max="5000" value="2500" />
        <h3>Avaliação</h3>
        <FormInput label="★★★★★" type="checkbox" />
        <FormInput label="★★★★☆ ou mais" type="checkbox" />
        <FormInput label="★★★☆☆ ou mais" type="checkbox" />
        <h3>Entrega</h3>
        <FormInput label="Frete grátis" type="checkbox" />
        <FormInput label="Full" type="checkbox" />
        <FormInput label="Entrega rápida" type="checkbox" />
      </form>
    </div>
  );
}
