"use client";

import { Certification } from "@/types/services/certification";
import { useState } from "react";
import { FormInput, FormButton, Alert } from "@/components/ui";

interface CertificationError {
  name: string;
  description: string;
  institution: string;
  year: string;
}

function handleChange(
  e: React.ChangeEvent<HTMLInputElement>,
  setForm: React.Dispatch<React.SetStateAction<Certification>>,
) {
  const { name, value } = e.target;

  setForm((prev) => ({
    ...prev,
    [name]: value,
  }));
}

function validate(certification: Certification): Partial<CertificationError> {
  const newErrors: Partial<CertificationError> = {};

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
  certification1: Certification,
  certification2: Certification,
): boolean {
  return (
    certification1.name === certification2.name &&
    certification1.description === certification2.description &&
    certification1.institution === certification2.institution &&
    certification1.year === certification2.year
  );
}

function validateCertification(
  certifications: Certification[],
  certification: Certification,
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
    React.SetStateAction<Partial<CertificationError>>
  >,
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<Certification[]>
  >,
  certificationsAnalyst: Certification[],
  setCertification: React.Dispatch<React.SetStateAction<Certification>>,
  certification: Certification,
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
        id: 0,
        name: "",
        description: "",
        institution: "",
        year: -1,
      });
      console.log("Dados salvados:", certification);
    }
  }
}

function removeCertification(
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<Certification[]>
  >,
  certification: Certification,
) {
  setCertificationsAnalyst((prev) => [
    ...prev.filter((item) => !isEqualCertification(item, certification)),
  ]);
}

function CertificationSaves({
  certification,
  setCertificationsAnalyst,
}: {
  certification: Certification;
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<Certification[]>
  >;
}) {
  return (
    <FormButton
      variant="outline"
      fullWidth={false}
      onClick={() =>
        removeCertification(setCertificationsAnalyst, certification)
      }
      className="m-1 py-2! px-4! text-sm"
      type="button"
    >
      {certification.name}/{certification.year} ✕
    </FormButton>
  );
}

export function RegisterCertifications({
  certificationsAnalyst,
  setCertificationsAnalyst,
}: {
  certificationsAnalyst: Certification[];
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<Certification[]>
  >;
}) {
  const [certification, setCertification] = useState<Certification>({
    id: -1,
    name: "",
    description: "",
    institution: "",
    year: 0,
  });

  const [error, setError] = useState<string>("");

  const [errorsCertification, setErrorsCertification] = useState<
    Partial<CertificationError>
  >({});

  return (
    <div>
      <p className="text-sm font-medium text-white/80 mb-2">Certificações</p>
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
        className="mb-8 mt-2 p-5 pt-3 border border-white/10 rounded-xl"
      >
        <div className="flex flex-col gap-4">
          <FormInput
            label="Nome"
            name="name"
            value={certification.name}
            onChange={(e) => handleChange(e, setCertification)}
            error={errorsCertification.name}
            required
          />
          <FormInput
            label="Descrição"
            name="description"
            value={certification.description}
            onChange={(e) => handleChange(e, setCertification)}
            error={errorsCertification.description}
            required
          />
          <FormInput
            label="Instituição"
            name="institution"
            value={certification.institution}
            onChange={(e) => handleChange(e, setCertification)}
            error={errorsCertification.institution}
            required
          />
          <FormInput
            label="Ano"
            name="year"
            value={certification.year > 0 ? certification.year : ""}
            onChange={(e) => handleChange(e, setCertification)}
            error={errorsCertification.year}
            required
          />
        </div>
        {error && <Alert variant="error">{error}</Alert>}
        <FormButton type="submit" fullWidth={false} className="mt-4 px-6">
          Adicionar
        </FormButton>
        {certificationsAnalyst.length > 0 && (
          <p className="m-2 text-sm text-neutral-slate">
            Certificações salvas:
          </p>
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
