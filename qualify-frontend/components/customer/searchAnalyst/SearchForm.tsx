"use client";

import { FormButton, FormInput } from "@/components/ui";
import { IFormResponse } from "@/types/customer/formResponse";
import { Dispatch, SetStateAction, useEffect } from "react";

function removeSkill(
  setFormResponse: Dispatch<SetStateAction<IFormResponse | null>>,
  skill: string,
) {
  setFormResponse((prev) =>
    prev
      ? { ...prev, skills: prev.skills.filter((item) => item !== skill) }
      : null,
  );
}

function SkillsSaves({
  skill,
  setFormResponse,
}: {
  skill: string;
  setFormResponse: Dispatch<SetStateAction<IFormResponse | null>>;
}) {
  return (
    <FormButton
      variant="outline"
      fullWidth={false}
      onClick={() => removeSkill(setFormResponse, skill)}
      className="m-1 py-2! px-4! text-sm"
      type="button"
    >
      {skill}
    </FormButton>
  );
}

export function SearchForm({
  formResponse,
  setFormResponse,
}: {
  formResponse: IFormResponse | null;
  setFormResponse: Dispatch<SetStateAction<IFormResponse | null>>;
}) {
  useEffect(() => {
    if (!formResponse) {
      setFormResponse({
        min_hourly_rate: 0,
        max_hourly_rate: 0,
        rating: 0,
        skills: [],
        localization: {
          country: "",
          state: "",
          city: "",
        },
      });
    }

    if (formResponse) {
      if (formResponse.min_hourly_rate > formResponse?.max_hourly_rate) {
        setFormResponse({
          ...formResponse,
          max_hourly_rate: formResponse.min_hourly_rate,
        });
      }
    }
  }, [formResponse]);
  return (
    <form>
      <h3>Rating</h3>

      <FormInput
        label="Rating mínimo"
        onChange={(e) =>
          setFormResponse({
            min_hourly_rate: formResponse?.min_hourly_rate || 0,
            max_hourly_rate: formResponse?.max_hourly_rate || 0,
            rating: Number(e.target.value),
            skills: formResponse?.skills || [],
            localization: formResponse?.localization || {
              country: "",
              state: "",
              city: "",
            },
          })
        }
      />

      <h3>Habilidades</h3>

      <FormInput
        label="Habilidades"
        onChange={(e) =>
          setFormResponse({
            min_hourly_rate: formResponse?.min_hourly_rate || 0,
            max_hourly_rate: formResponse?.max_hourly_rate || 0,
            rating: formResponse?.rating || 0,
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

      {formResponse?.skills.length ? (
        <p className="m-2 text-sm text-neutral-slate">Habilidades salvas:</p>
      ) : null}
      <div className="flex flex-row justify-start flex-wrap">
        {formResponse?.skills.map((skill) => (
          <SkillsSaves
            key={skill}
            skill={skill}
            setFormResponse={setFormResponse}
          />
        ))}
      </div>

      <div></div>

      <h3>Preço por hora</h3>

      <FormInput
        label="Preço mínimo por hora"
        value={formResponse?.min_hourly_rate || 0}
        onChange={(e) =>
          setFormResponse({
            min_hourly_rate: Number(e.target.value),
            max_hourly_rate: formResponse?.max_hourly_rate || 0,
            rating: formResponse?.rating || 0,
            skills: formResponse?.skills || [],
            localization: formResponse?.localization || {
              country: "",
              state: "",
              city: "",
            },
          })
        }
      />

      <FormInput
        label="Preço máximo por hora"
        value={formResponse?.max_hourly_rate || 0}
        onChange={(e) =>
          setFormResponse({
            min_hourly_rate: formResponse?.min_hourly_rate || 0,
            max_hourly_rate: Number(e.target.value),
            rating: formResponse?.rating || 0,
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
            min_hourly_rate: formResponse?.min_hourly_rate || 0,
            max_hourly_rate: formResponse?.max_hourly_rate || 0,
            rating: formResponse?.rating || 0,
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
            min_hourly_rate: formResponse?.min_hourly_rate || 0,
            max_hourly_rate: formResponse?.max_hourly_rate || 0,
            rating: formResponse?.rating || 0,
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
            min_hourly_rate: formResponse?.min_hourly_rate || 0,
            max_hourly_rate: formResponse?.max_hourly_rate || 0,
            rating: formResponse?.rating || 0,
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
  );
}
