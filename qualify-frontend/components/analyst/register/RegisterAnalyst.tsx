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
  const [hourlyRateAnalyst, setHourlyRateAnalyst] = useState<number>(0);

  async function handleSubmitAll() {
    const payload = {
      certificationsAnalyst,
      technologiesAnalyst,
      hourlyRateAnalyst,
    };

    // Primeiro vamos promover o usuário para analista.
    // E cade o id?
    try {
      const hourlyRateDTO = { hourly_rate: hourlyRateAnalyst };
      const response = await fetch(`http://localhost:8001/users/${0}/analyst`,
        {
          method: 'POST',
          headers: {
            'Content-type': 'application/json',
          },
          body: JSON.stringify(hourlyRateDTO)
        }
      );
      const data = await response.json();
      if(response.ok) {
        
      } else {
        switch (response.status) {
          case 400:
            alert(data.message || 'Dados inválidos');
            break;
          case 401:
            alert('Email ou senha incorretos');
            break;
          case 422:
            alert('Dados mal formatados');
            break;
          default:
            alert('Erro no servidor. Tente novamente.');
        }
      }
    } catch(err) {}

    // Se deu certo, adicionaremos o resto.
    try {
      // Checamos se os certificados estão lá
      for(const cert of certificationsAnalyst) {
        const response = await fetch(`http://localhost:8001/certifications/`,
          {
            method: 'POST',
            headers: {
              'Content-type': 'application/json',
            },
            body: JSON.stringify(cert)
          }
        );
        const data = await response.json();
        if(!response.ok) {
          return;
        }
      }

      for(const tech of technologiesAnalyst) {
        const response = await fetch(`http://localhost:8001/skills/`,
          {
            method: 'POST',
            headers: {
              'Content-type': 'application/json',
            },
            body: JSON.stringify(tech)
          }
        );
        const data = await response.json();
        if(!response.ok) {
          return;
        }
      }
    } catch(err) {}
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
          Cadastrar-se
        </button>
      </div>
    </section>
  );
}
