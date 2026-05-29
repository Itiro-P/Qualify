"use client";

import { useState } from "react";
import { Certification } from "@/types/services/certification";
import { Skill } from "@/types/services/skill";
import { User } from "@/types/services";
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

export function RegisterAnalyst({ userSession }: { userSession: User }) {
  const [certificationsAnalyst, setCertificationsAnalyst] = useState<
    Certification[]
  >([]);
  const [skillsAnalyst, setSkillsAnalyst] = useState<Skill[]>([]);
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<number>(-1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  async function createAnalyst(): Promise<number> {
    try {
      await analystService.create(userSession.id, {
        hourly_rate: hourlyRateAnalyst,
      });
    } catch (err) {
      console.error("Erro ao atualizar usuário para analista:", err);
      setError(
        "Ocorreu um erro ao atualizar o usuário para analista. Tente novamente.",
      );
      return -1; // retorna -1 para indicar falha
    }
    return 0;
  }

  async function updateCertifications(): Promise<number> {
    try {
      for (const certification of certificationsAnalyst) {
        // 'anda' por cada certificado em certificationAnalyst
        const certificationOfDataBase = await certificationService.list(
          certification.name,
        ); // pega certificação que correspondem ao nome da certificação
        let createdCertication: Certification | null;
        if (certificationOfDataBase == null) {
          // não tem certificação no banco de dados
          createdCertication = await certificationService.create(certification); // cria certificação no banco de dados
        } else {
          // tem certificação no banco de dados
          createdCertication = certificationOfDataBase[0]; // pega primeira correspondência
        }
        if (createdCertication != null) {
          // verifica se criou ou pegou certification corretamente
          await analystService.addCertification(userSession.id, {
            certification_id: createdCertication.id,
          }); // adiciona certificação ao analista
        }
      }
    } catch (err) {
      console.error("Erro ao atualizar certificações:", err);
      setError(
        "Ocorreu um erro ao atualizar as certificações. Tente novamente.",
      );
      return -1; // retorna -1 para indicar falha
    }
    return 0;
  }

  async function updateSkills(): Promise<number> {
    try {
      for (const skill of skillsAnalyst) {
        const skillOfDataBase = await skillService.list(skill.name); // pega skill que correspondem ao nome da skill
        let createdSkill: Skill | null;
        if (skillOfDataBase == null) {
          // não tem skill no banco de dados
          createdSkill = await skillService.create(skill); // cria skill no banco de dados
        } else {
          // tem skill no banco de dados
          createdSkill = skillOfDataBase[0]; // pega primeira correspondência
        }
        if (createdSkill != null) {
          // verifica se criou ou pegou skill corretamente
          await analystService.addSkill(userSession.id, {
            skill_id: createdSkill.id,
          }); // adiciona skill no analista
        }
      }
    } catch (err) {
      console.error("Erro ao atualizar skills:", err);
      setError("Ocorreu um erro ao atualizar as skills. Tente novamente.");
      return -1;
    }
    return 0;
  }
  async function handleSubmitAll() {
    setLoading(true);
    setError("");
    setSuccess("");
    if (!userSession.id) {
      setError("Usuário não autenticado.");
      return -1;
    }

    // 1. Promover usuário a analista
    const createAnalystResult = await createAnalyst();
    if (createAnalystResult == -1) {
      setLoading(false);
      return; // se der erro na criação do analista, para o processo
    }
    // 2. Colocar Certificações
    const certificationResult = await updateCertifications();
    if (certificationResult == -1) {
      setLoading(false);
      return; // se der erro na atualização das certificações, para o processo
    }

    // 3. Colocar Skills
    const skillsResult = await updateSkills();
    if (skillsResult == -1) {
      setLoading(false);
      return; // se der erro na atualização das skills, para o processo
    }

    setSuccess("Cadastro de analista realizado com sucesso!");
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
        analystCertifications={certificationsAnalyst}
        setAnalystCertifications={setCertificationsAnalyst}
      />
      <RegisterSkills
        analystSkills={skillsAnalyst}
        setAnalystSkills={setSkillsAnalyst}
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
