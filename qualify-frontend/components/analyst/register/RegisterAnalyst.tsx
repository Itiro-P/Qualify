"use client";

import { useState } from "react";
import { ICertification } from "@/types/analyst/profile/certification";
import { ITechnology } from "@/types/analyst/profile/technology";
import {
  RegisterCertifications,
  RegisterHourlyRate,
  RegisterTechnologies,
} from "@/components/analyst/register";

export function RegisterAnalyst() {
  const [certificationsAnalyst, setCertificationsAnalyst] = useState<
    ICertification[]
  >([]);
  const [technologiesAnalyst, setTechnologiesAnalyst] = useState<ITechnology[]>(
    [],
  );
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<string>("");

  function handleSubmitAll() {
    const payload = {
      certificationsAnalyst,
      technologiesAnalyst,
      hourlyRateAnalyst,
    };

    console.log("ENVIANDO TUDO:", payload);

    // aqui você chama API depois
  }

  return (
    <section>
      <div className="my-8 p-5 pt-3 border border-solid rounded-xl">
        <RegisterCertifications
          certificationsAnalyst={certificationsAnalyst}
          setCertificationsAnalyst={setCertificationsAnalyst}
        />

        <RegisterTechnologies
          technologiesAnalyst={technologiesAnalyst}
          setTechnologiesAnalyst={setTechnologiesAnalyst}
        />

        <RegisterHourlyRate
          hourlyRateAnalyst={hourlyRateAnalyst}
          setHourlyRateAnalyst={setHourlyRateAnalyst}
        />

        <button
          onClick={handleSubmitAll}
          className="mt-4 bg-blue-600 text-white font-medium px-5 py-2 rounded-lg hover:bg-blue-700 active:scale-95 transition-all duration-200"
        >
          Enviar tudo
        </button>
      </div>
    </section>
  );
}
