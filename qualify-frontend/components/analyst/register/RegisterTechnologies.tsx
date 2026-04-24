"use client";

import { ITechnology } from "@/types/analyst/technology";
import { useState } from "react";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<ITechnology>>
) {
  const { name, value } = e.target;

  setForm(
        (prev) =>
        ({
            ...prev,
            [name]: value,
        })
    );
};


function validate(data: ITechnology): Partial<ITechnology> {
    const newErrors: Partial<ITechnology> = {};
    
    if (!data.technology) {
        newErrors.technology = "Tecnologia é obrigatório";
    }
    
    return newErrors;
}

function handleSubmit(
    e: React.FormEvent,
    setErrors: React.Dispatch<React.SetStateAction<Partial<ITechnology>>>,
    form: ITechnology ) {
    e.preventDefault();

    const validationErrors = validate(form);
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length === 0) {
      console.log("Dados enviados:", form);
    }
}

export function RegisterCertifications(){
    const [form, setForm] = useState<ITechnology>({
        technology: ""
    });

    const [errors, setErrors] = useState<Partial<ITechnology>>({});

    return(
        <form 
            onSubmit={(e) => handleSubmit(e, setErrors,form)}
        >
            <div className="flex flex-col gap-4">
                <div>
                <label className="text-sm font-medium">Nome</label>
                <input
                    name="name" 
                    value={form.technology}
                    onChange={(e) => handleChange(e, setForm)}
                    className="w-full border rounded px-3 py-2 mt-1"
                />
                {errors.technology && (
                    <p className="text-red-500 text-sm">{errors.technology}</p>
                )}
                </div>
            </div>
        </form>

    );
}