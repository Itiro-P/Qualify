"use client";

import { ICertification } from "@/types/analyst/certification";
import { useState } from "react";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<ICertification>>
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


function validate(data: ICertification): Partial<ICertification> {
    const newErrors: Partial<ICertification> = {};
    
    if (!data.name) {
        newErrors.name = "Nome é obrigatório";
    }
    
    if (!data.description) {
        newErrors.description = "Descrição é obrigatório";
    }
    
    if (!data.institution) {
        newErrors.institution = "Instituição é obrigatório";
    }
    
    if (!data.year) {
        newErrors.year = "Ano é obrigatório";
    }
    
    return newErrors;
}

function handleSubmit(
    e: React.FormEvent,
    setErrors: React.Dispatch<React.SetStateAction<Partial<ICertification>>>,
    form: ICertification ) {
    e.preventDefault();

    const validationErrors = validate(form);
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length === 0) {
      console.log("Dados enviados:", form);
    }
}

export function RegisterCertifications(){
    const [form, setForm] = useState<ICertification>({
        name: "" ,description: "" ,institution: "" , year: ""
    });

    const [errors, setErrors] = useState<Partial<ICertification>>({});

    return(
        <form 
            onSubmit={(e) => handleSubmit(e, setErrors,form)}
        >
            <div className="flex flex-col gap-4">
                <div>
                <label className="text-sm font-medium">Nome</label>
                <input
                    name="name" 
                    value={form.name}
                    onChange={(e) => handleChange(e, setForm)}
                    className="w-full border rounded px-3 py-2 mt-1"
                />
                {errors.name && (
                    <p className="text-red-500 text-sm">{errors.name}</p>
                )}
                </div>
            </div>
        </form>

    );
}