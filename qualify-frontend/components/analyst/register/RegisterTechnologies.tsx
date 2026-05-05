"use client";

import { ITechnology } from "@/types/analyst/profile/technology";
import { ITechnologySavesProps } from "@/types/analyst/register/technologySavesProps";
import { IRegisterTechnology } from "@/types/analyst/register/registerTechnology";
import { useState } from "react";

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
}: ITechnologySavesProps) {
  return (
    <button
      onClick={() => removeTechnology(setTechnologiesAnalyst, technology)}
      className="px-4 py-2 m-2 rounded-xl bg-blue-600"
    >
      {technology.technology}
    </button>
  );
}

export function RegisterTechnologies({technologiesAnalyst, setTechnologiesAnalyst}:IRegisterTechnology) {
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
      className="my-4 p-5 pt-3 border border-solid rounded-xl"
    >
      <div className="flex flex-col gap-4">
        <div>
          <label className="text-sm mx-3 font-medium">Tecnologia</label>
          <input
            name="technology"
            value={technology.technology}
            onChange={(e) => handleChange(e, setTechnology)}
            className="w-full border rounded px-3 py-2 mt-1"
          />
          {errors && <p className="text-red-500 text-sm mt-2">{errors}</p>}
        </div>
      </div>
      <button
        type="submit"
        className="mt-4 bg-blue-600 text-white font-medium px-5 py-2 rounded-lg hover:bg-blue-700 active:scale-95 transition-all duration-200"
      >
        Adicionar
      </button>
      {technologiesAnalyst.length > 0 && (
        <p className="m-2">Tecnologias salvas:</p>
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
