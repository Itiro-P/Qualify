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

  // TODO: pegar o userId da sessão/auth
  const userId = 1;

  async function handleSubmitAll() {
    setLoading(true);
    setError("");
    setSuccess("");
    try {
      // 1. Promover usuário a analista
      await analystService.create(userId, { hourly_rate: hourlyRateAnalyst });

      // 2. Colocar Certificações
      for (const certification of certificationsAnalyst) {
        // 'anda' por cada certificado em certificationAnalyst
        const certificationOfDataBase = await certificationService.list(
          certification.name,
        ); // pega certificação que correspondem ao nome da certificação
        let createdCertication: Certification;
        if (certificationOfDataBase.count == 0) {
          // não tem certificação no banco de dados
          createdCertication = (
            await certificationService.create(certification)
          ).certification; // cria certificação no banco de dados
        } else {
          // tem certificação no banco de dados
          createdCertication = certificationOfDataBase.certifications[0]; // pega primeira correspondência
        }

        await analystService.addCertification(userId, {
          certification_id: createdCertication.id,
        }); // adiciona certificação ao analista
      }

      // 3. Colocar Skills //necessario mudar depois
      for (const skill of skillsAnalyst) {
        const skillOfDataBase = await skillService.list(skill.name); // pega skill que correspondem ao nome da skill
        let createdSkill: Skill;
        if (skillOfDataBase.count == 0) {
          // não tem skill no banco de dados
          createdSkill = (await skillService.create(skill)).skill; // cria skill no banco de dados
        } else {
          // tem skill no banco de dados
          createdSkill = skillOfDataBase.skills[0]; // pega primeira correspondência
        }
        await analystService.addSkill(userId, {
          skill_id: createdSkill.id,
        }); // adiciona
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
