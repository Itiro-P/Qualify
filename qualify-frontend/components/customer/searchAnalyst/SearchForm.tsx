"use client";

import { FormButton, FormInput } from "@/components/ui";
import { IFormResponse } from "@/types/customer/formResponse";
import { X } from "lucide-react";
import { Dispatch, SetStateAction, useState } from "react";

function SkillTag({
  skill,
  onRemove,
}: {
  skill: string;
  onRemove: (skill: string) => void;
}) {
  return (
    <button
      type="button"
      onClick={() => onRemove(skill)}
      className="inline-flex items-center gap-1 rounded-full border border-accent/30 bg-accent/10 px-3 py-1 text-xs font-medium text-accent hover:bg-accent/20 transition-colors cursor-pointer"
    >
      {skill}
      <X className="w-3 h-3" />
    </button>
  );
}

export function SearchForm({
  formResponse,
  setFormResponse,
  onSearch,
}: {
  formResponse: IFormResponse;
  setFormResponse: Dispatch<SetStateAction<IFormResponse>>;
  onSearch: () => void;
}) {
  const [skill, setSkill] = useState<string>("");

  function update(patch: Partial<IFormResponse>) {
    setFormResponse((prev) => ({ ...prev, ...patch }));
  }

  function updateLocation(patch: Partial<IFormResponse["localization"]>) {
    setFormResponse((prev) => ({
      ...prev,
      localization: { ...prev.localization, ...patch },
    }));
  }

  function addSkill() {
    const value = skill.trim();
    if (!value) return;
    setFormResponse((prev) =>
      prev.skills.includes(value)
        ? prev
        : { ...prev, skills: [...prev.skills, value] },
    );
    setSkill("");
  }

  function removeSkill(target: string) {
    setFormResponse((prev) => ({
      ...prev,
      skills: prev.skills.filter((item) => item !== target),
    }));
  }

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        onSearch();
      }}
      className="glass-panel rounded-2xl p-6 flex flex-col gap-6 lg:sticky lg:top-24"
    >
      <div>
        <h3 className="text-sm font-semibold text-white mb-3">Competências</h3>
        <FormInput
          label="Adicionar competência"
          placeholder="Ex.: Java, Selenium..."
          value={skill}
          onChange={(e) => setSkill(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addSkill();
            }
          }}
          hint="Pressione Enter para adicionar"
        />
        {formResponse.skills.length > 0 && (
          <div className="mt-3 flex flex-wrap gap-2">
            {formResponse.skills.map((s) => (
              <SkillTag key={s} skill={s} onRemove={removeSkill} />
            ))}
          </div>
        )}
      </div>

      <div>
        <h3 className="text-sm font-semibold text-white mb-3">
          Avaliação mínima
        </h3>
        <FormInput
          label="Nota (0 a 5)"
          type="number"
          min={0}
          max={5}
          step={0.5}
          value={formResponse.rating || 0}
          onChange={(e) => update({ rating: Number(e.target.value) })}
        />
      </div>

      <div>
        <h3 className="text-sm font-semibold text-white mb-3">
          Preço por hora
        </h3>
        <div className="flex flex-col gap-3">
          <FormInput
            label="Mínimo"
            type="number"
            min={0}
            value={formResponse.min_hourly_rate || 0}
            onChange={(e) =>
              update({ min_hourly_rate: Number(e.target.value) })
            }
          />
          <FormInput
            label="Máximo"
            type="number"
            min={0}
            value={formResponse.max_hourly_rate || 0}
            onChange={(e) =>
              update({ max_hourly_rate: Number(e.target.value) })
            }
          />
        </div>
      </div>

      <div>
        <h3 className="text-sm font-semibold text-white mb-3">Localização</h3>
        <div className="flex flex-col gap-3">
          <FormInput
            label="País"
            value={formResponse.localization.country}
            onChange={(e) => updateLocation({ country: e.target.value })}
          />
          <FormInput
            label="Estado"
            value={formResponse.localization.state}
            onChange={(e) => updateLocation({ state: e.target.value })}
          />
          <FormInput
            label="Cidade"
            value={formResponse.localization.city}
            onChange={(e) => updateLocation({ city: e.target.value })}
          />
        </div>
      </div>

      <FormButton type="submit">Buscar analistas</FormButton>
    </form>
  );
}
