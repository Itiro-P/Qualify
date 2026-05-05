"use client";

import { useState } from "react";
import { IRegisterHourlyRate } from "@/types/analyst/register/registerHourlyRate";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<string>>,
) {
  const { value } = e.target;

  setForm(value);
}

function validate(data: string): string {
  let newErrors: string = "";

  if (!data) {
    newErrors = "Valor por hora é obrigatório";
  }
  if (Number.isNaN(Number(data))) {
    newErrors = "Apenas números são aceitos";
  }

  return newErrors;
}

function handleSubmit(
  e: React.FormEvent,
  setErrors: React.Dispatch<React.SetStateAction<string>>,
  hourlyRateAnalyst: string,
) {
  e.preventDefault();

  const validationErrors = validate(hourlyRateAnalyst);
  setErrors(validationErrors);

  if (!validationErrors) {
    console.log("Dados enviados:", hourlyRateAnalyst);
  }
}

export function RegisterHourlyRate({hourlyRateAnalyst, setHourlyRateAnalyst}:IRegisterHourlyRate) {

  const [errors, setErrors] = useState<string>("");

  return (
    <form
      onSubmit={(e) => handleSubmit(e, setErrors, hourlyRateAnalyst)}
      className="my-8 p-5 pt-3 border border-solid rounded-xl"
    >
      <div className="flex flex-col gap-4">
        <div>
          <label className="text-sm font-medium">Valor por hora ($)</label>
          <input
            name="name"
            value={hourlyRateAnalyst}
            onChange={(e) => handleChange(e, setHourlyRateAnalyst)}
            className="w-full border rounded px-3 py-2 mt-1"
          />
          {errors && <p className="text-red-500 text-sm">{errors}</p>}
        </div>
      </div>
      <button
        type="submit"
        className="mt-4 bg-blue-600 text-white font-medium px-5 py-2 rounded-lg 
             hover:bg-blue-700 active:scale-95 transition-all duration-200"
      >
        Salvar
      </button>
    </form>
  );
}
