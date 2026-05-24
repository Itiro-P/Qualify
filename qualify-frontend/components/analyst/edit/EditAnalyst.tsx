"use client";
import { useState, useEffect } from "react";
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
import { Analyst } from "@/types/services";

export function EditAnalyst({ analyst }: { analyst: Analyst }) {
  const [certificationsAnalyst, setCertificationsAnalyst] = useState<
    Certification[]
  >([]);
  const [skillsAnalyst, setSkillsAnalyst] = useState<Skill[]>([]);
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<number>(-1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const analystId = analyst.id;
  async function get_info(analystId: number) {
    const analyst = await analystService.getByUserId(analystId); // pega analista do banco de dados

    if (analyst != null) {
      setHourlyRateAnalyst(analyst?.hourly_rate); // pega preço por hora do analista e coloca na variavel

      const analystCertifications =
        await analystService.listCertifications(analystId); // pega as certificações e quantas tem com o id do analista
      setCertificationsAnalyst(analystCertifications.certifications); // pega apenas as certificações e coloca elas na variavel

      const analysytSkills = await analystService.listSkills(analystId); // pega os ids das skills do analist
      if (analysytSkills != null) {
        for (const skillResponse of analysytSkills) {
          // para cada item dos ids dos analistas
          const skill = await skillService.getById(skillResponse.skill_id); // passa para a função apenas o id e pega apenas a skill
          if (skill != null) {
            // verifica se retornou uma skill
            setSkillsAnalyst((prev) => [...prev, skill]); // coloca a skill na variavel preservando as skills anteriores na variavel
          }
        }
      }
    }
  }

  useEffect(() => {
    get_info(analystId);
  }, [analystId]);

  async function handleSubmitAll() {
    setLoading(true);
    setError("");
    setSuccess("");
    try {
      // 1. atualiza o preço por hora do analista
      analystService.patch(analystId, { hourly_rate: hourlyRateAnalyst });

      // 2. atualiza certificações
      const analystCertifications = (
        await analystService.listCertifications(analystId)
      ).certifications; // pega certificações do analista no banco de dados

      const removedCertificates = analystCertifications.filter(
        (dbCert) =>
          !certificationsAnalyst.some((userCert) => userCert.id === dbCert.id),
      ); // pega as certificações que foram removidas

      for (const certification of removedCertificates) {
        analystService.removeCertification(analystId, certification.id);
      } // remove as certificações do analista

      for (const certification of certificationsAnalyst) {
        // 'anda' por cada certificado em certificationAnalyst
        if (certification.id == -1) {
          // verifica se o id é igual a -1, ou seja é novo
          const certificationOfDataBase = await certificationService.list(
            certification.name,
          ); // pega certificação que correspondem ao nome da certificação
          let createdCertication: Certification | null;
          if (certificationOfDataBase == null) {
            // não tem certificação no banco de dados
            createdCertication =
              await certificationService.create(certification); // cria certificação no banco de dados
          } else {
            // tem certificação no banco de dados
            createdCertication = certificationOfDataBase[0]; // pega primeira correspondência
          }
          if (createdCertication != null) {
            // verifica se criou ou pegou certification corretamente
            await analystService.addCertification(analystId, {
              certification_id: createdCertication.id,
            }); // adiciona certificação ao analista
          }
        }
      }

      // 3. atualiza skills
      const analystSkills = await analystService.listSkills(analystId); // pega skills do analista no banco de dados

      if (analystSkills != null) {
        const removedSkills = analystSkills.filter(
          (dbSkill) =>
            !skillsAnalyst.some(
              (analystSkill) => analystSkill.id === dbSkill.skill_id,
            ),
        ); // pega skills que foram removidas pelo analista

        for (const skill of removedSkills) {
          //analystService.removeSkill(analystId, skill.skill_id); // tem que corrigir a interface do remover skill do analista, pois não esta pegando o id da skill
        } // remove as skills do analista
      }

      for (const skill of skillsAnalyst) {
        if (skill.id != -1) {
          // se o id for igual a -1 significa que foi criado, se for diferente de -1 já foi criado
          continue;
        }
        const skillsOfDataBase = await skillService.list(skill.name); // pega skill que correspondem ao nome da skill
        let dataBaseSkill: Skill | null;
        if (skillsOfDataBase == null) {
          // não tem skill no banco de dados
          dataBaseSkill = await skillService.create(skill); // cria skill no banco de dados
        } else {
          // tem skill no banco de dados
          dataBaseSkill = skillsOfDataBase[0]; // pega primeira correspondência
        }
        if (dataBaseSkill != null) {
          // verifica se criou ou pegou skill corretamente
          await analystService.addSkill(analystId, {
            skill_id: dataBaseSkill.id,
          }); // adiciona ao analista
        }
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
