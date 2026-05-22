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
import type { ApiError } from "@/libs/api";
import { Alert } from "@/components/ui";

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

  const [error, setError] = useState("");

  useEffect(() => {
    async function loadData() {
      try {
        // Biography
        const profileResp = await analystService.getProfile(analyst.id);
        if (profileResp != null) {
          setBiography(profileResp.biography);
        }

        // Reviews
        const reviewsResp = await reviewService.list({
          analyst_id: analyst.id,
        });

        setReviewCardsVector(reviewsResp.reviews);

        // Skills
        const skills: Skill[] = [];
        const indexedList = await analystService.listSkills(analyst.id);
        if (indexedList != null) {
          indexedList.forEach(async (ref) => {
            const skill = await skillService.getById(ref.skill_id);
            if (skill != null) {
              skills.push(skill.skill);
            }
          });

          setSkillsCardsVector(skills);
        }

        // Certifications
        const certResp = await analystService.listCertifications(analyst.id);

        setCertificationsCardsVector(certResp.certifications);
      } catch (err) {
        const apiError = err as ApiError;
        if (apiError.status === 401) {
          setError("E-mail ou senha incorretos.");
        } else if (apiError.status === 400) {
          setError(apiError.message || "Dados inválidos.");
        }
      }
    }

    loadData();
  }, [analyst.id]);

  return (
    <div>
      <div className="flex">
        {error && <Alert variant="error">{error}</Alert>}
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
