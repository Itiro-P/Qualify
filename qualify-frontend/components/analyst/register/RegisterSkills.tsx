"use client";

import { Skill } from "@/types/services";
import { useState } from "react";
import { FormInput, FormButton, Alert } from "@/components/ui";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<Skill>>,
) {
  const { value } = e.target;

  setForm({ id: 0, name: value });
}

function validateSkill(skills: Skill[], skill: string): string {
  let newErrors: string = "";

  if (
    skills.length > 0 &&
    skills.find((element) => element.name === skill)
  ) {
    newErrors = "Essa tecnologia já foi incluída";
  }

  return newErrors;
}

function handleSubmit(
  e: React.FormEvent<HTMLFormElement>,
  setErrors: React.Dispatch<React.SetStateAction<string>>,
  setSkillsAnalyst: React.Dispatch<React.SetStateAction<Skill[]>>,
  skillsAnalyst: Skill[],
  skill: Skill,
  setskill: React.Dispatch<React.SetStateAction<Skill>>,
) {
  e.preventDefault();

  const validationErrors = validateSkill(
    skillsAnalyst,
    skill.name,
  );

  setErrors(validationErrors);

  if (!validationErrors) {
    console.log("Dados salvados:", skill.name);

    setSkillsAnalyst((prev) => [
      ...prev,
      { id: 0, name: skill.name },
    ]);

    setskill({ id: 0, name: "" });
  }
}

function removeskill(
  setSkillsAnalyst: React.Dispatch<React.SetStateAction<Skill[]>>,
  skill: Skill,
) {
  setSkillsAnalyst((prev) => [
    ...prev.filter((item) => item.name !== skill.name),
  ]);
}

function SkillSaves({
  skill,
  setSkillsAnalyst,
}: {
  skill: Skill;
  setSkillsAnalyst: React.Dispatch<React.SetStateAction<Skill[]>>;
}) {
  return (
    <FormButton
      variant="outline"
      fullWidth={false}
      onClick={() => removeskill(setSkillsAnalyst, skill)}
      className="m-1 py-2! px-4! text-sm"
      type="button"
    >
      {skill.name}
    </FormButton>
  );
}

export function RegisterSkills({
  skillsAnalyst,
  setSkillsAnalyst,
}: {
  skillsAnalyst: Skill[];
  setSkillsAnalyst: React.Dispatch<React.SetStateAction<Skill[]>>;
}) {
  const [errors, setErrors] = useState<string>("");
  const [skill, setskill] = useState<Skill>({
    id: 0,
    name: "",
  });

  return (
    <form
      onSubmit={(e) =>
        handleSubmit(
          e,
          setErrors,
          setSkillsAnalyst,
          skillsAnalyst,
          skill,
          setskill,
        )
      }
      className="my-4 p-5 pt-3 border border-white/10 rounded-xl"
    >
      <div className="flex flex-col gap-4">
        <FormInput
          label="Tecnologia"
          name="skill"
          value={skill.name}
          onChange={(e) => handleChange(e, setskill)}
          required
        />
        {errors && <Alert variant="error">{errors}</Alert>}
      </div>
      <FormButton type="submit" fullWidth={false} className="mt-4 px-6">
        Adicionar
      </FormButton>
      {skillsAnalyst.length > 0 && (
        <p className="m-2 text-sm text-neutral-slate">Tecnologias salvas:</p>
      )}
      <div className="flex flex-row justify-start flex-wrap">
        {skillsAnalyst.map((skill) => (
          <SkillSaves
            key={skill.name}
            {...{ setSkillsAnalyst, skill }}
          />
        ))}
      </div>
    </form>
  );
}
