"use client";

import { useState, useEffect } from "react";
import { Reviews } from "@/components/analyst/profile/Reviews";
import { Biography } from "@/components/analyst/profile/Biography";
import { Skills } from "@/components/analyst/profile/Skills";
import { Certifications } from "@/components/analyst/profile/Certifications";
import { analystService } from "@/libs/services/analystService";
import { reviewService } from "@/libs/services/reviewService";
import { Analyst } from "@/types/services";
import { Skill } from "@/types/services";
import { Review } from "@/types/services/review";
import { Certification } from "@/types/services/certification";
import { skillService } from "@/libs";

function TabsSystem({ analyst }: { analyst: Analyst }) {
  const [error, setError] = useState("");

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
      // Biography
      const profileResp = await analystService.getProfile(analyst.id);
      if (profileResp != null) {
        setBiography(
          profileResp.biography
            ? profileResp.biography
            : "É um usuário de poucas palavras",
        );
      }

      // Reviews
      const reviewsResp = await reviewService.list({
        analyst_id: analyst.id,
      });
      if (reviewsResp?.reviews != null) {
        setReviewCardsVector(reviewsResp.reviews);
      }

      // Skills
      const anlystSkillIds = await analystService.listSkills(analyst.id);
      if (anlystSkillIds != null) {
        const analystSkills: Skill[] = [];
        anlystSkillIds.forEach(async (ref) => {
          const analystSkill = await skillService.getById(ref.id);
          if (analystSkill != null) {
            analystSkills.push(analystSkill);
          }
        });

        setSkillsCardsVector(analystSkills);
      }

      // Certifications
      const certResp = await analystService.listCertifications(analyst.id);
      if (certResp != null) {
        setCertificationsCardsVector(certResp);
      } else {
        setError("Erro ao carregar as certificações do analista.");
      }
    }
    loadData();
    console.log("Analyst ID:", analyst.id);
    console.log("Biography:", biography);
    console.log("Reviews:", reviewCardsVector);
    console.log("Skills:", skillsCardsVector);
    console.log("Certifications:", certificationsCardsVector);
  }, []);

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

      {error && <p className="text-red-500">{error}</p>}

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
