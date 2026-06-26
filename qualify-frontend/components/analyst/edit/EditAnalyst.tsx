"use client";

import { useState, useEffect } from "react";
import { Certification } from "@/types/services/certification";
import { Skill } from "@/types/services/skill";
import { EditImage } from "@/components/analyst/editImage";
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
  const [analystImage, setAnalystImage] = useState<File | undefined>();
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

      const analystSkillsList = await analystService.listSkills(analyst.id); // pega as skills completas do analista
      if (analystSkillsList != null) {
        setAnalystSkills(analystSkillsList);
      }
    }

    getInfo();
  }, []);

  async function updateImage(): Promise<number> {
    if (analystImage) {
      const result = await analystService.postImage(analyst.id, analystImage); // atualiza a imagem do analista
      if (result == null) {
        setError("Ocorreu um erro ao atualizar a imagem. Tente novamente.");
        return -1; // retorna -1 para indicar falha
      }
    }
    setError("Imagem inválida. Tente novamente.");
    return -1; // retorna -1 para indicar falha
  }

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

        await analystService.addCertification(
          analyst.id,
          createdCertification.id,
        ); // adiciona certificação ao analista
      }
    }
    return 0; // retorna 0 para indicar sucesso
  }

  async function updateSkills(): Promise<number> {
    const analystSkillsIds = await analystService.listSkills(analyst.id); // pega skills do analista no banco de dados

    if (analystSkillsIds != null) {
      const removedSkills = analystSkillsIds.filter(
        (dbSkill) =>
          !anlystSkills.some((analystSkill) => analystSkill.id === dbSkill.id),
      ); // pega skills que foram removidas pelo analista

      for (const skill of removedSkills) {
        await analystService.removeSkill(analyst.id, skill.id); // remove a skill do analista
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
        await analystService.addSkill(analyst.id, dataBaseSkill.id); // adiciona ao analista
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

    // 4. atualiza a imagem do analista
    const imageResult = await updateImage();
    if (imageResult == -1) {
      setLoading(false);
      return; // se der erro na atualização do preço por hora, para o processo
    }

    setSuccess("Analista atualizado com sucesso!");
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
      <EditImage analyst={analyst} setAnalystImage={setAnalystImage} />

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
        loadingText="Salvando..."
        className="mt-4"
      >
        Salvar alterações
      </FormButton>
    </FormPanel>
  );
}
