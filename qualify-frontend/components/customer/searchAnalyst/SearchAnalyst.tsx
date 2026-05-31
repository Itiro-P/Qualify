"use client";

import { Alert, FormInput } from "@/components/ui";
import { Loading } from "@/components/ui/Loading";
import { Analyst } from "@/types/services/analyst";
import { useState } from "react";

interface Location {
  country: string;
  state: string;
  city: string;
}

interface IFormResponse {
  value: number;
  skills: string[];
  localization: Location;
}

export function SearchAnalyst() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [analysts, setAnalysts] = useState<Analyst[]>([]);
  const [formResponse, setFormResponse] = useState<IFormResponse | null>(null);

  return loading ? (
    <Loading />
  ) : (
    <div>
      {error && <Alert variant="error">{error}</Alert>}
      <form>
        <h3>Habilidades</h3>
        <FormInput
          label="Habilidades"
          onChange={(e) =>
            setFormResponse({
              value: formResponse?.value || 0,
              skills: formResponse
                ? [...formResponse.skills, e.target.value]
                : [e.target.value],
              localization: formResponse?.localization || {
                country: "",
                state: "",
                city: "",
              },
            })
          }
        />
        <h3>Preço por hora</h3>
        <FormInput
          label="Preço"
          type="range"
          min="0"
          max="5000"
          value={formResponse?.value || 0}
          onChange={(e) =>
            setFormResponse({
              value: Number(e.target.value),
              skills: formResponse?.skills || [],
              localization: formResponse?.localization || {
                country: "",
                state: "",
                city: "",
              },
            })
          }
        />
        <h3>Localização</h3>
        <FormInput
          label="País"
          value={formResponse?.localization?.country || ""}
          onChange={(e) =>
            setFormResponse({
              value: formResponse?.value || 0,
              skills: formResponse?.skills || [],
              localization: formResponse?.localization
                ? {
                    ...formResponse.localization,
                    country: e.target.value,
                  }
                : {
                    country: e.target.value,
                    state: "",
                    city: "",
                  },
            })
          }
        />
        <FormInput
          label="Estado"
          value={formResponse?.localization?.state || ""}
          onChange={(e) =>
            setFormResponse({
              value: formResponse?.value || 0,
              skills: formResponse?.skills || [],
              localization: formResponse?.localization
                ? {
                    ...formResponse.localization,
                    state: e.target.value,
                  }
                : {
                    country: "",
                    state: e.target.value,
                    city: "",
                  },
            })
          }
        />
        <FormInput
          label="Cidade"
          value={formResponse?.localization?.city || ""}
          onChange={(e) =>
            setFormResponse({
              value: formResponse?.value || 0,
              skills: formResponse?.skills || [],
              localization: formResponse?.localization
                ? {
                    ...formResponse.localization,
                    city: e.target.value,
                  }
                : {
                    country: "",
                    state: "",
                    city: e.target.value,
                  },
            })
          }
        />
      </form>
    </div>
  );
}
