"use client";

import { useState } from "react";
import { FormInput, FormButton } from "@/components/ui";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<number>>,
) {
  const { value } = e.target;

  setForm(Number(value));
}

function validate(data: number): string {
  let newErrors: string = "";

  if (!data) {
    newErrors = "Valor por hora é obrigatório";
  }
  if (Number.isNaN(data)) {
    newErrors = "Apenas números são aceitos";
  }

  return newErrors;
}

function handleSubmit(
  e: React.FormEvent,
  setErrors: React.Dispatch<React.SetStateAction<string>>,
  hourlyRateAnalyst: number,
) {
  e.preventDefault();

  const validationErrors = validate(hourlyRateAnalyst);
  setErrors(validationErrors);

  if (!validationErrors) {
    console.log("Dados enviados:", hourlyRateAnalyst);
  }
}

export function RegisterHourlyRate({
  hourlyRateAnalyst,
  setHourlyRateAnalyst,
}: {
  hourlyRateAnalyst: number;
  setHourlyRateAnalyst: React.Dispatch<React.SetStateAction<number>>;
}) {
  const [errors, setErrors] = useState<string>("");

  return (
    <form
      onSubmit={(e) => handleSubmit(e, setErrors, hourlyRateAnalyst)}
      className="my-8 p-5 pt-3 border border-white/10 rounded-xl"
    >
      <div className="flex flex-col gap-4">
        <FormInput
          label="Valor por hora ($)"
          name="hourlyRate"
          type="number"
          value={hourlyRateAnalyst}
          onChange={(e) => handleChange(e, setHourlyRateAnalyst)}
          error={errors}
          required
        />
      </div>
      <FormButton type="submit" fullWidth={false} className="mt-4 px-6">
        Salvar
      </FormButton>
    </form>
  );
}
