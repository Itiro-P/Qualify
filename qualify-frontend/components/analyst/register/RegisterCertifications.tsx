"use client";

import { ICertification } from "@/types/analyst/profile/certification";
import { ICertificationSavesProps } from "@/types/analyst/register/certificationSavesProps";
import { IRegisterCertifications } from "@/types/analyst/register/registerCertification";
import { useState } from "react";

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<ICertification>>,
) {
  const { name, value } = e.target;

  setForm((prev) => ({
    ...prev,
    [name]: value,
  }));
}

function validate(certification: ICertification): Partial<ICertification> {
  const newErrors: Partial<ICertification> = {};

  if (!certification.name) {
    newErrors.name = "Nome é obrigatório";
  }

  if (!certification.description) {
    newErrors.description = "Descrição é obrigatório";
  }

  if (!certification.institution) {
    newErrors.institution = "Instituição é obrigatório";
  }

  if (!certification.year) {
    newErrors.year = "Ano é obrigatório";
  } else if (Number.isNaN(Number(certification.year))) {
    newErrors.year = "Ano tem que ser numérico";
  }

  return newErrors;
}

function isEqualCertification(
  certification1: ICertification,
  certification2: ICertification,
): boolean {
  return (
    certification1.name === certification2.name &&
    certification1.description === certification2.description &&
    certification1.institution === certification2.institution &&
    certification1.year === certification2.year
  );
}

function validateCertification(
  certifications: ICertification[],
  certification: ICertification,
): string {
  let newErrors = "";

  if (
    certifications.length > 0 &&
    certifications.find((element) =>
      isEqualCertification(element, certification),
    )
  ) {
    newErrors = "Essa certificação já foi incluída";
  }

  return newErrors;
}

function handleSubmit(
  e: React.FormEvent,
  setError: React.Dispatch<React.SetStateAction<string>>,
  setErrorsCertification: React.Dispatch<
    React.SetStateAction<Partial<ICertification>>
  >,
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<ICertification[]>
  >,
  certificationsAnalyst: ICertification[],
  setCertification: React.Dispatch<React.SetStateAction<ICertification>>,
  certification: ICertification,
) {
  e.preventDefault();

  const validationErrors = validate(certification);
  setErrorsCertification(validationErrors);

  if (Object.keys(validationErrors).length === 0) {
    const error: string = validateCertification(
      certificationsAnalyst,
      certification,
    );
    setError(error);
    if (!error) {
      setCertificationsAnalyst((prev) => [...prev, certification]);
      setCertification({
        name: "",
        description: "",
        institution: "",
        year: "",
      });
      console.log("Dados salvados:", certification);
    }
  }
}

function removeTehnology(
  setTecnologiesAnalyst: React.Dispatch<React.SetStateAction<ICertification[]>>,
  certification: ICertification,
) {
  setTecnologiesAnalyst((prev) => [
    ...prev.filter((item) => !isEqualCertification(item, certification)),
  ]);
}

function CertificationSaves({
  certification,
  setCertificationsAnalyst,
}: ICertificationSavesProps) {
  return (
    <button
      onClick={() => removeTehnology(setCertificationsAnalyst, certification)}
      className="px-4 py-2 m-2 rounded-xl bg-blue-600"
    >
      {certification.name}/{certification.year}
    </button>
  );
}

export function RegisterCertifications({certificationsAnalyst, setCertificationsAnalyst} : IRegisterCertifications) {
  const [certification, setCertification] = useState<ICertification>({
    name: "",
    description: "",
    institution: "",
    year: "",
  });

  const [error, setError] = useState<string>("");

  const [errorsCertification, setErrorsCertification] = useState<
    Partial<ICertification>
  >({});

  return (
    <div>
      <p className="text-sm font-medium">Certificações</p>
      <form
        onSubmit={(e) =>
          handleSubmit(
            e,
            setError,
            setErrorsCertification,
            setCertificationsAnalyst,
            certificationsAnalyst,
            setCertification,
            certification,
          )
        }
        className="mb-8 mt-2 p-5 pt-3 border border-solid rounded-xl"
      >
        <div className="flex flex-col gap-4">
          <div>
            <label className="text-sm font-medium">Nome</label>
            <input
              name="name"
              value={certification.name}
              onChange={(e) => handleChange(e, setCertification)}
              className="w-full border rounded px-3 py-2 mt-1"
            />
            {errorsCertification.name && (
              <p className="text-red-500 text-sm">{errorsCertification.name}</p>
            )}
          </div>
          <div>
            <label className="text-sm font-medium">Descrição</label>
            <input
              name="description"
              value={certification.description}
              onChange={(e) => handleChange(e, setCertification)}
              className="w-full border rounded px-3 py-2 mt-1"
            />
            {errorsCertification.description && (
              <p className="text-red-500 text-sm">
                {errorsCertification.description}
              </p>
            )}
          </div>
          <div>
            <label className="text-sm font-medium">Instituição</label>
            <input
              name="institution"
              value={certification.institution}
              onChange={(e) => handleChange(e, setCertification)}
              className="w-full border rounded px-3 py-2 mt-1"
            />
            {errorsCertification.institution && (
              <p className="text-red-500 text-sm">
                {errorsCertification.institution}
              </p>
            )}
          </div>
          <div>
            <label className="text-sm font-medium">Ano</label>
            <input
              name="year"
              value={certification.year}
              onChange={(e) => handleChange(e, setCertification)}
              className="w-full border rounded px-3 py-2 mt-1"
            />
            {errorsCertification.year && (
              <p className="text-red-500 text-sm">{errorsCertification.year}</p>
            )}
          </div>
        </div>
        {error && <p className="text-red-500 text-sm">{error}</p>}
        <button
          type="submit"
          className="mt-4 bg-blue-600 text-white font-medium px-5 py-2 rounded-lg 
              hover:bg-blue-700 active:scale-95 transition-all duration-200"
        >
          Adicionar
        </button>
        {certificationsAnalyst.length > 0 && (
          <p className="m-2">Certificações salvas:</p>
        )}
        <div className="flex flex-row justify-start flex-wrap">
          {certificationsAnalyst.map((certification) => (
            <CertificationSaves
              key={certification.description}
              {...{ setCertificationsAnalyst, certification }}
            />
          ))}
        </div>
      </form>
    </div>
  );
}
