"use client";

import { useState } from "react";
import { Review } from "@/components/analyst/profile/Review";
import { Biography } from "@/components/analyst/profile/Biography";
import { Technologies } from "@/components/analyst/profile/Technologies";
import { Certifications } from "@/components/analyst/profile/Certifications";

function TabsSystem() {
  const [abaAtiva, setAbaAtiva] = useState<
    "biography" | "reviews" | "technologies" | "certifications"
  >("biography");

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

        {abaAtiva === "reviews" && <Review />}

        {abaAtiva === "technologies" && <Technologies />}

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
