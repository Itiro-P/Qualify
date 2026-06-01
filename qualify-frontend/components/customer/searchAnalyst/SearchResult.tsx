"use client";

import { Alert } from "@/components/ui";
import { Loading } from "@/components/ui/Loading";
import { Analyst } from "@/types/services/analyst";
import { analystService } from "@/libs/services/analystService";
import { skillService } from "@/libs/services/skillService";
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

  useState(async () => {
    if (!formResponse) return;

    setLoading(true);
    setError("");
    setAnalysts([]);

    async function haveSkill(
      Analyst: Analyst,
      skills: string[],
    ): Promise<boolean> {
      const analystSkillsId = await analystService.listSkills(Analyst.id);
      const analystSkills: string[] = [];

      for (const skill of analystSkillsId ?? []) {
        try {
          const resp = await skillService.getById(skill.skill_id);

          if (resp?.name) {
            analystSkills.push(resp.name);
          }
        } catch {
          continue;
        }
      }
      return skills.every((skill) => analystSkills.includes(skill));
    }

    // Pegar os analistas recomendados com base na resposta do formulário
    const analystOnDataBase = await analystService
      .list({
        country: formResponse.localization.country,
        country_state: formResponse.localization.state,
        city: formResponse.localization.city,
        min_hourly_rate: formResponse.min_hourly_rate,
        max_hourly_rate: formResponse.max_hourly_rate,
        min_mean_rating: formResponse.rating,
      })
      .then(
        (resp) => {
          return resp.analysts;
        },
        () => {
          return null;
        },
      );

    if (analystOnDataBase) {
      analystOnDataBase.forEach(async (analyst) => {
        if (await haveSkill(analyst, formResponse.skills)) {
          setAnalysts((prev) => [...prev, analyst]);
        }
      });
    }

    setLoading(false);
  });

  return (
    <div>
      {error && <Alert variant="error">{error}</Alert>}
      {loading ? <Loading /> : <div></div>}
    </div>
  );
}
