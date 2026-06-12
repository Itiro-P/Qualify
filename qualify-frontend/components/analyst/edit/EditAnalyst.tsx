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
import { FormButton, FormPanel, Alert, Loading } from "@/components/ui";
import { Analyst } from "@/types/services";

export function EditAnalyst({ analyst }: { analyst: Analyst }) {
  const [analystCertifications, setAnalystCertifications] = useState<
    Certification[]
  >([]);
  const [anlystSkills, setAnalystSkills] = useState<Skill[]>([]);
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<number>(-1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    async function getInfo() {
      setHourlyRateAnalyst(analyst?.hourly_rate); // pega preço por hora do analista e coloca na variavel

      const analystCertifications = await analystService.listCertifications(
        analyst.id,
      ); // pega as certificações e quantas tem com o id do analista

      if (analystCertifications != null) {
        setAnalystCertifications(analystCertifications); // pega apenas as certificações e coloca elas na variavel
      }

      const analysytSkills = await analystService.listSkills(analyst.id); // pega os ids das skills do analist
      if (analysytSkills != null) {
        for (const skillResponse of analysytSkills) {
          // para cada item dos ids dos analistas
          const skill = await skillService.getById(skillResponse.skill_id); // passa para a função apenas o id e pega apenas a skill

          if (skill != null) {
            // verifica se retornou uma skill
            setAnalystSkills((prev) => [...prev, skill]); // coloca a skill na variavel preservando as skills anteriores na variavel
          }
        }
      }
    }

    getInfo();
    console.log("Analyst ID:", analyst.id);
    console.log("Hourly Rate:", hourlyRateAnalyst);
    console.log("Certifications:", analystCertifications);
    console.log("Skills:", anlystSkills);
  }, []);

  async function updateHourlyRate(): Promise<number> {
    const result = await analystService.patch(analyst.id, {
      hourly_rate: hourlyRateAnalyst,
    }); // atualiza o preço por hora do analista
    if (result == null) {
      setError(
        "Ocorreu um erro ao atualizar o preço por hora. Tente novamente.",
      );
      return -1; // retorna -1 para indicar falha
    }

    return 0; // retorna 0 para indicar sucesso
  }

  async function updateCertifications(): Promise<number> {
    const analystCertificationsDB = await analystService.listCertifications(
      analyst.id,
    ); // pega certificações do analista no banco de dados

    if (analystCertificationsDB != null) {
      const removedCertificates = analystCertificationsDB.filter(
        (dbCert) =>
          !analystCertifications.some((userCert) => userCert.id === dbCert.id),
      ); // pega as certificações que foram removidas

      for (const certification of removedCertificates) {
        const result = await analystService.removeCertification(
          analyst.id,
          certification.id,
        );
        if (result == null) {
          setError(
            "Deu erro ao remover a certificação do analista. Tente novamente.",
          );
          return -1; // retorna -1 para indicar falha
        }
      } // remove as certificações do analista
    }

    for (const certification of analystCertifications) {
      // 'anda' por cada certificado em certificationAnalyst
      if (certification.id == -1) {
        // verifica se o id é igual a -1, ou seja é novo
        const certificationOfDataBase = await certificationService.list(
          certification.name,
        ); // pega certificação que correspondem ao nome da certificação

        let createdCertification: Certification | null;

        if (certificationOfDataBase == null) {
          // não tem certificação no banco de dados
          createdCertification = await certificationService.create({
            id: -1,
            name: certification.name,
            description: certification.description,
            institution: certification.institution,
            year: Number(certification.year),
          }); // cria certificação no banco de dados

          if (createdCertification == null) {
            setError("Deu erro ao criar a certificação. Tente novamente.");
            return -1; // retorna -1 para indicar falha
          }
        } else {
          // tem certificação no banco de dados
          createdCertification = certificationOfDataBase[0]; // pega primeira correspondência
        }

        await analystService.addCertification(analyst.id, {
          certification_id: createdCertification.id,
        }); // adiciona certificação ao analista
      }
    }
    return 0; // retorna 0 para indicar sucesso
  }

  async function updateSkills(): Promise<number> {
    const analystSkillsIds = await analystService.listSkills(analyst.id); // pega skills do analista no banco de dados

    if (analystSkillsIds != null) {
      const removedSkills = analystSkillsIds.filter(
        (dbSkill) =>
          !anlystSkills.some(
            (analystSkill) => analystSkill.id === dbSkill.skill_id,
          ),
      ); // pega skills que foram removidas pelo analista

      for (const skill of removedSkills) {
        const result = await analystService.removeSkill(
          analyst.id,
          skill.skill_id,
        ); // remove a skill do analista
        if (!result) {
          setError("Deu erro ao remover a skill do analista. Tente novamente.");
          return -1; // retorna -1 para indicar falha
        }
      } // remove as skills do analista
    }

    for (const skill of anlystSkills) {
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
        await analystService.addSkill(analyst.id, {
          skill_id: dataBaseSkill.id,
        }); // adiciona ao analista
      }
    }
    return 0; // retorna 0 para indicar sucesso
  }

  async function handleSubmitAll() {
    setLoading(true);
    setError("");
    setSuccess("");
    // 1. atualiza o preço por hora do analista
    const hourlyRateResult = await updateHourlyRate();
    if (hourlyRateResult == -1) {
      setLoading(false);
      return; // se der erro na atualização do preço por hora, para o processo
    }

    // 2. atualiza certificações
    const certificationResult = await updateCertifications();
    if (certificationResult == -1) {
      setLoading(false);
      return; // se der erro na atualização das certificações, para o processo
    }

    // 3. atualiza skills
    const skillsResult = await updateSkills();
    if (skillsResult == -1) {
      setLoading(false);
      return; // se der erro na atualização das skills, para o processo
    }
    setSuccess("Cadastro de analista realizado com sucesso!");
    setLoading(false);
  }

  return loading ? (
    <Loading />
  ) : (
    <FormPanel
      title="Editar Analista"
      description="Preencha suas certificações, tecnologias e valor por hora."
      maxWidth="max-w-2xl"
    >
      {error && <Alert variant="error">{error}</Alert>}
      {success && <Alert variant="success">{success}</Alert>}

      <RegisterCertifications
        analystCertifications={analystCertifications}
        setAnalystCertifications={setAnalystCertifications}
      />
      <RegisterSkills
        analystSkills={anlystSkills}
        setAnalystSkills={setAnalystSkills}
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
