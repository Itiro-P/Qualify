"use client";

import { useState } from "react";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<string>>
) {
  const { value } = e.target;

  setForm(
        value,
    );
};


function validate(data: string): Partial<string> {
    let newErrors: string = "";
    
    if (!data) {
        newErrors = "Valor por hora é obrigatório";
    }
    
    return newErrors;
}

function handleSubmit(
    e: React.FormEvent,
    setErrors: React.Dispatch<React.SetStateAction<Partial<string>>>,
    form: string ) {
    e.preventDefault();

    const validationErrors = validate(form);
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length === 0) {
      console.log("Dados enviados:", form);
    }
}

export function RegisterCertifications(){
    const [form, setForm] = useState<string>("");

    const [errors, setErrors] = useState<Partial<string>>("");

    return(
        <form 
            onSubmit={(e) => handleSubmit(e, setErrors,form)}
        >
            <div className="flex flex-col gap-4">
                <div>
                <label className="text-sm font-medium">Nome</label>
                <input
                    name="name" 
                    value={form}
                    onChange={(e) => handleChange(e, setForm)}
                    className="w-full border rounded px-3 py-2 mt-1"
                />
                {errors && (
                    <p className="text-red-500 text-sm">{errors}</p>
                )}
                </div>
            </div>
        </form>

    );
}