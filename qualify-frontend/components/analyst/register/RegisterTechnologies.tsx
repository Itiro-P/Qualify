"use client";

import { ITechnology } from "@/types/analyst/profile/technology";
import { useState } from "react";
import { FormInput, FormButton, Alert } from "@/components/ui";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<ITechnology>>,
) {
  const { value } = e.target;

  setForm({ technology: value });
}

function validateTechnology(
  technologys: ITechnology[],
  technology: string,
): string {
  let newErrors: string = "";

  if (
    technologys.length > 0 &&
    technologys.find((element) => element.technology === technology)
  ) {
    newErrors = "Essa tecnologia já foi incluída";
  }

  return newErrors;
}

function handleSubmit(
  e: React.FormEvent<HTMLFormElement>,
  setErrors: React.Dispatch<React.SetStateAction<string>>,
  setTechnologiesAnalyst: React.Dispatch<React.SetStateAction<ITechnology[]>>,
  technologiesAnalyst: ITechnology[],
  technology: ITechnology,
  setTechnology: React.Dispatch<React.SetStateAction<ITechnology>>,
) {
  e.preventDefault();

  const validationErrors = validateTechnology(
    technologiesAnalyst,
    technology.technology,
  );

  setErrors(validationErrors);

  if (!validationErrors) {
    console.log("Dados salvados:", technology.technology);

    setTechnologiesAnalyst((prev) => [
      ...prev,
      { technology: technology.technology },
    ]);

    setTechnology({ technology: "" });
  }
}

function removeTechnology(
  setTechnologiesAnalyst: React.Dispatch<React.SetStateAction<ITechnology[]>>,
  technology: ITechnology,
) {
  setTechnologiesAnalyst((prev) => [
    ...prev.filter((item) => item.technology !== technology.technology),
  ]);
}

function TechnologySaves({
  technology,
  setTechnologiesAnalyst,
}: {
  technology: ITechnology;
  setTechnologiesAnalyst: React.Dispatch<React.SetStateAction<ITechnology[]>>;
}) {
  return (
    <FormButton
      variant="outline"
      fullWidth={false}
      onClick={() => removeTechnology(setTechnologiesAnalyst, technology)}
      className="m-1 py-2! px-4! text-sm"
      type="button"
    >
      {technology.technology} ✕
    </FormButton>
  );
}

export function RegisterTechnologies({
  technologiesAnalyst,
  setTechnologiesAnalyst,
}: {
  technologiesAnalyst: ITechnology[];
  setTechnologiesAnalyst: React.Dispatch<React.SetStateAction<ITechnology[]>>;
}) {
  const [errors, setErrors] = useState<string>("");
  const [technology, setTechnology] = useState<ITechnology>({
    technology: "",
  });

  return (
    <form
      onSubmit={(e) =>
        handleSubmit(
          e,
          setErrors,
          setTechnologiesAnalyst,
          technologiesAnalyst,
          technology,
          setTechnology,
        )
      }
      className="my-4 p-5 pt-3 border border-white/10 rounded-xl"
    >
      <div className="flex flex-col gap-4">
        <FormInput
          label="Tecnologia"
          name="technology"
          value={technology.technology}
          onChange={(e) => handleChange(e, setTechnology)}
          required
        />
        {errors && <Alert variant="error">{errors}</Alert>}
      </div>
      <FormButton type="submit" fullWidth={false} className="mt-4 px-6">
        Adicionar
      </FormButton>
      {technologiesAnalyst.length > 0 && (
        <p className="m-2 text-sm text-neutral-slate">Tecnologias salvas:</p>
      )}
      <div className="flex flex-row justify-start flex-wrap">
        {technologiesAnalyst.map((technology) => (
          <TechnologySaves
            key={technology.technology}
            {...{ setTechnologiesAnalyst, technology }}
          />
        ))}
      </div>
    </form>
  );
}
