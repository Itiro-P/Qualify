"use client";

import { useState } from "react";
import { Reviews } from "@/components/analyst/profile/Reviews";
import { Biography } from "@/components/analyst/profile/Biography";
import { Technologies } from "@/components/analyst/profile/Technologies";
import { Certifications } from "@/components/analyst/profile/Certifications";
import { analystService } from "@/libs/services/analystService";
import { Analyst } from "@/types/services";
import { ITechnology } from "@/types/analyst/technology";
import { Review } from "@/types/services/review";
import { Certification } from "@/types/services/certification";

function TabsSystem({ analyst }: { analyst: Analyst }) {
  const [biography, setBiography] = useState<string>("");

  const [reviewCardsVector, setReviewCardsVector] = useState<Review[]>([]);

  const [technologiesCardsVector, setTechnologiesCardsVector] = useState<
    ITechnology[]
  >([]);

  const [certificationsCardsVector, setCertificationsCardsVector] = useState<
    Certification[]
  >([]);

  const [abaAtiva, setAbaAtiva] = useState<
    "biography" | "reviews" | "technologies" | "certifications"
  >("biography");

  analystService.getProfile(analyst.id).then((resp) => {
    setBiography(resp.analyst_profile.biography);
  });

  //colocar metodo para atualizar as reviewCardsVector de acordo com o 'analyst'

  analystService.listSkills(analyst.id!).then((resp) => {
    for (const skill in resp.analyst_skills) {
      setTechnologiesCardsVector((prev) => [...prev, { technology: skill }]);
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
          onClick={() => setAbaAtiva("technologies")}
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

        {abaAtiva === "technologies" && (
          <Technologies technologiesCardsVector={technologiesCardsVector} />
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
