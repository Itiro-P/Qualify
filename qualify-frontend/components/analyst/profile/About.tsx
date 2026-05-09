"use client";

import { useState } from "react";
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

  analystService.getProfile(analyst.id).then((resp) => {
    setBiography(resp.analyst_profile.biography);
  });

  //colocar metodo para atualizar as reviewCardsVector de acordo com o 'analyst'

  analystService.listSkills(analyst.id).then((resp) => {
    for (const skillid of resp.analyst_skills) {
      skillService.getById(skillid.skill_id).then((respskill) => {
        setSkillsCardsVector((prev) => [...prev, respskill.skill]);
      });
    }
  });

  analystService.listCertifications(analyst.id).then((resp) => {
    setCertificationsCardsVector(resp.certifications);
  });

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
