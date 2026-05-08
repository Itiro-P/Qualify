"use client";

import { useState } from "react";
import { Reviews } from "@/components/analyst/profile/Reviews";
import { Biography } from "@/components/analyst/profile/Biography";
import { Technologies } from "@/components/analyst/profile/Technologies";
import { Certifications } from "@/components/analyst/profile/Certifications";
import { analystService } from "@/libs/services/analystService";
import { Analyst } from "@/types/services";
import { ITechnology } from "@/types/analyst/profile/technology";

function TabsSystem(analyst: Analyst) {
  const [technologiesCardsVector, setTechnologiesCardsVector] = useState<
    ITechnology[]
  >([]);
  const [abaAtiva, setAbaAtiva] = useState<
    "biography" | "reviews" | "technologies" | "certifications"
  >("biography");

  analystService.listSkills(analyst.id!).then((resp) => {
    for (const skill in resp.analyst_skills) {
      setTechnologiesCardsVector((prev) => [...prev, { technology: skill }]);
    }
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
        {abaAtiva === "biography" && <Biography />}

        {abaAtiva === "reviews" && <Reviews />}

        {abaAtiva === "technologies" && <Technologies technologiesCardsVector={technologiesCardsVector}/>}

        {abaAtiva === "certifications" && <Certifications />}
      </div>
    </div>
  );
}

export function About() {
  return (
    <section id="sobre" className="flex flex-col px-3 w-8/10">
      <TabsSystem />
    </section>
  );
}
