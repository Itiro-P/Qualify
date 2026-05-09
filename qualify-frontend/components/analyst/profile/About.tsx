"use client";

import { useState, useEffect } from "react";
import { Reviews } from "@/components/analyst/profile/Reviews";
import { Biography } from "@/components/analyst/profile/Biography";
import { Skills } from "@/components/analyst/profile/Skills";
import { Certifications } from "@/components/analyst/profile/Certifications";
import { analystService } from "@/libs/services/analystService";
import { Analyst } from "@/types/services";
import { Skill } from "@/types/services";
import { Review } from "@/types/services/review";
import { Certification } from "@/types/services/certification";
import { skillService } from "@/libs";

function TabsSystem({ analyst }: { analyst: Analyst }) {
  const [biography, setBiography] = useState<string>("");

  const [reviewCardsVector, setReviewCardsVector] = useState<Review[]>([]);

  const [skillsCardsVector, setSkillsCardsVector] = useState<Skill[]>([]);

  const [certificationsCardsVector, setCertificationsCardsVector] = useState<
    Certification[]
  >([]);

  const [abaAtiva, setAbaAtiva] = useState<
    "biography" | "reviews" | "skills" | "certifications"
  >("biography");

  useEffect(() => {
    async function loadData() {
      try {
        // Biography
        const profileResp = await analystService.getProfile(analyst.id);

        setBiography(profileResp.analyst_profile.biography);

        // Reviews
        // precisa implementar

        // Skills
        const skillsResp = await analystService.listSkills(analyst.id);

        const skills = await Promise.all(
          skillsResp.analyst_skills.map(async (skillid) => {
            const respskill = await skillService.getById(skillid.skill_id);

            return respskill.skill;
          }),
        );

        setSkillsCardsVector(skills);

        // Certifications
        const certResp = await analystService.listCertifications(analyst.id);

        setCertificationsCardsVector(certResp.certifications);
      } catch (error: unknown) {
        if (error instanceof Error) {
          console.error(error.message);
        } else {
          //implementar erros
        }
      }
    }

    loadData();
  }, [analyst.id]);

  return (
    <div>
      <div className="flex">
        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("biography")}
        >
          Biografia
        </button>
        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("reviews")}
        >
          Reviews
        </button>
        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("skills")}
        >
          Tecnologias
        </button>
        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("certifications")}
        >
          Certificações
        </button>
      </div>

      <div className="mt-4">
        {abaAtiva === "biography" && <Biography biography={biography} />}

        {abaAtiva === "reviews" && (
          <Reviews reviewCardsVector={reviewCardsVector} />
        )}

        {abaAtiva === "skills" && (
          <Skills skillsCardsVector={skillsCardsVector} />
        )}

        {abaAtiva === "certifications" && (
          <Certifications
            certificationsCardsVector={certificationsCardsVector}
          />
        )}
      </div>
    </div>
  );
}

export function About({ analyst }: { analyst: Analyst }) {
  return (
    <section id="sobre" className="flex flex-col px-3 w-8/10">
      <TabsSystem analyst={analyst} />
    </section>
  );
}
