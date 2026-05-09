"use client";
import { useState } from "react";
import { Certification } from "@/types/services/certification";
import { Skill } from "@/types/services/skill";
import {
  RegisterCertifications,
  RegisterHourlyRate,
  RegisterSkills,
} from "@/components/analyst/register";
import {
  analystService,
  certificationService,
  skillService,
} from "@/libs/services";
import type { ApiError } from "@/libs/api";
import { FormButton, FormPanel, Alert } from "@/components/ui";

export function RegisterAnalyst() {
  const [certificationsAnalyst, setCertificationsAnalyst] = useState<
    Certification[]
  >([]);
  const [skillsAnalyst, setSkillsAnalyst] = useState<Skill[]>([]);
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<number>(-1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // TODO: pegar o userId da sessão/auth (ex: useSession do NextAuth)
  const userId = 1;

  async function handleSubmitAll() {
    setLoading(true);
    setError("");
    setSuccess("");
    try {
      // 1. Promover usuário a analista
      await analystService.create(userId, { hourly_rate: hourlyRateAnalyst });

      // 2. Colocar Certificações //necessario mudar depois
      for (const certification of certificationsAnalyst) {
        const temp = await certificationService.create(certification);
        await analystService.addCertification(userId, {
          certification_id: temp.certification.id,
        });
      }

      // 3. Colocar Skills //necessario mudar depois
      for (const skill of skillsAnalyst) {
        const temp = await skillService.create(skill);
        await analystService.addSkill(userId, { skill_id: temp.skill.id });
      }

      setSuccess("Cadastro de analista realizado com sucesso!");
    } catch (err) {
      const apiError = err as ApiError;
      if (apiError.status) {
        switch (apiError.status) {
          case 400:
            setError(apiError.message || "Dados inválidos");
            break;
          case 401:
            setError("Não autorizado");
            break;
          case 422:
            setError("Dados mal formatados");
            break;
          default:
            setError("Erro no servidor. Tente novamente.");
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
    <FormPanel
      title="Cadastro de Analista"
      description="Preencha suas certificações, tecnologias e valor por hora."
      maxWidth="max-w-2xl"
    >
      {error && <Alert variant="error">{error}</Alert>}
      {success && <Alert variant="success">{success}</Alert>}
      <RegisterCertifications
        certificationsAnalyst={certificationsAnalyst}
        setCertificationsAnalyst={setCertificationsAnalyst}
      />
      <RegisterSkills
        skillsAnalyst={skillsAnalyst}
        setSkillsAnalyst={setSkillsAnalyst}
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
