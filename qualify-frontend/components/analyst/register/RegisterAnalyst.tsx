"use client";
import { useState } from "react";
import { ICertification } from "@/types/analyst/profile/certification";
import { ITechnology } from "@/types/analyst/profile/technology";
import {
  RegisterCertifications,
  RegisterHourlyRate,
  RegisterTechnologies,
} from "@/components/analyst/register";
import { analystService, certificationService, skillService } from "@/libs/services";
import type { ApiError } from "@/libs/api";
import { FormButton, FormPanel, Alert } from "@/components/ui";

export function RegisterAnalyst() {
  const [certificationsAnalyst, setCertificationsAnalyst] = useState<ICertification[]>([]);
  const [technologiesAnalyst, setTechnologiesAnalyst] = useState<ITechnology[]>([]);
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // TODO: pegar o userId da sessão/auth (ex: useSession do NextAuth)
  const userId = 1;

  async function getOrCreateCertification(cert: ICertification): Promise<number> {
    const searchData = await certificationService.list(cert.name);

    if (searchData.count > 0) {
      return searchData.certifications[0].id!;
    }

    const created = await certificationService.create({
      description: cert.description,
      institution: cert.institution,
      name: cert.name,
      year: Number(cert.year),
    });

    return created.certification.id!;
  }

  async function getOrCreateSkill(tech: ITechnology): Promise<number> {
    const searchData = await skillService.list(tech.technology);

    if (searchData.count > 0) {
      return searchData.skills[0].id!;
    }

    const created = await skillService.create({ name: tech.technology });
    return created.skill.id!;
  }

  async function handleSubmitAll() {
    setLoading(true);
    setError("");
    setSuccess("");
    try {
      // 1. Promover usuário a analista
      await analystService.create(userId, { hourly_rate: hourlyRateAnalyst });

      // 2. Certificações: busca ou cria, depois vincula ao analista
      const certIds = await Promise.all(
        certificationsAnalyst.map(getOrCreateCertification)
      );

      await Promise.all(
        certIds.map((certId) =>
          analystService.addCertification(userId, { certification_id: certId })
        )
      );

      // 3. Skills: busca ou cria, depois vincula ao analista
      const skillIds = await Promise.all(
        technologiesAnalyst.map(getOrCreateSkill)
      );

      await Promise.all(
        skillIds.map((skillId) =>
          analystService.addSkill(userId, { skill_id: skillId })
        )
      );

      setSuccess("Cadastro de analista realizado com sucesso!");
    } catch (err) {
      const apiError = err as ApiError;
      if (apiError.status) {
        switch (apiError.status) {
          case 400: setError(apiError.message || "Dados inválidos"); break;
          case 401: setError("Não autorizado"); break;
          case 422: setError("Dados mal formatados"); break;
          default:  setError("Erro no servidor. Tente novamente.");
        }
      } else {
        console.error(err);
        setError("Ocorreu um erro inesperado. Verifique o console.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <FormPanel title="Cadastro de Analista" description="Preencha suas certificações, tecnologias e valor por hora." maxWidth="max-w-2xl">
      {error && <Alert variant="error">{error}</Alert>}
      {success && <Alert variant="success">{success}</Alert>}
      <RegisterCertifications
        certificationsAnalyst={certificationsAnalyst}
        setCertificationsAnalyst={setCertificationsAnalyst}
      />
      <RegisterTechnologies
        technologiesAnalyst={technologiesAnalyst}
        setTechnologiesAnalyst={setTechnologiesAnalyst}
      />
      <RegisterHourlyRate
        hourlyRateAnalyst={hourlyRateAnalyst}
        setHourlyRateAnalyst={setHourlyRateAnalyst}
      />
      <FormButton
        onClick={handleSubmitAll}
        loading={loading}
        loadingText="Cadastrando..."
        className="mt-4"
      >
        Cadastrar-se
      </FormButton>
    </FormPanel>
  );
}